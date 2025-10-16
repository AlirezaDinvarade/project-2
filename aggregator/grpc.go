package main

import (
	"context"
	"tolling/types"
)

type AggregatorGRPCServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewAggregatorGRPCServer(svc Aggregator) *AggregatorGRPCServer {
	return &AggregatorGRPCServer{
		svc: svc,
	}
}

func (s *AggregatorGRPCServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.OBUID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(distance)
}
