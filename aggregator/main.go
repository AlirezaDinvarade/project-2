package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	types "tolling/types"

	"google.golang.org/grpc"
)

func main() {
	httplistenAddress := flag.String("httplistenAddress", ":3000", "the listen address of the HTTP server")
	grpclistenAddress := flag.String("grpclistenAddress", ":3001", "the listen address of the GRPC server")

	flag.Parse()

	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	go makeGRPCTransport(*grpclistenAddress, svc)
	log.Fatal(makeHTTPTransport(*httplistenAddress, svc))
}

func makeGRPCTransport(listenAddress string, svc Aggregator) error {
	fmt.Println("GRPC transport runnig on port ", listenAddress)
	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return err
	}
	defer ln.Close()
	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewAggregatorGRPCServer(svc))
	return server.Serve(ln)

}

func makeHTTPTransport(listenAddress string, svc Aggregator) error {
	fmt.Println("HTTP transport runnig on port ", listenAddress)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	return http.ListenAndServe(listenAddress, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBU ID"})
			return
		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid OBU ID"})
			return
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
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
