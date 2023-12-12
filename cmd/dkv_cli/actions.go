package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

func get(ctx *cli.Context) error {
	endpoint := ctx.String("endpoint")
	key := ctx.String("key")

	resp, err := restyClient.R().
		SetPathParam("key", key).
		Get(fmt.Sprintf("http://%s/kv/{key}", endpoint))
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("gon non-ok response (%s): %s", resp.Status(), resp.String())
	}

	data := struct {
		Key   string
		Value map[string]any
	}{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return fmt.Errorf("unmarshalling response data: %w", err)
	}

	logger.Info().Msg("Got KV successfully:")
	fmt.Printf("Key: %s\nValue: %+v", data.Key, data.Value)

	return nil
}

func put(ctx *cli.Context) error {
	endpoint := ctx.String("endpoint")
	key := ctx.String("key")
	valueStr := ctx.String("value")

	value := map[string]any{}
	if err := json.Unmarshal([]byte(valueStr), &value); err != nil {
		return fmt.Errorf("value must be struct, unable to marshal: %w", err)
	}

	resp, err := restyClient.R().
		SetBody(map[string]any{
			"key":   key,
			"value": value,
		}).
		Post(fmt.Sprintf("http://%s/kv", endpoint))
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("gon non-ok response (%s): %s", resp.Status(), resp.String())
	}

	logger.Info().Msg("Put KV successfully")

	return nil
}
