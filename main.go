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

package main

import (
	"os"

	"github.com/njcx/packetbeat6_dpdk/cmd"
	"github.com/njcx/packetbeat6_dpdk/dpdkinit"
)

var Name = "packetbeat"

// Setups and Runs Packetbeat
func main() {

	dpdkinit.DpdkInit()
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
