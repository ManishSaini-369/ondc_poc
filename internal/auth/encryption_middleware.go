package auth

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"log"
)

// A hardcoded 32-byte key for AES-256. 
// For production, use a secure key management system.
var encryptionKey = []byte("12345678901234567890123456789012")

// responseWriter is a wrapper for http.ResponseWriter to capture the response.
type responseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.body.Write(b)
}

// EncryptionMiddleware handles decryption of requests and encryption of responses.
func EncryptionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("FATAL: Panic recovered in middleware: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Decrypt the request body
		decryptedBody, err := Decrypt(string(body), encryptionKey)
		if err != nil {
			http.Error(w, "Failed to decrypt request body", http.StatusBadRequest)
			return
		}

		// Replace the request body with the decrypted data
		r.Body = ioutil.NopCloser(bytes.NewBuffer(decryptedBody))

		// Create a response writer to capture the response
		rpw := &responseWriter{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}

		// Call the next handler
		next.ServeHTTP(rpw, r)

		// Encrypt the response
		encryptedResponse, err := Encrypt(rpw.body.Bytes(), encryptionKey)
		if err != nil {
			http.Error(w, "Failed to encrypt response", http.StatusInternalServerError)
			return
		}

		// Write the encrypted response
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(encryptedResponse))
	}
}
