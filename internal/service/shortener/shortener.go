package shortener

import (
	"crypto/sha256"
	"encoding/binary"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const shortURLLength = 10

func encodeBase62(num uint64) string {
	var sb strings.Builder
	for i := 0; i < shortURLLength; i++ {
		sb.WriteByte(charset[num%uint64(len(charset))])
		num /= uint64(len(charset))
	}
	return sb.String()
}

func GenerateShortURL(originalURL string) string {
	hash := sha256.Sum256([]byte(originalURL))
	num := binary.BigEndian.Uint64(hash[:8])
	return encodeBase62(num)
}
