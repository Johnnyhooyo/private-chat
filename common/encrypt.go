package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

type Encrypt interface {
	Encrypt(data []byte) []byte
	Decrypt(data []byte) []byte
}

type aesEncryptor struct {
	key []byte
}

func NewAesEncryptor(key []byte) Encrypt {
	return &aesEncryptor{key: key}
}

func GenerateAESKey(passphrase string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(passphrase))
	key := hasher.Sum(nil)
	return key
}

func (a *aesEncryptor) Encrypt(plaintext []byte) []byte {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		fmt.Println("err NewCipher" + err.Error())
		return nil
	}

	// 使用AES的CTR模式进行加密
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = rand.Read(iv); err != nil {
		fmt.Println("err rand.Read" + err.Error())
		return nil
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext
}

func (a *aesEncryptor) Decrypt(ciphertext []byte) []byte {
	// 使用相同的密钥和IV进行解密
	block, err := aes.NewCipher(a.key)
	if err != nil {
		fmt.Println("err NewCipher" + err.Error())
		return nil
	}

	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	iv := ciphertext[:aes.BlockSize]
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])

	return plaintext
}
