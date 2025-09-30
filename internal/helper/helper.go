// internal/auth/response.go

package helper

import (
	"net/http"
	"encoding/json"
)


func WriteACK(w http.ResponseWriter, messageID string) {
	resp := map[string]any{
		"message": map[string]string{
			"ack":        "ACK",
			"message_id": messageID,
		},
		"status": "success",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}


func WriteNACK(w http.ResponseWriter, errMsg, code string) {
	resp := map[string]any{
		"message": map[string]string{
			"ack": "NACK",
		},
		"error": map[string]string{
			"code":    code,
			"message": errMsg,
		},
		"status": "failure",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest) // 🔴 Protocol/validation failure
	json.NewEncoder(w).Encode(resp)
}
