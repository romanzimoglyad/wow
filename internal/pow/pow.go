package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	Ver      = 1
	zeroByte = 48
)

type HashBlock struct {
	Ver      int
	Bits     int
	Date     int64
	Resource string
	Rand     string
	Counter  int
}

func CalculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func CheckValidHash(hash string, zerosCount int) bool {
	if zerosCount > len(hash) {
		return false
	}
	for _, ch := range hash[:zerosCount] {
		if ch != zeroByte {
			return false
		}
	}
	return true
}

func (h *HashBlock) String() string {
	return fmt.Sprintf("%d:%d:%d:%s:%s:%d", h.Ver, h.Bits, h.Date, h.Resource, h.Rand, h.Counter)
}
func (h *HashBlock) DoWork(max int) error {
	for h.Counter < max {
		hashBlock := h.String()
		hash := CalculateHash(hashBlock)
		if CheckValidHash(hash, h.Bits) {
			return nil
		}
		h.Counter++
	}
	return fmt.Errorf("hash not found")
}
