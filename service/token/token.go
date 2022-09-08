package token

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Manager struct {
	privateKey *rsa.PrivateKey
}

func New() Manager {
	p, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Panicf("failed to generate private key: %s\n", err)
	}

	return Manager{
		privateKey: p,
	}
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

func (t *Manager) FromString(token string) (jwt.Token, error) {
	return t.parsePayload(token)
}

func (t *Manager) parsePayload(payload string) (jwt.Token, error) {
	return jwt.Parse(
		[]byte(payload),
		jwt.WithValidate(true),
		jwt.WithVerify(jwa.RS256, &t.privateKey.PublicKey),
	)
}

func extractToken(v string) string {
	return strings.TrimPrefix(v, "Bearer ")
}
