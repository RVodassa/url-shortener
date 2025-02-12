package random_test

import (
	"fmt"
	"testing"

	"github.com/RVodassa/url-shortener/internal/lib/random"
	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	r := random.New()

	tests := []struct {
		name        string
		length      int
		expectedErr error
	}{
		{
			name:        "успешная генерация строки длиной 10",
			length:      10,
			expectedErr: nil,
		},
		{
			name:        "успешная генерация строки длиной 256",
			length:      256,
			expectedErr: nil,
		},
		{
			name:        "успешная генерация строки длиной 1",
			length:      1,
			expectedErr: nil,
		},
		{
			name:        "ошибка при длине 0",
			length:      0,
			expectedErr: fmt.Errorf("random.RandomString: length=%d. %w", 0, random.ErrShortLength),
		},
		{
			name:        "ошибка при отрицательной длине",
			length:      -5,
			expectedErr: fmt.Errorf("random.RandomString: length=%d. %w", -5, random.ErrShortLength),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := r.RandomString(tt.length)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.length, len(result))
			}
		})
	}
}
