package auth

import (
	"crypto/ed25519"
	"encoding/base64"

	"golang.org/x/crypto/blake2b"
)

// GenerateKeys creates a new ed25519 key pair.
func GenerateKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(nil)
}

// CreateDigest creates a Blake2b hash of the given payload.
func CreateDigest(payload []byte) (string, error) {
	hash, err := blake2b.New512(nil)
	if err != nil {
		return "", err
	}
	_, err = hash.Write(payload)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

// SignRequest signs the digest with the given private key.
func SignRequest(privateKey ed25519.PrivateKey, digest string) (string, error) {
	signature := ed25519.Sign(privateKey, []byte(digest))
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySignature verifies the signature of the request.
func VerifySignature(publicKey ed25519.PublicKey, signature string, digest string) bool {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	return ed25519.Verify(publicKey, []byte(digest), sig)
}