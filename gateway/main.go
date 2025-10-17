package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
	"tolling/aggregator/client"

	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listrenAddress := flag.String("listrenAddress", ":6000", "the listen address of gateway")
	aggregatorServiceAdd := flag.String("aggServiceAdd", "http://localhost:3000", "the listen address of aggregator")
	flag.Parse()
	var (
		client         = client.NewHTTPClient(*aggregatorServiceAdd)
		invoiceHandler = NewInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeAPIFunc(invoiceHandler.HandleGetInvoice))
	logrus.Infof("gateway HTTP server on port %s", *listrenAddress)
	log.Fatal(http.ListenAndServe(*listrenAddress, nil))

}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{client: c}
}

func (h *InvoiceHandler) HandleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	inv, err := h.client.GetInvoice(context.Background(), -867618446)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)

}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func (start time.Time)  {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri": r.RequestURI,
			}).Info("REQ :: ")
		}(time.Now())
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
