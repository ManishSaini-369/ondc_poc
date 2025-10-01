package helper

import (
	"fmt"
	"time"
    "encoding/base64"
	"crypto/aes"
	"crypto/cipher"
	"ondc-poc/internal/database"
	"golang.org/x/net/context"
	"ondc-poc/internal/models"

)

// Redis TTL for cached public keys
const publicKeyTTL = time.Hour

// GetPublicKeyByUKID fetches public key from Redis cache first, then DB

func GetPublicKeyByUKID(ukid string) (string, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("participant:ukid:%s", ukid)

	// Try Redis first
	if database.RDB != nil {
		cachedKey, err := database.RDB.Get(ctx, cacheKey).Result()
		if err == nil && cachedKey != "" {
			return cachedKey, nil // Base64 string from cache
		}
	}

	// Fallback: Query PostgreSQL via GORM
	var participant models.Participant
if err := database.DB.Table("ondc.participants").
	Select("signing_public_key").
	Where("ukid = ?", ukid).
	First(&participant).Error; err != nil {
	return "", fmt.Errorf("ukid not found or db error: %w", err)
}

	publicKeyBase64 := participant.SigningPublicKey

	// Store in Redis
	if database.RDB != nil {
		_ = database.RDB.Set(ctx, cacheKey, publicKeyBase64, publicKeyTTL).Err()
	}

	return publicKeyBase64, nil
}




func DecryptAES(encryptedBase64, keyString, ivString string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", err
	}

	key := []byte(keyString)
	iv := []byte(ivString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// remove PKCS7 padding
	padding := int(plaintext[len(plaintext)-1])
	if padding > len(plaintext) {
		return "", fmt.Errorf("invalid padding")
	}
	plaintext = plaintext[:len(plaintext)-padding]

	return string(plaintext), nil
}

