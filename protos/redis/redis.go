// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package redis

import (
	"bytes"
	"time"

	"github.com/njcx/libbeat_v6/beat"
	"github.com/njcx/libbeat_v6/common"
	"github.com/njcx/libbeat_v6/logp"
	"github.com/njcx/libbeat_v6/monitoring"

	"github.com/njcx/packetbeat6_dpdk/procs"
	"github.com/njcx/packetbeat6_dpdk/protos"
	"github.com/njcx/packetbeat6_dpdk/protos/applayer"
	"github.com/njcx/packetbeat6_dpdk/protos/tcp"
)

type stream struct {
	applayer.Stream
	parser   parser
	tcptuple *common.TCPTuple
}

type redisConnectionData struct {
	streams   [2]*stream
	requests  MessageQueue
	responses MessageQueue
}

// Redis protocol plugin
type redisPlugin struct {
	// config
	ports              []int
	sendRequest        bool
	sendResponse       bool
	transactionTimeout time.Duration
	queueConfig        MessageQueueConfig

	results protos.Reporter
}

var (
	debugf  = logp.MakeDebug("redis")
	isDebug = false
)

var (
	unmatchedResponses = monitoring.NewInt(nil, "redis.unmatched_responses")
	unmatchedRequests  = monitoring.NewInt(nil, "redis.unmatched_requests")
)

func init() {
	protos.Register("redis", New)
}

func New(
	testMode bool,
	results protos.Reporter,
	cfg *common.Config,
) (protos.Plugin, error) {
	p := &redisPlugin{}
	config := defaultConfig
	if !testMode {
		if err := cfg.Unpack(&config); err != nil {
			return nil, err
		}
	}

	if err := p.init(results, &config); err != nil {
		return nil, err
	}
	return p, nil
}

func (redis *redisPlugin) init(results protos.Reporter, config *redisConfig) error {
	redis.setFromConfig(config)

	redis.results = results
	isDebug = logp.IsDebug("redis")

	return nil
}

func (redis *redisPlugin) setFromConfig(config *redisConfig) {
	redis.ports = config.Ports
	redis.sendRequest = config.SendRequest
	redis.sendResponse = config.SendResponse
	redis.transactionTimeout = config.TransactionTimeout
	redis.queueConfig = config.QueueLimits
}

func (redis *redisPlugin) GetPorts() []int {
	return redis.ports
}

func (s *stream) PrepareForNewMessage() {
	parser := &s.parser
	s.Stream.Reset()
	parser.reset()
}

func (redis *redisPlugin) ConnectionTimeout() time.Duration {
	return redis.transactionTimeout
}

func (redis *redisPlugin) Parse(
	pkt *protos.Packet,
	tcptuple *common.TCPTuple,
	dir uint8,
	private protos.ProtocolData,
) protos.ProtocolData {
	defer logp.Recover("ParseRedis exception")

	conn := redis.ensureRedisConnection(private)
	conn = redis.doParse(conn, pkt, tcptuple, dir)
	if conn == nil {
		return nil
	}
	return conn
}
func (redis *redisPlugin) newConnectionData() *redisConnectionData {
	return &redisConnectionData{
		requests:  NewMessageQueue(redis.queueConfig),
		responses: NewMessageQueue(redis.queueConfig),
	}
}

func (redis *redisPlugin) ensureRedisConnection(private protos.ProtocolData) *redisConnectionData {
	if private == nil {
		return redis.newConnectionData()
	}

	priv, ok := private.(*redisConnectionData)
	if !ok {
		logp.Warn("redis connection data type error, create new one")
		return redis.newConnectionData()
	}
	if priv == nil {
		logp.Warn("Unexpected: redis connection data not set, create new one")
		return redis.newConnectionData()
	}

	return priv
}

func (redis *redisPlugin) doParse(
	conn *redisConnectionData,
	pkt *protos.Packet,
	tcptuple *common.TCPTuple,
	dir uint8,
) *redisConnectionData {

	st := conn.streams[dir]
	if st == nil {
		st = newStream(pkt.Ts, tcptuple)
		conn.streams[dir] = st
		if isDebug {
			debugf("new stream: %p (dir=%v, len=%v)", st, dir, len(pkt.Payload))
		}
	}

	if err := st.Append(pkt.Payload); err != nil {
		if isDebug {
			debugf("%v, dropping TCP stream: ", err)
		}
		return nil
	}
	if isDebug {
		debugf("stream add data: %p (dir=%v, len=%v)", st, dir, len(pkt.Payload))
	}

	for st.Buf.Len() > 0 {
		if st.parser.message == nil {
			st.parser.message = newMessage(pkt.Ts)
		}

		ok, complete := st.parser.parse(&st.Buf)
		if !ok {
			// drop this tcp stream. Will retry parsing with the next
			// segment in it
			conn.streams[dir] = nil
			if isDebug {
				debugf("Ignore Redis message. Drop tcp stream. Try parsing with the next segment")
			}
			return conn
		}

		if !complete {
			// wait for more data
			break
		}

		msg := st.parser.message
		if isDebug {
			if msg.isRequest {
				debugf("REDIS (%p) request message: %s", conn, msg.message)
			} else {
				debugf("REDIS (%p) response message: %s", conn, msg.message)
			}
		}

		// all ok, go to next level and reset stream for new message
		redis.handleRedis(conn, msg, tcptuple, dir)
		st.PrepareForNewMessage()
	}

	return conn
}

func newStream(ts time.Time, tcptuple *common.TCPTuple) *stream {
	s := &stream{
		tcptuple: tcptuple,
	}
	s.parser.message = newMessage(ts)
	s.Stream.Init(tcp.TCPMaxDataInStream)
	return s
}

func newMessage(ts time.Time) *redisMessage {
	return &redisMessage{ts: ts}
}

func (redis *redisPlugin) handleRedis(
	conn *redisConnectionData,
	m *redisMessage,
	tcptuple *common.TCPTuple,
	dir uint8,
) {
	m.tcpTuple = *tcptuple
	m.direction = dir
	m.cmdlineTuple = procs.ProcWatcher.FindProcessesTupleTCP(tcptuple.IPPort())

	if m.isRequest {
		// wait for response
		if evicted := conn.requests.Append(m); evicted > 0 {
			unmatchedRequests.Add(int64(evicted))
		}
	} else {
		if evicted := conn.responses.Append(m); evicted > 0 {
			unmatchedResponses.Add(int64(evicted))
		}
		redis.correlate(conn)
	}
}

func (redis *redisPlugin) correlate(conn *redisConnectionData) {
	// drop responses with missing requests
	if conn.requests.IsEmpty() {
		for !conn.responses.IsEmpty() {
			debugf("Response from unknown transaction. Ignoring")
			unmatchedResponses.Add(1)
			conn.responses.Pop()
		}
		return
	}

	// merge requests with responses into transactions
	for !conn.responses.IsEmpty() && !conn.requests.IsEmpty() {
		requ, okReq := conn.requests.Pop().(*redisMessage)
		resp, okResp := conn.responses.Pop().(*redisMessage)
		if !okReq || !okResp {
			logp.Err("invalid type found in message queue")
			continue
		}
		if redis.results != nil {
			event := redis.newTransaction(requ, resp)
			redis.results(event)
		}
	}
}

func (redis *redisPlugin) newTransaction(requ, resp *redisMessage) beat.Event {
	error := common.OK_STATUS
	if resp.isError {
		error = common.ERROR_STATUS
	}

	var returnValue map[string]common.NetString
	if resp.isError {
		returnValue = map[string]common.NetString{
			"error": resp.message,
		}
	} else {
		returnValue = map[string]common.NetString{
			"return_value": resp.message,
		}
	}

	source, destination := common.MakeEndpointPair(requ.tcpTuple.BaseTuple, requ.cmdlineTuple)
	src, dst := &source, &destination
	if requ.direction == tcp.TCPDirectionReverse {
		src, dst = dst, src
	}

	// resp_time in milliseconds
	responseTime := int32(resp.ts.Sub(requ.ts).Nanoseconds() / 1e6)

	fields := common.MapStr{
		"type":         "redis",
		"status":       error,
		"responsetime": responseTime,
		"redis":        returnValue,
		"method":       common.NetString(bytes.ToUpper(requ.method)),
		"resource":     requ.path,
		"query":        requ.message,
		"bytes_in":     uint64(requ.size),
		"bytes_out":    uint64(resp.size),
		"src":          src,
		"dst":          dst,
	}
	if redis.sendRequest {
		fields["request"] = requ.message
	}
	if redis.sendResponse {
		fields["response"] = resp.message
	}

	return beat.Event{
		Timestamp: requ.ts,
		Fields:    fields,
	}
}

func (redis *redisPlugin) GapInStream(tcptuple *common.TCPTuple, dir uint8,
	nbytes int, private protos.ProtocolData) (priv protos.ProtocolData, drop bool) {

	// tsg: being packet loss tolerant is probably not very useful for Redis,
	// because most requests/response tend to fit in a single packet.

	return private, true
}

func (redis *redisPlugin) ReceivedFin(tcptuple *common.TCPTuple, dir uint8,
	private protos.ProtocolData) protos.ProtocolData {

	// TODO: check if we have pending data that we can send up the stack

	return private
}
