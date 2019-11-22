package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func getKey(passphrase string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(passphrase))
	return hasher.Sum(nil)
}

func encrypt(key []byte, text string) (encodedText string, err error) {
	plainText := []byte(text)

	keyHash := getKey(string(key))

	block, err := aes.NewCipher(keyHash)
	if err != nil {
		return
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	encodedText = base64.URLEncoding.EncodeToString(cipherText)
	return
}

func decrypt(key []byte, encryptedText string) (text string, err error) {
	cipherText, err := base64.URLEncoding.DecodeString(encryptedText)
	if err != nil {
		return
	}

	keyHash := getKey(string(key))

	block, err := aes.NewCipher(keyHash)
	if err != nil {
		return
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	text = string(cipherText)
	return
}
