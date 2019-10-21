package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gitlab.com/NebulousLabs/Sia/build"
	sia "gitlab.com/NebulousLabs/Sia/node/api/client"
)

var (
	debug  bool
	module string

	log *logrus.Logger
)

//initSiaClient sets the values of the client so we can communicate with the
//Sia Daemon.
func findPassword() string {
	// Check environment variables
	apiPassword := os.Getenv("SIA_API_PASSWORD")
	if apiPassword != "" {
		log.Info("Using SIA_API_PASSWORD environment variable")
		return apiPassword
	}

	// Check .apipassword file
	var siaDir = build.DefaultSiaDir()
	pw, err := ioutil.ReadFile(build.APIPasswordFile(siaDir))
	if err != nil {
		log.Info("Could not read API password file:", err)
		return ""
	} else {
		return strings.TrimSpace(string(pw))
	}
}

// initLogger initializes the logger
func initLogger(debug bool) {
	log = logrus.New()

	// Define logger level
	if debug {
		log.SetLevel(logrus.DebugLevel)
		// Print out file names and line numbers
		log.SetReportCaller(true)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

}

// boolToFloat64 converts a bool to a float64
func boolToFloat64(b bool) float64 {
	if b {
		return float64(1)
	}
	return float64(0)
}

// startMonitor refreshes the Sia metrics periodically as defined by refreshRate
func startMonitor(refreshRate time.Duration, sc *sia.Client) {
	for range time.Tick(time.Minute * refreshRate) {
		updateMetrics(sc)
	}
}

// updateMetrics calls the various metric collection functions
func updateMetrics(sc *sia.Client) {

	log.Debug("Updating metrics for modules:", module)

	log.Debug("Updating Daemon Metrics")
	daemonMetrics(sc)

	if strings.Contains(module, "r") {
		log.Debug("Updating Renter Metrics")
		renterMetrics(sc)
		log.Debug("Updating hostdb Metrics")
		hostdbMetrics(sc)
	}

	if strings.Contains(module, "c") {
		log.Debug("Updating Consensus Metrics")
		consensusMetrics(sc)
	}

	if strings.Contains(module, "w") {
		log.Debug("Updating Wallet Metrics")
		walletMetrics(sc)
	}

	if strings.Contains(module, "g") {
		log.Debug("Updating Gateway Metrics")
		gatewayMetrics(sc)
	}

	if strings.Contains(module, "h") {
		log.Debug("Updating Host Metrics")
		hostMetrics(sc)
	}

	if strings.Contains(module, "m") {
		log.Info("Miner metrics are not implemented yet")
	}

	if strings.Contains(module, "t") {
		log.Info("Transactionpool metrics are not implemented yet")
	}

}

func main() {

	// Flags
	flag.BoolVar(&debug, "debug", false, "Enable debug mode. Warning: generates a lot of output.")
	address := flag.String("address", "127.0.0.1:9980", "Sia's API address")
	agent := flag.String("agent", "Sia-Agent", "Sia agent")
	refresh := flag.Int("refresh", 5, "Frequency to get Metrics from Sia (minutes)")
	port := flag.Int("port", 9983, "Port to serve Prometheus Metrics on")
	flag.StringVar(&module, "modules", "cghmrtw", "Sia Modules to monitor")
	flag.Parse()

	// Initialize the logger
	initLogger(debug)

	// Set the Sia Client connection information
	sc := sia.New(*address)
	sc.UserAgent = *agent
	sc.Password = findPassword()

	// Set the metrics initially before starting the monitor and HTTP server
	// If you don't do this all the metrics start with a "0" until they are set
	updateMetrics(sc)

	// start the metrics collector
	go startMonitor(time.Duration(*refresh), sc)

	// This section will start the HTTP server and expose
	// any metrics on the /metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to metrics at http://<your ip address>:", *port, "/metrics")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
