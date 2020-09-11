package transport

type Transport struct {
	Number        *int64
	OriginHost    string
	OriginRealm   string
	Applications  []string
	HostIps       []string
	LocalIP       string
	LocalPort     int64
	SendBuffer    int64
	ReceiveBuffer int64
	Protocol      string
	// Transport can act as client or server, in case of client in call only be one peer, as server there could be multiple
	Peers []Peer
	// for internal logic
	LastKey     string
	CurrentPeer Peer
}

func NewTransport(number int64) Transport {
	return Transport{Number: &number}
}

type Peer struct {
	Number           *int64
	DestinationHost  string
	DestinationRealm string
	RemoteIP         string
	RemotePort       int64
	State            State
	Statistics       Statistics
}

func NewPeer(number int64) Peer {
	return Peer{Number: &number}
}

type State struct {
	// can be okay, waiting and down
	Name   string
	Number int
}

type Statistics struct {
	ReceiveCnt  int64
	ReceiveMax  int64
	ReceiveAvg  int64
	ReceiveOct  int64
	ReceiveDvi  int64
	SendCnt     int64
	SendMax     int64
	SendAvg     int64
	SendOct     int64
	SendPending int64
}
