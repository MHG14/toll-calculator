package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/mhg14/toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type InvoiceHandler struct {
	client client.Client
}

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "the listen address of the http server")
	flag.Parse()

	aggregatorServiceAddr := flag.String("aggServiceAddr", "http://localhost:5000", "the listen address of the aggregator service")
	var (
		client     = client.NewHTTPClient(*aggregatorServiceAddr)
		invHandler = newInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeApiFunc(invHandler.handleGetInvoice))
	logrus.Infof("gateway http server running on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	inv, err := h.client.GetInvoice(context.Background(), 862870508)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func newInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeApiFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
