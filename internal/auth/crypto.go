package auth

import (
	"crypto/ed25519"
	"encoding/base64"
    "golang.org/x/crypto/blake2b"
	// "net/http"
    "fmt"
    "log"
    "errors"
    "regexp"
	"bytes"
	"encoding/json"
	// "io"
    // "strings"
)


func NormalizeJSON(body []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.Compact(&buf, body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}



// GenerateKeys creates a new ed25519 key pair.
func GenerateKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(nil)
}

func CreateDigest(body []byte) (string, error) {
	// Step 1: Normalize JSON
	normalized, err := NormalizeJSON(body)
	if err != nil {
		return "", err
	}

	// Step 2: Generate BLAKE2b-256 hash
	hash := blake2b.Sum256(normalized)

	// Step 3: Return Base64 encoded hash
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}


// SignRequest signs the digest with the given private key.
func SignRequest(privateKey ed25519.PrivateKey, digestBytes []byte) (string, error) {
    signature := ed25519.Sign(privateKey, digestBytes)
    return base64.StdEncoding.EncodeToString(signature), nil
}

func VerifySignature(pubKeyBase64, signatureBase64, digestBase64 string) bool {
	// Decode public key

	fmt.Println("jjjjjjjjjxs ", pubKeyBase64)
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKeyBase64)
    fmt.Println("dfndfndjfd ", len(pubKeyBytes))
	if err != nil {
		log.Println("Invalid public key:", err)
		return false
	}

	// Decode signature
	sigBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		log.Println(" Invalid signature:", err)
		return false
	}

	// Decode digest
	digestBytes, err := base64.StdEncoding.DecodeString(digestBase64)
	if err != nil {
		log.Println(" Invalid digest:", err)
		return false
	}

	// Verify signature
	ok := ed25519.Verify(ed25519.PublicKey(pubKeyBytes), digestBytes, sigBytes)
	if !ok {
		log.Println(" Signature verification failed!")
	}
	return ok
}






func SignDigest(privateKeyBase64 string, digestBase64 string) (string, error) {
    privKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
    if err != nil {
        return "", err
    }

    digestBytes, err := base64.StdEncoding.DecodeString(digestBase64)
    if err != nil {
        return "", err
    }

    signature := ed25519.Sign(ed25519.PrivateKey(privKeyBytes), digestBytes)
    return base64.StdEncoding.EncodeToString(signature), nil
}






// func ExtractUKID(r *http.Request) (string, error) {
//     authHeader := r.Header.Get("Authorization")
//     if authHeader == "" {
//         return "", http.ErrNoCookie // just reusing a simple error
//     }

//     // Example: Authorization: Signature keyId="UKID001",algorithm="ed25519",signature="..."
//     parts := strings.Split(authHeader, ",")
//     for _, part := range parts {
//         if strings.Contains(part, "keyId=") {
//             // clean up: remove keyId= and quotes
//             keyID := strings.TrimSpace(strings.Split(part, "=")[1])
//             keyID = strings.Trim(keyID, `"`)
//             return keyID, nil
//         }
//     }

//     return "", http.ErrNoCookie
// }


func ExtractSignature(authHeader string) (string, error) {
	re := regexp.MustCompile(`signature="([^"]+)"`)
	matches := re.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", errors.New("signature not found")
	}
	return matches[1], nil
}


func ExtractUKID(authHeader string) (string, error) {
	re := regexp.MustCompile(`keyId="([^"]+)"`)
	matches := re.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", errors.New("keyId not found")
	}
	return matches[1], nil
}