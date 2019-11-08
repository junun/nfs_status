package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"nfs_status/collector"
)

const (
	nfsMountPath 	= "/data/images/lighting"
)

var (
	// commandline arguments
	listenAddr  		= flag.String("web.listen-port", "9001", "An port to listen on for web interface and telemetry.")
	metricsPath 		= flag.String("web.telemetry-path", "/metrics", "A path under which to expose metrics.")
	metricsNamespace 	= flag.String("metric.namespace", "nfs", "Prometheus metrics namespace, as the prefix of metrics name")
	nfsPath       		= flag.String("nfs.storage-path", nfsMountPath, "Path to nfs storage volume.")

	num           int
)


func main() {
	flag.Parse()

	metrics 	:= collector.NewMetrics(*metricsNamespace, *nfsPath)
	registry 	:= prometheus.NewRegistry()
	registry.MustRegister(metrics)

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Printf("Starting Server at http://localhost:%s%s", *listenAddr, *metricsPath)
	if err := http.ListenAndServe(":"+*listenAddr, nil); err != nil {
		fmt.Printf("Error occur when start server %v", err)
	}
}