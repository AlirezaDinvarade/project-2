package main

import "tolling/types"

type AggregatorGRPCServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewAggregatorGRPCServer(svc Aggregator) *AggregatorGRPCServer {
	return &AggregatorGRPCServer{
		svc: svc,
	}
}

func (s *AggregatorGRPCServer) AggregateDistance(req *types.AggregateRequest) error {
	distance := types.Distance{
		OBUID: int(req.OnuID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return  s.svc.AggregateDistance(distance)
}
