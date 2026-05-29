package format_test

import (
	"testing"
	"time"

	"github.com/laschenkov67/docgen/format"
	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	tm := time.Date(2025, 3, 12, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "12.03.2025", format.Date(tm))
	assert.Equal(t, "12 марта 2025 г.", format.DateLong(tm))
	assert.Equal(t, "", format.Date(time.Time{}))
}

func TestValidINN(t *testing.T) {
	assert.True(t, format.ValidINN("7707083893")) // Сбербанк, реальный
	assert.False(t, format.ValidINN("1234567890"))
	assert.False(t, format.ValidINN("abc"))
}
