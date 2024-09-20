package entity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type UserPersonalData struct {
	Firstname *string `json:"firstname,omitempty" gorm:"column:firstname_encrypted"`
	Lastname  *string `json:"lastname,omitempty" gorm:"column:lastname_encrypted"`
}

func (u *UserPersonalData) HasEmptyValues() bool {
	return u.Firstname == nil ||
		u.Lastname == nil
}

func (u *UserPersonalData) UpdateNonEmpty(d UserPersonalData) {
	if d.Firstname != nil {
		u.Firstname = d.Firstname
	}

	if d.Lastname != nil {
		u.Lastname = d.Lastname
	}

}

func (u *UserPersonalData) Encrypt(baseKey string, userKey *string) error {
	if userKey == nil {
		return errors.New("user encryption key not found")
	}
	return u.transformPersonalData(baseKey, userKey, encryptString)
}

func (u *UserPersonalData) Decrypt(baseKey string, userKey *string) error {
	if userKey == nil {
		return errors.New("user encryption key not found")
	}
	return u.transformPersonalData(baseKey, userKey, decryptString)
}

func (u *UserPersonalData) transformPersonalData(baseKey string, userKey *string,
	transformFunc func(key []byte, text string) (string, error)) error {
	key := u.generateKey(baseKey, userKey)

	if u.Firstname != nil {
		transformed, err := transformFunc(key, *u.Firstname)
		if err != nil {
			return err
		}
		u.Firstname = &transformed
	}

	if u.Lastname != nil {
		transformed, err := transformFunc(key, *u.Lastname)
		if err != nil {
			return err
		}
		u.Lastname = &transformed
	}

	return nil
}

func encryptString(key []byte, text string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	b := base64.StdEncoding.EncodeToString([]byte(text))
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptString(key []byte, text string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(b) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := b[:aes.BlockSize]
	b = b[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(b, b)
	data, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (u *UserPersonalData) generateKey(baseKey string, userKey *string) []byte {
	h := sha256.New()
	h.Write([]byte(baseKey + *userKey))
	return h.Sum(nil)
}

func (u *UserPersonalData) GetFullName() string {
	return fmt.Sprintf("%s %s", u.getFirstName(), u.getLastName())
}
func (u *UserPersonalData) getFirstName() string {
	if u.Firstname == nil {
		return ""
	}
	return *u.Firstname
}

func (u *UserPersonalData) getLastName() string {
	if u.Lastname == nil {
		return ""
	}
	return *u.Lastname
}
