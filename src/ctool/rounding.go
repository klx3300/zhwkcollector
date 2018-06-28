package ctool

import (
	"crypto/aes"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// AESKeyRound may not be so safe.
func AESKeyRound(ky string) []byte {
	if len(ky) < 16 {
		buf := make([]byte, 16)
		copy(buf, []byte(ky))
		return buf
	} else if len(ky) < 24 {
		buf := make([]byte, 24)
		copy(buf, []byte(ky))
		return buf
	} else {
		buf := make([]byte, 32)
		copy(buf, []byte(ky)[:max(len(ky), 32)])
		return buf
	}
}

// AESIV generate the iv for AES encryption.
func AESIV(iv string) []byte {
	result := make([]byte, aes.BlockSize)
	for i := 0; i < aes.BlockSize; i++ {
		result[i] = iv[i%len(iv)]
	}
	return result
}
