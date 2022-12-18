package models

type Secret struct {
	Key    string `json:"key,omitempty"`
	Store  string `json:"store,omitempty"`
	Secret string `json:"secret,omitempty"`
}
