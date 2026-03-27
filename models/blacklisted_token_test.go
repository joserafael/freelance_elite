package models

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlacklistedTokenFields(t *testing.T) {
	expiresAt := time.Now().Add(time.Hour)
	token := BlacklistedToken{
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
		ExpiresAt: expiresAt,
	}

	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test", token.Token)
	assert.Equal(t, expiresAt, token.ExpiresAt)
}

func TestBlacklistedTokenGormTags(t *testing.T) {
	token := BlacklistedToken{}
	tokenType := reflect.TypeOf(token)

	tokenField, _ := tokenType.FieldByName("Token")
	assert.Contains(t, string(tokenField.Tag), `gorm:"type:varchar(500);uniqueIndex"`)

	expiresAtField, _ := tokenType.FieldByName("ExpiresAt")
	assert.Contains(t, string(expiresAtField.Tag), `gorm:"not null"`)
}

func TestBlacklistedTokenDefaultValues(t *testing.T) {
	token := BlacklistedToken{}

	assert.Empty(t, token.Token)
	assert.Equal(t, time.Time{}, token.ExpiresAt)
}
