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

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	invReq := types.GetInvoiceRequest{
		OBUID: int32(id),
	}
	b, err := json.Marshal(&invReq)
	if err != nil {
		return nil, err
	}
	endPoint := fmt.Sprintf("%s/%s?obu=%d", c.EndPoint, "invoice", id)
	req, err := http.NewRequest("POST", endPoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("the service responded with non 200 status code %d", resp.StatusCode)
	}
	var inv types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &inv, nil
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.EndPoint+"/aggregate", bytes.NewReader(b))
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
