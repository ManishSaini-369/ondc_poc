package auth

import (
	"net/http"
    "log"
	"bytes"
	"io"
	"strings"
	"ondc-poc/internal/database"
	"ondc-poc/internal/helper"
	
)


func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1️⃣ Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		log.Println("Step 1: Body read successfully")

		// Reset body for next handler
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// 2️⃣ Compute Digest (SHA-256, then Base64)
		computedDigest := CreateDigest(body)
		log.Println("Computed Digest:", computedDigest)

		// 3️⃣ Get Digest from client header
		clientDigest := strings.TrimPrefix(r.Header.Get("Digest"), "SHA-256=")
		if clientDigest == "" {
			http.Error(w, "Digest header missing", http.StatusUnauthorized)
			return
		}

		// 4️⃣ Compare Digests
		if clientDigest != computedDigest {
			http.Error(w, "Digest mismatch", http.StatusUnauthorized)
			return
		}
		log.Println("Step 2: Digest verified ")

		// 5️⃣ Extract ukid (keyId) from Authorization header
		authHeader := r.Header.Get("Authorization")
		ukid, err := ExtractUKID(authHeader)
		if err != nil {
			http.Error(w, "Missing or invalid keyId in Authorization header", http.StatusUnauthorized)
			return
		}
		log.Println("Extracted ukid:", ukid)

		// 6️⃣ Extract Signature from Authorization header
		signature, err := ExtractSignature(authHeader)
		if err != nil {
			http.Error(w, "Signature not found in Authorization header", http.StatusUnauthorized)
			return
		}
		log.Println("Signature from client:", signature)

		// 7️ Fetch public key from participants table using ukid
		pubKey, err := helper.GetPublicKeyByUKID(database.DB, ukid)
		if err != nil {
			http.Error(w, "Unable to fetch public key for ukid", http.StatusUnauthorized)
			return
		}
		log.Println("Fetched public key for ukid:", ukid)

		// 8️⃣ Verify Signature
		if !VerifySignature(pubKey, signature, computedDigest) {
			http.Error(w, "Signature verification failed", http.StatusUnauthorized)
			return
		}
		log.Println("Step 3: Signature verified ")

		// All checks passed
		next.ServeHTTP(w, r)
	}
}



