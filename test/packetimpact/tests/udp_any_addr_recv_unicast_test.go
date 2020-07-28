// Copyright 2020 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package udp_any_addr_recv_unicast_test

import (
	"flag"
	"net"
	"testing"

	"golang.org/x/sys/unix"
	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/test/packetimpact/testbench"
)

func init() {
	testbench.RegisterFlags(flag.CommandLine)
}

func TestAnyRecvUnicastUDP(t *testing.T) {
	dut := testbench.NewDUT(t)
	defer dut.TearDown()
	boundFD, remotePort := dut.CreateBoundSocket(t, unix.SOCK_DGRAM, unix.IPPROTO_UDP, net.IPv4(0, 0, 0, 0))
	defer dut.Close(t, boundFD)
	conn := testbench.NewUDPIPv4(t, testbench.UDP{DstPort: &remotePort}, testbench.UDP{SrcPort: &remotePort})
	defer conn.Close(t)

	payload := testbench.GenerateRandomPayload(t, 1<<10)
	conn.SendIP(
		t,
		testbench.IPv4{DstAddr: testbench.Address(tcpip.Address(net.ParseIP(testbench.RemoteIPv4).To4()))},
		testbench.UDP{},
		&testbench.Payload{Bytes: payload},
	)
	if got, want := string(dut.Recv(t, boundFD, int32(len(payload)), 0)), string(payload); got != want {
		t.Errorf("received payload does not match sent payload got: %s, want: %s", got, want)
	}
}
