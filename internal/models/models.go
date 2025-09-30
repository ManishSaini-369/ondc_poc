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



// Search ---------flow 

// Context is part of ONDC protocol

type ONDCSearchRequest struct {
	Context Context `json:"context"`
	Message struct {
		Intent map[string]interface{} `json:"intent"`
	} `json:"message"`
}


type Context struct {
	TransactionID string `json:"transaction_id"`
	MessageID     string `json:"message_id"`
	BapID         string `json:"bap_id"`
	BapURI        string `json:"bap_uri"`
	Domain        string `json:"domain"`
	City          string `json:"city"`
	Timestamp     string `json:"timestamp"`
}


type AckResponse struct {
	Context Context   `json:"context"`
	Message AckMessage `json:"message"`
}



type Ack struct {
	Status string `json:"status"`
}



type AckMessage struct {
	Ack struct {
		Status string `json:"status"`
	} `json:"ack"`
}