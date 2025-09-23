package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ondc-poc/internal/models"
	"ondc-poc/internal/registry"
	"ondc-poc/internal/auth"
)

func LookupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req models.LookupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	results, err := registry.FindParticipants(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find participants: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func VlookupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "vlookup not implemented yet"}`))
}


func SignHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	signature, err := auth.SignDigest(req.PrivateKey, req.Digest)
	if err != nil {
		http.Error(w, "Failed to sign digest: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := models.SignResponse{Signature: signature}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Failed to write response:", err)
	}
}