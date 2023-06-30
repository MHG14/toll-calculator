package main

import (
	"net"
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/mhg14/toll-calculator/gokit-aggregator/aggsvc/aggendpoint"
	"github.com/mhg14/toll-calculator/gokit-aggregator/aggsvc/aggservice"
	"github.com/mhg14/toll-calculator/gokit-aggregator/aggsvc/aggtransport"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	service := aggservice.New(logger)
	endpoints := aggendpoint.New(service, logger)
	httpHandler := aggtransport.NewHTTPHandler(endpoints, logger)

	// The HTTP listener mounts the Go kit HTTP handler we created.
	httpListener, err := net.Listen("tcp", ":5000")
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}
	logger.Log("transport", "HTTP", "addr", ":5000")
	err = http.Serve(httpListener, httpHandler)
	if err != nil {
		panic(err)
	}
}
