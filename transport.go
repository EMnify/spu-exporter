package main

type Transport struct {
	Number        int
	OriginHost    string
	OriginRealm   string
	Applications  []string
	HostIps       []string
	LocalIp       string
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

type Peer struct {
	DestinationHost  string
	DestinationRealm string
	RemoteIp         string
	RemotePort       int64
	State            State
	Statistics       Statistics
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
