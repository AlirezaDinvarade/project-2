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
	
	var (
		store = NewMemoryStore()
		svc = NewInvoiceAggregator(store)
		
	)
	svc = NewLogMiddleware(svc)
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
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
