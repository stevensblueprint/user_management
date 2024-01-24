package models

import "sync"

type TokenPool struct {
	tokens map[string]*Token
	mu     sync.Mutex
}

func NewTokenPool() *TokenPool {
	return &TokenPool{
		tokens: make(map[string]*Token),
	}
}

func (p *TokenPool) Add(token *Token) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.tokens[token.Value] = token
}

func (p *TokenPool) Get(token string) *Token {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.tokens[token]
}

func (p *TokenPool) ValidateToken(tokenValue string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	token, exists := p.tokens[tokenValue]
	if !exists {
		return false
	}

	if token.IsValid() {
		return true
	}

	delete(p.tokens, tokenValue)
	return false
}

func (p *TokenPool) InvalidateToken(tokenValue string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.tokens, tokenValue)
}
