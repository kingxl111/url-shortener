package shortener

import (
	"crypto/sha256"
	"encoding/binary"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const ShortURLLength = 10

func encodeBase62(num uint64) string {
	var sb strings.Builder
	for i := 0; i < ShortURLLength; i++ {
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

func IsValidShortURL(s string) bool {
	for _, c := range s {
		if !isAllowedShortURLChar(c) {
			return false
		}
	}
	return true
}

func isAllowedShortURLChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}
