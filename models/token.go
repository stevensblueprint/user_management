package models

import (
	"time"
)

type Token struct {
	Value     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewToken(value string, duration time.Duration) *Token {
	return &Token{
		Value:     value,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
}

func (t *Token) IsValid() bool {
	return t.ExpiresAt.After(time.Now())
}
