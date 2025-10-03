package models

type VlookupRequest struct {
	SenderSubscriberID string `json:"sender_subscriber_id"`
	RequestID          string `json:"request_id"`
	Timestamp          string `json:"timestamp"`
	Signature          string `json:"signature"`
	SearchParameters   struct {
		Country      string `json:"country"`
		Domain       string `json:"domain"`
		Type         string `json:"type"`
		City         string `json:"city"`
		SubscriberID string `json:"subscriber_id"`
	} `json:"search_parameters"`
}
