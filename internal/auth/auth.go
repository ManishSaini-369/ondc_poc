package auth

import (
	"net/http"
    "log"
	"bytes"
	"fmt"

	"io"
	"strings"
	// "os"
	"encoding/json"
	// "ondc-poc/internal/database"
	"ondc-poc/internal/helper"
	"regexp"
	// "crypto/sha256"
	// "encoding/base64"

	
)


func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1️⃣ Read raw request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		log.Println("Step 1: Body read successfully")

		// Reset body for downstream (important!)
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// 2️⃣ Compute Digest from body
		computedDigest, err := CreateDigest(body)
		if err != nil {
			http.Error(w, "Failed to compute digest", http.StatusInternalServerError)
			return
		}
		log.Println("Computed Digest:", computedDigest)

		// 3️⃣ Compare with client Digest header
		clientDigest := strings.TrimPrefix(r.Header.Get("Digest"), "SHA-256=")
		if clientDigest == "" {
			http.Error(w, "Digest header missing", http.StatusUnauthorized)
			return
		}
		if clientDigest != computedDigest {
			http.Error(w, "Digest mismatch", http.StatusUnauthorized)
			return
		}
		log.Println("Step 2: Digest verified")

		// 4️⃣ Extract ukid (keyId) from Authorization header
		authHeader := r.Header.Get("Authorization")
		ukid, err := ExtractUKID(authHeader)
		if err != nil {
			http.Error(w, "Missing or invalid keyId in Authorization header", http.StatusUnauthorized)
			return
		}
		log.Println("Extracted ukid:", ukid)

		// 5️⃣ Extract Signature
		signature, err := ExtractSignature(authHeader)
		if err != nil {
			http.Error(w, "Signature not found in Authorization header", http.StatusUnauthorized)
			return
		}
		log.Println("Signature from client:", signature)

		// 6️⃣ Fetch public key for ukid
		pubKey, err := helper.GetPublicKeyByUKID(ukid)
		if err != nil {
			http.Error(w, "Unable to fetch public key for ukid", http.StatusUnauthorized)
			return
		}
		log.Println("Fetched public key for ukid:", ukid)

		// 7️⃣ Verify Signature
		if !VerifySignature(pubKey, signature, computedDigest) {
			http.Error(w, "Signature verification failed", http.StatusUnauthorized)
			return
		}
		log.Println("Step 3: Signature verified")

		// ✅ All good → pass plain body to handler
		next.ServeHTTP(w, r)
	}
}




func AuthMiddlewareSearch(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1️⃣ Read Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.WriteNACK(w, "Missing Authorization header", "30001")
			return
		}

		// 2️⃣ Extract signature
		signature, err := ExtractSignature(authHeader)
		if err != nil {
			helper.WriteNACK(w, "Invalid Authorization header (no signature)", "30001")
			return
		}
		log.Println("Signature from client:", signature)

		// 3️⃣ Extract ukid (subscriber_id)
		re := regexp.MustCompile(`keyId="([^"]+)"`)
		matches := re.FindStringSubmatch(authHeader)
		if len(matches) != 2 {
			helper.WriteNACK(w, "Invalid keyId format", "30001")
			return
		}
		ukid := matches[1]
		log.Println("Extracted ukid:", ukid)

		// 4️⃣ Fetch public key
		pubKeyBase64, err := helper.GetPublicKeyByUKID(ukid)
		if err != nil {
			helper.WriteNACK(w, "Failed to fetch public key", "30002")
			return
		}
		log.Println("Fetched public key for ukid:", ukid)

		// 5️⃣ Read request body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			helper.WriteNACK(w, "Failed to read request body", "30003")
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // reset for next handler
		log.Println("Step 1: Body read successfully")

		// 6️⃣ Compute Digest (BLAKE2b-256)
		computedDigest, err := CreateDigest(bodyBytes)
		if err != nil {
			helper.WriteNACK(w, "Failed to compute digest", "30003")
			return
		}
		log.Println("Computed Digest:", computedDigest)

		// 7️⃣ Compare with Digest header
		clientDigest := strings.TrimPrefix(r.Header.Get("Digest"), "BLAKE-256=")
		if clientDigest == "" {
			helper.WriteNACK(w, "Digest header missing", "30001")
			return
		}
		if clientDigest != computedDigest {
			helper.WriteNACK(w, "Digest mismatch", "30004")
			return
		}
		log.Println("Step 2: Digest verified")

		// 8️⃣ Verify signature
		if !VerifySignature(pubKeyBase64, signature, computedDigest) {
			helper.WriteNACK(w, "Signature verification failed", "30005")
			return
		}
		log.Println("Step 3: Signature verified")

		// ✅ All good → pass request to handler
		next.ServeHTTP(w, r)
	}
}








// handles the /vlookup endpoint with signature verification





func AuthMiddlewareVlookup(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// 2) Parse only what we need from request body
		var payload struct {
			SenderSubscriberID string `json:"sender_subscriber_id"`
			Signature          string `json:"signature"`
			SearchParameters   struct {
				Country      string `json:"country"`
				Domain       string `json:"domain"`
				Type         string `json:"type"`
				City         string `json:"city"`
				SubscriberID string `json:"subscriber_id"`
			} `json:"search_parameters"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		// check required fields
		if payload.SenderSubscriberID == "" || payload.Signature == "" {
			http.Error(w, "Missing sender_subscriber_id or signature", http.StatusBadRequest)
			return
		}

		sp := payload.SearchParameters
		if sp.Country == "" || sp.Domain == "" || sp.Type == "" || sp.City == "" || sp.SubscriberID == "" {
			http.Error(w, "Missing search_parameters fields", http.StatusBadRequest)
			return
		}

		// 3) Build signing string
		signingString := fmt.Sprintf("%s|%s|%s|%s|%s",
			sp.Country, sp.Domain, sp.Type, sp.City, sp.SubscriberID,
		)

		log.Println("Signing string:", signingString)

		// 4) Fetch public key using sender_subscriber_id
		pubKey, err := helper.GetPublicKeyBySubscriberID(payload.SenderSubscriberID)
		if err != nil {
			http.Error(w, "Unable to fetch public key for sender_subscriber_id", http.StatusUnauthorized)
			return
		}

		log.Println("Fetched subscribver_id:", payload.SenderSubscriberID)
		log.Println("Fetched public key for sender_subscriber_id:", pubKey)

		// 5) Verify signature
		if !VerifySignatureString(pubKey, payload.Signature, signingString) {
			http.Error(w, "Signature verification failed", http.StatusUnauthorized)
			return
		}
		log.Println("Signature verified for:", payload.SenderSubscriberID)

		// 6) Pass request forward
		next.ServeHTTP(w, r)
	}
}





