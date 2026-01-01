package utils

type BraviaRequest struct {
	Method  string `json:"method"`
	ID      int    `json:"id"`
	Params  []any  `json:"params"`
	Version string `json:"version"`
}
