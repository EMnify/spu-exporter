package prom

import (
	"strconv"

	"github.com/EMnify/spu-exporter/pkg/transport"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	state *prometheus.GaugeVec

	recvCnt *prometheus.GaugeVec
	recvAvg *prometheus.GaugeVec
	recvMax *prometheus.GaugeVec
	recvOct *prometheus.GaugeVec
	recvDvi *prometheus.GaugeVec

	sendAvg  *prometheus.GaugeVec
	sendCnt  *prometheus.GaugeVec
	sendPend *prometheus.GaugeVec
	sendMax  *prometheus.GaugeVec
	sendOct  *prometheus.GaugeVec

	memoryTotal         prometheus.Gauge
	memoryProcesses     prometheus.Gauge
	memoryProcessesUsed prometheus.Gauge
	memoryTotal         prometheus.Gauge
	memoryTotal         prometheus.Gauge
	memoryTotal         prometheus.Gauge
	memoryTotal         prometheus.Gauge
)

func RegisterMetrics(reg *prometheus.Registry) {

	labels := []string{"transport", "origin_host", "destination_host", "remote_ip"}
	state = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_state", Help: "State of the transport (labels okay, waiting, down with 1 or 0)"}, append(labels, "state"))
	reg.MustRegister(state)

	recvCnt = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_receive_cnt_total", Help: "Number of packets received by the socket."}, labels)
	recvAvg = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_receive_avg", Help: "Average size of packets, in bytes, received by the socket."}, labels)
	recvMax = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_receive_max", Help: "Size of the largest packet, in bytes, received by the socket."}, labels)
	recvOct = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_receive_oct_total", Help: "Number of bytes received by the socket."}, labels)
	recvDvi = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_receive_dvi", Help: "Average packet size deviation, in bytes, received by the socket."}, labels)
	reg.MustRegister(recvCnt, recvAvg, recvDvi, recvMax, recvOct)

	sendAvg = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_send_avg", Help: "Average size of packets, in bytes, sent from the socket."}, labels)
	sendCnt = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_send_cnt_total", Help: "Number of packets sent from the socket."}, labels)
	sendPend = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_send_pending", Help: "Number of bytes waiting to be sent by the socket."}, labels)
	sendMax = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_send_max", Help: "Size of the largest packet, in bytes, sent from the socket."}, labels)
	sendOct = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "spu_transport_send_oct_total", Help: "Number of bytes sent from the socket."}, labels)
	reg.MustRegister(sendAvg, sendCnt, sendMax, sendOct, sendPend)

	memoryTotal = promauto.NewGauge(prometheus.GaugeOpts{Name: "spu_memory_total", Help: "Total size of spu memory, in bytes"})
	memoryProcesses = promauto.NewGauge(prometheus.GaugeOpts{Name: "spu_memory_processes", Help: "Processes size of spu memory, in bytes"})
	memoryProcessesUsed = promauto.NewGauge(prometheus.GaugeOpts{Name: "spu_memory_processes_used", Help: "Processes Used size of spu memory, in bytes"})
	memoryTotal = promauto.NewGauge(prometheus.GaugeOpts{Name: "spu_memory_total", Help: "Total size of spu memory, in bytes"})
	memoryTotal = promauto.NewGauge(prometheus.GaugeOpts{Name: "spu_memory_total", Help: "Total size of spu memory, in bytes"})
	memoryTotal = promauto.NewGauge(prometheus.GaugeOpts{Name: "spu_memory_total", Help: "Total size of spu memory, in bytes"})
	memoryTotal = promauto.NewGauge(prometheus.GaugeOpts{Name: "spu_memory_total", Help: "Total size of spu memory, in bytes"})
	memoryTotal = promauto.NewGauge(prometheus.GaugeOpts{Name: "spu_memory_total", Help: "Total size of spu memory, in bytes"})
	memoryTotal = promauto.NewGauge(prometheus.GaugeOpts{Name: "spu_memory_total", Help: "Total size of spu memory, in bytes"})

}

func CreateMetricLines(ts *[]transport.Transport, reg *prometheus.Registry) *prometheus.Registry {

	for _, t := range *ts {
		if t.OriginHost != "" {
			for _, p := range t.Peers {
				switch p.State.Name {
				case "okay":
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "okay").Set(1)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "waiting").Set(0)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "down").Set(0)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "initial").Set(0)
				case "waiting":
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "okay").Set(0)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "waiting").Set(1)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "down").Set(0)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "initial").Set(0)
				case "down":
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "okay").Set(0)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "waiting").Set(0)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "down").Set(1)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "initial").Set(0)
				case "initial":
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "okay").Set(0)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "waiting").Set(0)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "down").Set(0)
					state.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP, "initial").Set(1)
				}

				// receive stats
				recvCnt.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP).Set(float64(p.Statistics.ReceiveCnt))
				recvAvg.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP).Set(float64(p.Statistics.ReceiveAvg))
				recvMax.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP).Set(float64(p.Statistics.ReceiveMax))
				recvOct.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP).Set(float64(p.Statistics.ReceiveOct))
				recvDvi.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP).Set(float64(p.Statistics.ReceiveDvi))
				sendAvg.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP).Set(float64(p.Statistics.SendAvg))
				sendCnt.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP).Set(float64(p.Statistics.SendCnt))
				sendPend.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP).Set(float64(p.Statistics.SendPending))
				sendMax.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP).Set(float64(p.Statistics.SendMax))
				sendOct.WithLabelValues(strconv.FormatInt(*t.Number, 10), t.OriginHost, p.DestinationHost, p.RemoteIP).Set(float64(p.Statistics.SendOct))
			}
		}
	}
	return reg
}
func WriteToFile(gatherer prometheus.Gatherer, filename string) error {
	err := prometheus.WriteToTextfile(filename, gatherer)
	if err != nil {
		return err
	}
	return nil
}
