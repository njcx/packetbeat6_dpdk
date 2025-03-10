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

package protos

import (
	"testing"
	"time"

	"github.com/njcx/libbeat_v6/common"

	"github.com/stretchr/testify/assert"
)

type TestProtocol struct {
	Ports []int
}

type TCPProtocol TestProtocol

func (proto *TCPProtocol) Init(testMode bool, results Reporter) error {
	return nil
}

func (proto *TCPProtocol) GetPorts() []int {
	return proto.Ports
}

func (proto *TCPProtocol) Parse(pkt *Packet, tcptuple *common.TCPTuple,
	dir uint8, private ProtocolData) ProtocolData {
	return private
}

func (proto *TCPProtocol) ReceivedFin(tcptuple *common.TCPTuple, dir uint8,
	private ProtocolData) ProtocolData {
	return private
}

func (proto *TCPProtocol) GapInStream(tcptuple *common.TCPTuple, dir uint8,
	nbytes int, private ProtocolData) (priv ProtocolData, drop bool) {
	return private, true
}

func (proto *TCPProtocol) ConnectionTimeout() time.Duration { return 0 }

type UDPProtocol TestProtocol

func (proto *UDPProtocol) Init(testMode bool, results Reporter) error {
	return nil
}

func (proto *UDPProtocol) GetPorts() []int {
	return proto.Ports
}

func (proto *UDPProtocol) ParseUDP(pkt *Packet) {
	return
}

type TCPUDPProtocol TestProtocol

func (proto *TCPUDPProtocol) Init(testMode bool, results Reporter) error {
	return nil
}

func (proto *TCPUDPProtocol) GetPorts() []int {
	return proto.Ports
}

func (proto *TCPUDPProtocol) Parse(pkt *Packet, tcptuple *common.TCPTuple,
	dir uint8, private ProtocolData) ProtocolData {
	return private
}

func (proto *TCPUDPProtocol) ReceivedFin(tcptuple *common.TCPTuple, dir uint8,
	private ProtocolData) ProtocolData {
	return private
}

func (proto *TCPUDPProtocol) GapInStream(tcptuple *common.TCPTuple, dir uint8,
	nbytes int, private ProtocolData) (priv ProtocolData, drop bool) {
	return private, true
}

func (proto *TCPUDPProtocol) ParseUDP(pkt *Packet) {
	return
}

func (proto *TCPUDPProtocol) ConnectionTimeout() time.Duration { return 0 }

func TestProtocolNames(t *testing.T) {
	assert.Equal(t, "unknown", UnknownProtocol.String())
	assert.Equal(t, "impossible", Protocol(100).String())
}

func newProtocols() Protocols {
	p := ProtocolsStruct{}
	p.all = make(map[Protocol]protocolInstance)
	p.tcp = make(map[Protocol]TCPPlugin)
	p.udp = make(map[Protocol]UDPPlugin)

	tcp := &TCPProtocol{Ports: []int{80}}
	udp := &UDPProtocol{Ports: []int{5060}}
	tcpUDP := &TCPUDPProtocol{Ports: []int{53}}

	p.register(1, nil, tcp)
	p.register(2, nil, udp)
	p.register(3, nil, tcpUDP)
	return p
}

func TestBpfFilterWithoutVlanOnlyIcmp(t *testing.T) {
	p := ProtocolsStruct{}
	p.all = make(map[Protocol]protocolInstance)
	p.tcp = make(map[Protocol]TCPPlugin)
	p.udp = make(map[Protocol]UDPPlugin)

	filter := p.BpfFilter(false, true)
	assert.Equal(t, "icmp or icmp6", filter)
}

func TestBpfFilterWithoutVlanWithoutIcmp(t *testing.T) {
	p := newProtocols()
	filter := p.BpfFilter(false, false)
	assert.Equal(t, "tcp port 80 or udp port 5060 or port 53", filter)
}

func TestBpfFilterWithVlanWithoutIcmp(t *testing.T) {
	p := newProtocols()
	filter := p.BpfFilter(true, false)
	assert.Equal(t, "tcp port 80 or udp port 5060 or port 53 or "+
		"(vlan and (tcp port 80 or udp port 5060 or port 53))", filter)
}

func TestBpfFilterWithoutVlanWithIcmp(t *testing.T) {
	p := newProtocols()
	filter := p.BpfFilter(false, true)
	assert.Equal(t, "tcp port 80 or udp port 5060 or port 53 or icmp or icmp6", filter)
}

func TestBpfFilterWithVlanWithIcmp(t *testing.T) {
	p := newProtocols()
	filter := p.BpfFilter(true, true)
	assert.Equal(t, "tcp port 80 or udp port 5060 or port 53 or icmp or icmp6 or "+
		"(vlan and (tcp port 80 or udp port 5060 or port 53 or icmp or icmp6))", filter)
}

func TestGetAllTCP(t *testing.T) {
	p := newProtocols()
	tcp := p.GetAllTCP()
	assert.NotNil(t, tcp[1])
	assert.Nil(t, tcp[2])
	assert.NotNil(t, tcp[3])
}

func TestGetAllUDP(t *testing.T) {
	p := newProtocols()
	udp := p.GetAllUDP()
	assert.Nil(t, udp[1])
	assert.NotNil(t, udp[2])
	assert.NotNil(t, udp[3])
}

func TestGetTCP(t *testing.T) {
	p := newProtocols()
	tcp := p.GetTCP(1)
	assert.NotNil(t, tcp)
	assert.Contains(t, tcp.GetPorts(), 80)

	tcp = p.GetTCP(2)
	assert.Nil(t, tcp)

	tcp = p.GetTCP(3)
	assert.NotNil(t, tcp)
	assert.Contains(t, tcp.GetPorts(), 53)
}

func TestGetUDP(t *testing.T) {
	p := newProtocols()
	udp := p.GetUDP(1)
	assert.Nil(t, udp)

	udp = p.GetUDP(2)
	assert.NotNil(t, udp)
	assert.Contains(t, udp.GetPorts(), 5060)

	udp = p.GetUDP(3)
	assert.NotNil(t, udp)
	assert.Contains(t, udp.GetPorts(), 53)
}
