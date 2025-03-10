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

//go:build !integration
// +build !integration

package mongodb

import (
	"encoding/hex"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/njcx/libbeat_v6/beat"
	"github.com/njcx/libbeat_v6/common"
	"github.com/njcx/libbeat_v6/logp"
	"github.com/njcx/packetbeat6_dpdk/protos"
)

type eventStore struct {
	events []beat.Event
}

func (e *eventStore) publish(event beat.Event) {
	e.events = append(e.events, event)
}

// Helper function returning a Mongodb module that can be used
// in tests. It publishes the transactions in the results channel.
func mongodbModForTests() (*eventStore, *mongodbPlugin) {
	var mongodb mongodbPlugin
	results := &eventStore{}
	config := defaultConfig
	mongodb.init(results.publish, &config)
	return results, &mongodb
}

// Helper function that returns an example TcpTuple
func testTCPTuple() *common.TCPTuple {
	t := &common.TCPTuple{
		IPLength: 4,
		BaseTuple: common.BaseTuple{
			SrcIP: net.IPv4(192, 168, 0, 1), DstIP: net.IPv4(192, 168, 0, 2),
			SrcPort: 6512, DstPort: 27017,
		},
	}
	t.ComputeHashables()
	return t
}

// Helper function to read from the results Queue. Raises
// an error if nothing is found in the queue.
func expectTransaction(t *testing.T, e *eventStore) common.MapStr {
	if len(e.events) == 0 {
		t.Error("No transaction")
		return nil
	}

	event := e.events[0]
	e.events = e.events[1:]
	return event.Fields
}

// Test simple request / response.
func TestSimpleFindLimit1(t *testing.T) {
	logp.TestingSetup(logp.WithSelectors("mongodb", "mongodbdetailed"))

	results, mongodb := mongodbModForTests()

	// request and response from tests/pcaps/mongo_one_row.pcap
	reqData, err := hex.DecodeString(
		"320000000a000000ffffffffd4070000" +
			"00000000746573742e72667374617572" +
			"616e7473000000000001000000050000" +
			"0000")
	assert.Nil(t, err)
	respData, err := hex.DecodeString(
		"020200004a0000000a00000001000000" +
			"08000000000000000000000000000000" +
			"01000000de010000075f696400558beb" +
			"b45f075665d2ae862703616464726573" +
			"730069000000026275696c64696e6700" +
			"05000000313030370004636f6f726400" +
			"1b000000013000e6762ff7c97652c001" +
			"3100d5b14ae9996c4440000273747265" +
			"657400100000004d6f72726973205061" +
			"726b2041766500027a6970636f646500" +
			"060000003130343632000002626f726f" +
			"756768000600000042726f6e78000263" +
			"756973696e65000700000042616b6572" +
			"79000467726164657300eb0000000330" +
			"002b00000009646174650000703d8544" +
			"01000002677261646500020000004100" +
			"1073636f72650002000000000331002b" +
			"0000000964617465000044510a410100" +
			"00026772616465000200000041001073" +
			"636f72650006000000000332002b0000" +
			"00096461746500009cda693c01000002" +
			"6772616465000200000041001073636f" +
			"7265000a000000000333002b00000009" +
			"646174650000ccb8cd33010000026772" +
			"616465000200000041001073636f7265" +
			"0009000000000334002b000000096461" +
			"7465000014109d2e0100000267726164" +
			"65000200000042001073636f7265000e" +
			"0000000000026e616d6500160000004d" +
			"6f72726973205061726b2042616b6520" +
			"53686f70000272657374617572616e74" +
			"5f696400090000003330303735343435" +
			"0000")
	assert.Nil(t, err)

	tcptuple := testTCPTuple()
	req := protos.Packet{Payload: reqData}
	resp := protos.Packet{Payload: respData}

	private := protos.ProtocolData(new(mongodbConnectionData))

	private = mongodb.Parse(&req, tcptuple, 0, private)
	mongodb.Parse(&resp, tcptuple, 1, private)
	trans := expectTransaction(t, results)

	assert.Equal(t, "OK", trans["status"])
	assert.Equal(t, "find", trans["method"])
	assert.Equal(t, "mongodb", trans["type"])

	logp.Debug("mongodb", "Trans: %v", trans)
}

// Test simple request / response, where the response is split in
// 3 messages
func TestSimpleFindLimit1_split(t *testing.T) {
	logp.TestingSetup(logp.WithSelectors("mongodb", "mongodbdetailed"))

	results, mongodb := mongodbModForTests()
	mongodb.sendRequest = true
	mongodb.sendResponse = true

	// request and response from tests/pcaps/mongo_one_row.pcap
	reqData, err := hex.DecodeString(
		"320000000a000000ffffffffd4070000" +
			"00000000746573742e72667374617572" +
			"616e7473000000000001000000050000" +
			"0000")
	assert.Nil(t, err)
	respData1, err := hex.DecodeString(
		"020200004a0000000a00000001000000" +
			"08000000000000000000000000000000" +
			"01000000de010000075f696400558beb" +
			"b45f075665d2ae862703616464726573" +
			"730069000000026275696c64696e6700" +
			"05000000313030370004636f6f726400" +
			"1b000000013000e6762ff7c97652c001" +
			"3100d5b14ae9996c4440000273747265" +
			"657400100000004d6f72726973205061")

	respData2, err := hex.DecodeString(
		"726b2041766500027a6970636f646500" +
			"060000003130343632000002626f726f" +
			"756768000600000042726f6e78000263" +
			"756973696e65000700000042616b6572" +
			"79000467726164657300eb0000000330" +
			"002b00000009646174650000703d8544" +
			"01000002677261646500020000004100" +
			"1073636f72650002000000000331002b" +
			"0000000964617465000044510a410100" +
			"00026772616465000200000041001073" +
			"636f72650006000000000332002b0000")

	respData3, err := hex.DecodeString(
		"00096461746500009cda693c01000002" +
			"6772616465000200000041001073636f" +
			"7265000a000000000333002b00000009" +
			"646174650000ccb8cd33010000026772" +
			"616465000200000041001073636f7265" +
			"0009000000000334002b000000096461" +
			"7465000014109d2e0100000267726164" +
			"65000200000042001073636f7265000e" +
			"0000000000026e616d6500160000004d" +
			"6f72726973205061726b2042616b6520" +
			"53686f70000272657374617572616e74" +
			"5f696400090000003330303735343435" +
			"0000")
	assert.Nil(t, err)

	tcptuple := testTCPTuple()
	req := protos.Packet{Payload: reqData}

	private := protos.ProtocolData(new(mongodbConnectionData))

	private = mongodb.Parse(&req, tcptuple, 0, private)

	resp1 := protos.Packet{Payload: respData1}
	private = mongodb.Parse(&resp1, tcptuple, 1, private)

	resp2 := protos.Packet{Payload: respData2}
	private = mongodb.Parse(&resp2, tcptuple, 1, private)

	resp3 := protos.Packet{Payload: respData3}
	mongodb.Parse(&resp3, tcptuple, 1, private)

	trans := expectTransaction(t, results)

	assert.Equal(t, "OK", trans["status"])
	assert.Equal(t, "find", trans["method"])
	assert.Equal(t, "mongodb", trans["type"])

	logp.Debug("mongodb", "Trans: %v", trans)
}

func TestReconstructQuery(t *testing.T) {
	type io struct {
		Input  transaction
		Full   bool
		Output string
	}
	tests := []io{
		{
			Input: transaction{
				resource: "test.col",
				method:   "find",
				event: map[string]interface{}{
					"numberToSkip":   3,
					"numberToReturn": 2,
				},
				params: map[string]interface{}{
					"me": "you",
				},
			},
			Full:   true,
			Output: `test.col.find({"me":"you"}).skip(3).limit(2)`,
		},
		{
			Input: transaction{
				resource: "test.col",
				method:   "insert",
				params: map[string]interface{}{
					"documents": "you",
				},
			},
			Full:   true,
			Output: `test.col.insert({"documents":"you"})`,
		},
		{
			Input: transaction{
				resource: "test.col",
				method:   "insert",
				params: map[string]interface{}{
					"documents": "you",
				},
			},
			Full:   false,
			Output: `test.col.insert({})`,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Output,
			reconstructQuery(&test.Input, test.Full))
	}
}

// max_docs option should be respected
func TestMaxDocs(t *testing.T) {
	logp.TestingSetup(logp.WithSelectors("mongodb", "mongodbdetailed"))

	// more docs than configured
	trans := transaction{
		documents: []interface{}{
			1, 2, 3, 4, 5, 6, 7, 8,
		},
	}

	results, mongodb := mongodbModForTests()
	mongodb.sendResponse = true
	mongodb.maxDocs = 3

	mongodb.publishTransaction(&trans)

	res := expectTransaction(t, results)

	assert.Equal(t, "1\n2\n3\n[...]", res["response"])

	// exactly the same number of docs
	trans = transaction{
		documents: []interface{}{
			1, 2, 3,
		},
	}

	mongodb.publishTransaction(&trans)
	res = expectTransaction(t, results)
	assert.Equal(t, "1\n2\n3", res["response"])

	// less docs
	trans = transaction{
		documents: []interface{}{
			1, 2,
		},
	}

	mongodb.publishTransaction(&trans)
	res = expectTransaction(t, results)
	assert.Equal(t, "1\n2", res["response"])

	// unlimited
	trans = transaction{
		documents: []interface{}{
			1, 2, 3, 4,
		},
	}
	mongodb.maxDocs = 0
	mongodb.publishTransaction(&trans)
	res = expectTransaction(t, results)
	assert.Equal(t, "1\n2\n3\n4", res["response"])
}

func TestMaxDocSize(t *testing.T) {
	logp.TestingSetup(logp.WithSelectors("mongodb", "mongodbdetailed"))

	// more docs than configured
	trans := transaction{
		documents: []interface{}{
			"1234567",
			"123",
			"12",
		},
	}

	results, mongodb := mongodbModForTests()
	mongodb.sendResponse = true
	mongodb.maxDocLength = 5

	mongodb.publishTransaction(&trans)

	res := expectTransaction(t, results)

	assert.Equal(t, "\"1234 ...\n\"123\"\n\"12\"", res["response"])
}

func TestOpCodeNames(t *testing.T) {
	for _, testData := range []struct {
		code     int32
		expected string
	}{
		{1, "OP_REPLY"},
		{-1, "(value=-1)"},
	} {
		assert.Equal(t, testData.expected, opCode(testData.code).String())
	}
}

// Test for a (recovered) panic parsing document length in request/response messages
func TestDocumentLengthBoundsChecked(t *testing.T) {
	logp.TestingSetup(logp.WithSelectors("mongodb", "mongodbdetailed"))

	_, mongodb := mongodbModForTests()

	// request and response from tests/pcaps/mongo_one_row.pcap
	reqData, err := hex.DecodeString(
		// Request message with out of bounds document
		"320000000a000000ffffffffd4070000" +
			"00000000746573742e72667374617572" +
			"616e7473000000000001000000" +
			// Document length (including itself)
			"06000000" +
			// Document (1 byte instead of 2)
			"00")
	assert.Nil(t, err)

	tcptuple := testTCPTuple()
	req := protos.Packet{Payload: reqData}
	private := protos.ProtocolData(new(mongodbConnectionData))

	private = mongodb.Parse(&req, tcptuple, 0, private)
	assert.NotNil(t, private, "mongodb parser recovered from a panic")
}
