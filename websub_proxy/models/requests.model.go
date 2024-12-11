package models

type SubscribeRequest struct {
	Topic    string `json:"topic"`
	Hub      string `json:"hub"`
	Callback string `json:"callback"`
}

// type PublishRequest struct {
