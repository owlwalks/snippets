package snippets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func encryptSymmetric(hashedPassword []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(hashedPassword)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Seal(nil, nonce, plaintext, nil), nil
}
