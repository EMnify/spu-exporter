package collector

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/EMnify/spu-exporter/pkg/transport"
)

var parsetests = []struct {
	inFile        string
	expectFailure bool
	expectedTrans []transport.Transport
}{
	{"../../test/example", false, testExampleTransports()},
}

func TestParseLines(t *testing.T) {
	for _, tt := range parsetests {
		t.Run(tt.inFile, func(t *testing.T) {
			absPath, _ := filepath.Abs(tt.inFile)
			transports, err := parseLines(readFromFile(absPath))
			if tt.expectFailure {
				if err == nil {
					t.Errorf("expected failure, but passed")
				}
			} else {
				if err != nil {
					t.Errorf("got parse error when it should not fail")
				} else {
					if len(transports) != len(tt.expectedTrans) {
						t.Errorf("got different number of transports as parse result")
					} else {
						for i := range transports {
							compareTransport(transports[i], tt.expectedTrans[i], t)
						}
					}
				}
			}
		})
	}
}

func readFromFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}

func compareTransport(t1, t2 transport.Transport, t *testing.T) {
	assert.Equal(t, *t2.Number, *t1.Number, "Number expected to be the same")
	assert.Equal(t, t2.Protocol, t1.Protocol, "Protocol expected to be the same")
	assert.Equal(t, t2.OriginHost, t1.OriginHost, "OriginHost expected to be the same")
	assert.Equal(t, t2.OriginRealm, t1.OriginRealm, "OriginRealm expected to be the same")
	if len(t1.Peers) == len(t2.Peers) {
		for i := range t1.Peers {
			comparePeers(t1.Peers[i], t2.Peers[i], t)
		}
	} else {
		t.Errorf("Peer count expected to be the same, expected %d, was %d", len(t2.Peers), len(t1.Peers))
	}

}

func comparePeers(p1, p2 transport.Peer, t *testing.T) {
	assert.Equal(t, p2.RemoteIP, p1.RemoteIP, "RemoteIP expected to be the same")
	assert.Equal(t, p2.RemotePort, p1.RemotePort, "RemotePort expected to be the same")
	assert.Equal(t, p2.DestinationHost, p1.DestinationHost, "DestinationHost expected to be the same")
	assert.Equal(t, p2.DestinationRealm, p1.DestinationRealm, "DestinationRealm expected to be the same")
	assert.Equal(t, p2.State, p1.State, "State expected to be the same")
}

//func compareString(s1, s2, name string, diff []string, before bool) (bool, []string) {
//	if s1 != s2 {
//		return false, append(diff, fmt.Sprintf("%s was %s exptected %s", name, s1, s2))
//	}
//	return before, diff
//}
//func compareInt64(s1, s2 int64, name string, diff []string, before bool) (bool, []string) {
//	if s1 != s2 {
//		return false, append(diff, fmt.Sprintf("%s was %d exptected %d", name, s1, s2))
//	}
//	return before, diff
//}

func testExampleTransports() []transport.Transport {
	transports := []transport.Transport{}
	t1 := transport.NewTransport(0)
	t1.OriginHost = "hss.epc.mnc012.mcc901.3gppnetwork.org"
	t1.OriginRealm = "epc.mnc012.mcc901.3gppnetwork.org"
	t1.Protocol = "sctp"
	peers1 := []transport.Peer{}
	p1 := transport.NewPeer(0)
	p1.DestinationRealm = "dest.abc.3gppnetwork.org"
	p1.DestinationHost = "dest123.abc.3gppnetwork.org"
	p1.RemoteIP = "12.123.123.123"
	p1.RemotePort = 3868
	p1.State = transport.State{
		Name:   "okay",
		Number: 0,
	}
	t1.Peers = append(peers1, p1)
	transports = append(transports, t1)
	t2 := transport.NewTransport(1)
	t2.OriginHost = "hss.epc.mnc034.mcc123.3gppnetwork.org"
	t2.OriginRealm = "epc.mnc034.mcc123.3gppnetwork.org"
	t2.Protocol = "sctp"
	peers2 := []transport.Peer{}
	p2 := transport.NewPeer(0)
	p2.DestinationRealm = "abc.3gppnetwork.org"
	p2.DestinationHost = "dest2345.abc.3gppnetwork.org"
	p2.RemoteIP = "12.234.234.234"
	p2.RemotePort = 3868
	p2.State = transport.State{
		Name:   "okay",
		Number: 0,
	}
	t2.Peers = append(peers2, p2)

	transports = append(transports, t2)
	//t3 := transport.NewTransport(1)
	//transports = append(transports, t3)
	//t4 := transport.NewTransport(1)
	//transports = append(transports, t4)
	//t5 := transport.NewTransport(1)
	//transports = append(transports, t5)
	//t6 := transport.NewTransport(1)
	//transports = append(transports, t6)

	return transports
}
