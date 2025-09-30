package auth

import (
	"net/http"
    "log"
	"bytes"
	"io"
	"strings"
	"os"
	"encoding/json"
	// "ondc-poc/internal/database"
	"ondc-poc/internal/helper"
	"regexp"
	// "crypto/sha256"
	// "encoding/base64"

	
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
		pubKey, err := helper.GetPublicKeyByUKID(ukid)
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
		computedDigest := CreateDigest(bodyBytes)
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





