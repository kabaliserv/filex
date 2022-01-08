package token

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"os"
	"strings"
)

type Manager struct {
	privateKey *rsa.PrivateKey
}

func New() Manager {
	p, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("failed to generate private key: %s\n", err)
		os.Exit(1)
	}
	t := Manager{
		privateKey: p,
	}

	return t
}

func (t *Manager) NewToken() jwt.Token {
	return jwt.New()
}

func (t *Manager) GetTokenFromRequest(r *http.Request) (jwt.Token, error) {
	return t.getTokenFromHeaders(r.Header)
}

func (t *Manager) GetTokenFromHeaders(h http.Header) (jwt.Token, error) {
	return t.getTokenFromHeaders(h)
}

func (t *Manager) Sign(token jwt.Token) ([]byte, error) {
	return jwt.Sign(token, jwa.RS256, t.privateKey)
}

func (t *Manager) getTokenFromHeaders(h http.Header) (jwt.Token, error) {
	authHeader := h.Get("Authorization")
	tokenSigned := extractToken(authHeader)
	if tokenSigned == "" {
		return nil, nil
	}

	return t.parsePayload(tokenSigned)
}
func (t *Manager) parsePayload(payload string) (jwt.Token, error) {
	return jwt.Parse(
		[]byte(payload),
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, &t.privateKey.PublicKey),
	)
}

func extractToken(v string) string {
	fmt.Printf("token %#v", v)
	return strings.TrimPrefix(v, "Bearer ")
}
