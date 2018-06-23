package main

import (
	"context"
	"flag"
	"math"

	"github.com/Azure/azure-functions-go-worker/worker"
	log "github.com/Sirupsen/logrus"
)

var (
	flagDebug            bool
	host                 string
	port                 int
	workerID             string
	requestID            string
	grpcMaxMessageLength int
)

func init() {

	flag.BoolVar(&flagDebug, "debug", true, "enable verbose output")
	flag.StringVar(&host, "host", "127.0.0.1", "RPC Server Host")
	flag.IntVar(&port, "port", 0, "RPC Server Port")
	flag.StringVar(&workerID, "workerId", "", "RPC Server Worker ID")
	flag.StringVar(&requestID, "requestId", "", "Request ID")
	flag.IntVar(&grpcMaxMessageLength, "grpcMaxMessageLength", math.MaxInt32, "Max message length")

	flag.Parse()

	if flagDebug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	cfg := &worker.ClientConfig{
		Host:             host,
		Port:             port,
		WorkerID:         workerID,
		RequestID:        requestID,
		MaxMessageLength: grpcMaxMessageLength,
	}
	client := worker.NewClient(cfg)
	err := client.Connect()

	if err != nil {
		log.Fatalf("cannot create grpc connection: %v", err)
	}
	defer client.Disconnect()

	client.StartEventStream(context.Background())
}
