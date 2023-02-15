package cipher

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!_-")

func randomSequence(n int) string {
    rand.Seed(time.Now().UnixNano())
    b := make([]rune, n)
    for i := range b {
        b[i] = chars[rand.Intn(len(chars))]
    }
    return string(b)
}

func GeneratePassword(password, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func GenerateCode() (string, error) {
	code, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	
	return strings.Replace(code.String(), "-", "", -1), nil
}