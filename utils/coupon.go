package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"strconv"
	"time"
)

func GenerateCouponCode(id uint) (*string, error) {
	userID := strconv.FormatUint(uint64(id), 10)
	// Generate a random part of the code
	randomBytes := make([]byte, 4) // 4 bytes = 8 hex characters
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}
	randomPart := hex.EncodeToString(randomBytes)

	// Create a hash of the userID
	h := fnv.New32a()
	h.Write([]byte(userID))
	userHash := h.Sum32()

	// Combine the user hash, random part, and timestamp to ensure uniqueness
	timestamp := time.Now().Unix()
	couponCode := fmt.Sprintf("%x-%s-%x", userHash, randomPart, timestamp)

	return &couponCode, nil
}
