package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	types "tolling/Types"
)

func main() {
	listenAddress := flag.String("listenAddress", ":3000", "the listen address of the HTTP server")
	flag.Parse()
	store := NewMemoryStore()
	var (
		svc = NewInvoiceAggregator(store)
	)
	makeHTTPTransport(*listenAddress, svc)
}

func makeHTTPTransport(listenAddress string, svc Aggregator) {
	fmt.Println("HTTP transport runnig on port ", listenAddress)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(listenAddress, nil)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
