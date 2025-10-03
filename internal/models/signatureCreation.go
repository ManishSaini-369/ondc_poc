package models




type SignRequest struct {
    PrivateKey string `json:"private_key"`
    Data       string `json:"data"`
}

// models/SignResponse.go
type SignResponse struct {
    Signature string `json:"signature"`
}