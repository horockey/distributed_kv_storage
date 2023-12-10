package http_remote_storage

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/horockey/distributed_kv_storage/internal/adapter/gateway/remote_storage"
	"github.com/horockey/distributed_kv_storage/internal/controller/http_controller/dto"
)

var _ remote_storage.Gateway[map[string]any] = &httpRemoteStorage{}

type httpRemoteStorage struct {
	restClient *resty.Client
}

func New() *httpRemoteStorage {
	return &httpRemoteStorage{
		restClient: resty.New(),
	}
}

func (gw *httpRemoteStorage) Get(key string, node remote_storage.AppNode) (map[string]any, error) {
	resp, err := gw.restClient.R().
		SetPathParam("key", key).
		Get(fmt.Sprintf("http://%s/kv/{key}", node.Name))
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("got non-ok response (%s): %s", resp.Status(), resp.String())
	}

	kv := dto.KV{}
	if err := json.Unmarshal(resp.Body(), &kv); err != nil {
		return nil, fmt.Errorf("decoding response json: %w", err)
	}

	return kv.Value, nil
}

func (gw *httpRemoteStorage) Set(key string, val map[string]any, node remote_storage.AppNode) error {
	kv := dto.KV{
		Key:   key,
		Value: val,
	}

	resp, err := gw.restClient.R().
		SetBody(kv).
		Post(fmt.Sprintf("http://%s/kv", node.Name))
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("got non-ok response (%s): %s", resp.Status(), resp.String())
	}

	return nil
}
