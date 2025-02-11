package random

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var ErrShortLength = errors.New("short length")

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Random struct {
	rand *rand.Rand
	mu   sync.Mutex
}

func New() *Random {
	return &Random{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// RandomString генерирует случайную строку заданной длины.
func (r *Random) RandomString(length int) (string, error) {
	const op = "random.RandomString"

	r.mu.Lock()
	defer r.mu.Unlock()

	if length <= 0 {
		return "", fmt.Errorf("%s: %w: length=%d", op, ErrShortLength, length)
	}

	alias := make([]byte, length)

	for i := range alias {
		alias[i] = charset[r.rand.Intn(len(charset))] // Генерируем случайный символ из charset
	}
	return string(alias), nil
}
