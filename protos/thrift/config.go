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

package thrift

import (
	"github.com/njcx/packetbeat6_dpdk/config"
	"github.com/njcx/packetbeat6_dpdk/protos"
)

type thriftConfig struct {
	config.ProtocolCommon  `config:",inline"`
	StringMaxSize          int      `config:"string_max_size"`
	CollectionMaxSize      int      `config:"collection_max_size"`
	DropAfterNStructFields int      `config:"drop_after_n_struct_fields"`
	TransportType          string   `config:"transport_type"`
	ProtocolType           string   `config:"protocol_type"`
	CaptureReply           bool     `config:"capture_reply"`
	ObfuscateStrings       bool     `config:"obfuscate_strings"`
	IdlFiles               []string `config:"idl_files"`
}

var (
	defaultConfig = thriftConfig{
		ProtocolCommon: config.ProtocolCommon{
			TransactionTimeout: protos.DefaultTransactionExpiration,
		},
		StringMaxSize:          200,
		CollectionMaxSize:      15,
		DropAfterNStructFields: 500,
		TransportType:          "socket",
		ProtocolType:           "binary",
		CaptureReply:           true,
	}
)
