package models

type Chat struct {
	Id int64 `json:"id,omitempty"`

	Input string `json:"input,omitempty"`

	Output string `json:"output,omitempty"`

	Result string `json:"result,omitempty"`
}
