package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	types "tolling/types"
)

type HTTPClient struct {
	EndPoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		EndPoint: endpoint,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.EndPoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responded with non 200 status code %d", resp.StatusCode)
	}
	return nil
}
