package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mhg14/toll-calculator/aggregator/client"
	"github.com/mhg14/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		store          = NewMemoryStore()
		svc            = NewInvoiceAggregator(store)
		httpListenAddr = os.Getenv("AGG_HTTP_ENDPOINT")
		grpcListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)
	go makeGRPCTransport(grpcListenAddr, svc)

	c, err := client.NewGRPCClient(grpcListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	if err = c.Aggregate(context.Background(), &types.AggregateRequest{
		OBUID: 4234324,
		Value: 23.44,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
	log.Fatal(makeHTTPTransport(httpListenAddr, svc))
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport running on port", listenAddr)

	// make a TCP listener
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	// make a gRPC native server
	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(listener)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	aggMetricsHandler := newHTTPMetricHandler("aggregate")
	invoiceMetricsHandler := newHTTPMetricHandler("invoice")

	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/aggregate", aggMetricsHandler.instrument(handleAggregate(svc)))
	http.HandleFunc("/invoice", invoiceMetricsHandler.instrument(handleGetInvoice(svc)))
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(listenAddr, nil)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
