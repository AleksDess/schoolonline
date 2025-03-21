package crypto

import (
	crypto_rand "crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/oklog/ulid/v2"
)

func GetUlid() string {
	// Инициализация источника случайных чисел с использованием криптографически безопасного генератора
	entropy := ulid.Monotonic(crypto_rand.Reader, 0)

	// Генерация нового ULID
	ulidCode := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

	return ulidCode.String()
}

const (
	digits       = "0123456789"
	lowerLetters = "abcdefghijklmnopqrstuvwxyz"
	upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	allChars     = digits + lowerLetters + upperLetters
	passwordLen  = 10
)

func GeneratePassword() string {
	for {
		pass := generatePassword()
		if pass == "" {
			time.After(time.Second)
			continue
		}
		return pass
	}
}

func generatePassword() string {
	var password = make([]byte, passwordLen)

	mandatoryChars := []string{digits, lowerLetters, upperLetters}
	for i, chars := range mandatoryChars {
		char, err := randomChar(chars)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		password[i] = char
	}

	for i := len(mandatoryChars); i < passwordLen; i++ {
		char, err := randomChar(allChars)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		password[i] = char
	}

	shuffle(password)

	return string(password)
}

func randomChar(chars string) (byte, error) {
	max := big.NewInt(int64(len(chars)))
	num, err := crypto_rand.Int(crypto_rand.Reader, max)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return chars[num.Int64()], nil
}

func shuffle(data []byte) {
	for i := range data {
		j, _ := crypto_rand.Int(crypto_rand.Reader, big.NewInt(int64(i+1)))
		data[i], data[j.Int64()] = data[j.Int64()], data[i]
	}
}
