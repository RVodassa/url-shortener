package service

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// NewRandomString генерирует случайную строку заданной длины.
func NewRandomString(length int) string {
	if length <= 0 {
		return ""
	}

	rand.NewSource(time.Now().UnixNano()) // Инициализация генератора случайных чисел
	result := make([]byte, length)        // Срез байтов нужной длины

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))] // Генерируем случайный символ из charset
	}

	return string(result)
}
