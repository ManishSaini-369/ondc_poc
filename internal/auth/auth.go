package auth

import (
	"net/http"
    "log"
	"bytes"
	"io"
	"strings"
	"os"
	"encoding/json"
	"ondc-poc/internal/database"
	"ondc-poc/internal/helper"
	
)


func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1️⃣ Read encrypted request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}
		log.Println("Step 1: Body read successfully (Encrypted)")

		// 2️⃣ Get AES Key & IV from env
		key := os.Getenv("AES_KEY") // must be 32 chars
		iv := os.Getenv("AES_IV")   // must be 16 chars
		if len(key) != 32 || len(iv) != 16 {
			http.Error(w, "Invalid AES key or IV length", http.StatusInternalServerError)
			return
		}

		// 3️⃣ Parse {"payload":"..."}
		var req struct {
			Payload string `json:"payload"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// 4️⃣ Decrypt payload
		decrypted, err := helper.DecryptAES(req.Payload, key, iv)
		if err != nil {
			http.Error(w, "Decryption failed", http.StatusBadRequest)
			return
		}
		log.Println("Step 2: Body decrypted:", decrypted)

		// Reset body for downstream
		r.Body = io.NopCloser(bytes.NewBuffer([]byte(decrypted)))

		// 5️⃣ Compute Digest from decrypted JSON
		computedDigest := CreateDigest([]byte(decrypted))
		log.Println("Computed Digest:", computedDigest)

		// 6️⃣ Compare with client Digest header
		clientDigest := strings.TrimPrefix(r.Header.Get("Digest"), "SHA-256=")
		if clientDigest == "" {
			http.Error(w, "Digest header missing", http.StatusUnauthorized)
			return
		}
		if clientDigest != computedDigest {
			http.Error(w, "Digest mismatch", http.StatusUnauthorized)
			return
		}
		log.Println("Step 3: Digest verified")

		// 7️⃣ Extract ukid
		authHeader := r.Header.Get("Authorization")
		ukid, err := ExtractUKID(authHeader)
		if err != nil {
			http.Error(w, "Missing or invalid keyId in Authorization header", http.StatusUnauthorized)
			return
		}
		log.Println("Extracted ukid:", ukid)

		// 8️⃣ Extract Signature
		signature, err := ExtractSignature(authHeader)
		if err != nil {
			http.Error(w, "Signature not found in Authorization header", http.StatusUnauthorized)
			return
		}
		log.Println("Signature from client:", signature)

		// 9️⃣ Fetch public key for ukid
		pubKey, err := helper.GetPublicKeyByUKID(database.DB, ukid)
		if err != nil {
			http.Error(w, "Unable to fetch public key for ukid", http.StatusUnauthorized)
			return
		}
		log.Println("Fetched public key for ukid:", ukid)

		// 🔟 Verify Signature
		if !VerifySignature(pubKey, signature, computedDigest) {
			http.Error(w, "Signature verification failed", http.StatusUnauthorized)
			return
		}
		log.Println("Step 4: Signature verified")

		// ✅ All good → pass decrypted body to handler
		next.ServeHTTP(w, r)
	}
}





