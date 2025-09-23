package models

import "time"

// Participant represents a single entity in the ONDC registry.
type Participant struct {
	SubscriberID      string    `json:"subscriber_id"`
	Status            string    `json:"status"`
	UkID              string    `json:"ukId"`
	SubscriberURL     string    `json:"subscriber_url"`
	Country           string    `json:"country"`
	Domain            string    `json:"domain"`
	ValidFrom         time.Time `json:"valid_from"`
	ValidUntil        time.Time `json:"valid_until"`
	Type              string    `json:"type"`
	SigningPublicKey  string    `json:"signing_public_key"`
	EncrPublicKey     string    `json:"encr_public_key"`
	Created           time.Time `json:"created"`
	Updated           time.Time `json:"updated"`
	BrID              string    `json:"br_id"`
	City              string    `json:"city"`
}

// LookupRequest is the request body for the /lookup endpoint.
type LookupRequest struct {
	SubscriberID string `json:"subscriber_id,omitempty"`
	Domain       string `json:"domain,omitempty"`
	UkID         string `json:"ukId,omitempty"`
	Country      string `json:"country,omitempty"`
	City         string `json:"city,omitempty"`
	Type         string `json:"type,omitempty"`
}

// VLookupRequest is the request body for the /vlookup endpoint.
type VLookupRequest struct {
	SenderSubscriberID string `json:"sender_subscriber_id"`
	RequestID          string `json:"request_id"`
	Timestamp          string `json:"timestamp"`
	SearchParameters   struct {
		Domain       string `json:"domain"`
		SubscriberID string `json:"subscriber_id"`
		Country      string `json:"country"`
		Type         string `json:"type"`
		City         string `json:"city"`
	} `json:"search_parameters"`
	Signature string `json:"signature"`
}

type SignRequest struct {
	PrivateKey string `json:"privateKey"` // base64 private key
	Digest     string `json:"digest"`     // string digest to sign
}

type SignResponse struct {
	Signature string `json:"signature"`
}
