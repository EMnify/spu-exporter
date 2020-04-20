package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

func createMetricLines(ts []Transport) *prometheus.Registry {
	reg := prometheus.NewRegistry()
	state := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_state", Help: "State of the transport (labels okay, waiting, down with 1 or 0)"}, []string{"transport", "origin_host", "destination_host", "remote_ip", "state"})
	reg.MustRegister(state)

	recvCnt := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_receive_cnt_total", Help: "Number of packets received by the socket."}, []string{"transport", "origin_host", "destination_host", "remote_ip"})
	recvAvg := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_receive_avg", Help: "Average size of packets, in bytes, received by the socket."}, []string{"transport", "origin_host", "destination_host", "remote_ip"})
	recvMax := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_receive_max", Help: "Size of the largest packet, in bytes, received by the socket."}, []string{"transport", "origin_host", "destination_host", "remote_ip"})
	recvOct := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_receive_oct_total", Help: "Number of bytes received by the socket."}, []string{"transport", "origin_host", "destination_host", "remote_ip"})
	recvDvi := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_receive_dvi", Help: "Average packet size deviation, in bytes, received by the socket."}, []string{"transport", "origin_host", "destination_host", "remote_ip"})
	reg.MustRegister(recvCnt, recvAvg, recvDvi, recvMax, recvOct)

	sendAvg := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_send_avg", Help: "Average size of packets, in bytes, sent from the socket."}, []string{"transport", "origin_host", "destination_host", "remote_ip"})
	sendCnt := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_send_cnt_total", Help: "Number of packets sent from the socket."}, []string{"transport", "origin_host", "destination_host", "remote_ip"})
	sendPend := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_send_pending", Help: "Number of bytes waiting to be sent by the socket."}, []string{"transport", "origin_host", "destination_host", "remote_ip"})
	sendMax := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_send_max", Help: "Size of the largest packet, in bytes, sent from the socket."}, []string{"transport", "origin_host", "destination_host", "remote_ip"})
	sendOct := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_send_oct_total", Help: "Number of bytes sent from the socket."}, []string{"transport", "origin_host", "destination_host", "remote_ip"})
	reg.MustRegister(sendAvg, sendCnt, sendMax, sendOct, sendPend)

	for _, t := range ts {
		if t.OriginHost != "" {
			for _, p := range t.Peers {
				switch p.State.Name {
				case "okay":
					state.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp, "okay").Set(1)
					state.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp, "waiting").Set(0)
					state.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp, "down").Set(0)
				case "waiting":
					state.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp, "okay").Set(0)
					state.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp, "waiting").Set(1)
					state.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp, "down").Set(0)
				case "down":
					state.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp, "okay").Set(0)
					state.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp, "waiting").Set(0)
					state.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp, "down").Set(1)
				}

				// receive stats
				recvCnt.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp).Set(float64(p.Statistics.ReceiveCnt))
				recvAvg.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp).Set(float64(p.Statistics.ReceiveAvg))
				recvMax.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp).Set(float64(p.Statistics.ReceiveMax))
				recvOct.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp).Set(float64(p.Statistics.ReceiveOct))
				recvDvi.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp).Set(float64(p.Statistics.ReceiveDvi))
				sendAvg.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp).Set(float64(p.Statistics.SendAvg))
				sendCnt.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp).Set(float64(p.Statistics.SendCnt))
				sendPend.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp).Set(float64(p.Statistics.SendPending))
				sendMax.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp).Set(float64(p.Statistics.SendMax))
				sendOct.WithLabelValues(strconv.Itoa(t.Number), t.OriginHost, p.DestinationHost, p.RemoteIp).Set(float64(p.Statistics.SendOct))
			}
		}
	}
	return reg
}

func writeToFile(reg *prometheus.Registry, filename string) {
	gatherer := prometheus.Gatherers{reg}
	prometheus.WriteToTextfile(filename, gatherer)
}
