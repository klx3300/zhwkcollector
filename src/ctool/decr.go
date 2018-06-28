package ctool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func fulfill(src []byte) []byte {
	// calculate how much we need to append..
	theappend := aes.BlockSize - len(src)%aes.BlockSize
	empty := make([]byte, theappend-1)
	return append(append(src, empty...), byte(theappend))
}

func deflate(src []byte) []byte {
	rmcont := int(src[len(src)-1])
	if rmcont > len(src) {
		// no error reported.
		return make([]byte, 0)
	}
	return src[:len(src)-rmcont]
}

// AESEncrypt .
func AESEncrypt(src []byte, key []byte, iv []byte) (string, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	rsrc := fulfill(src)
	buffer := make([]byte, len(rsrc))
	cfbencrypter := cipher.NewCFBEncrypter(blk, iv)
	cfbencrypter.XORKeyStream(buffer, rsrc)
	return base64.URLEncoding.EncodeToString(buffer), nil
}

func emptyrd() *bytes.Reader {
	return bytes.NewReader(make([]byte, 0))
}

// AESDecrypt .
func AESDecrypt(ciph string, key []byte, iv []byte) (*bytes.Reader, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return emptyrd(), err
	}
	// where are you from ..?
	rciph, err := base64.URLEncoding.DecodeString(ciph)
	if err != nil {
		return emptyrd(), err
	}

	if (len(rciph) % aes.BlockSize) != 0 {
		return emptyrd(), fmt.Errorf("unmatched cipher length and aes blksize, %d %% %d != 0", len(rciph), aes.BlockSize)
	}
	buffer := make([]byte, len(rciph))
	cfbdecrypter := cipher.NewCFBDecrypter(blk, iv)
	cfbdecrypter.XORKeyStream(buffer, rciph)
	defed := deflate(buffer)
	return bytes.NewReader(defed), nil
}
