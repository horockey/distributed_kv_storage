package dto

type KV struct {
	Key   string         `json:"key"`
	Value map[string]any `json:"value"`
}
