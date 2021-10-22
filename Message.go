package Coral

type message struct {
	ID      string      `json:"id,omitempty"`
	Channel string      `json:"channel,omitempty"`
	Type    string      `json:"type"`
	Content interface{} `json:"content,omitempty"`
}
