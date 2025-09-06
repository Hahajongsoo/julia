package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

func NewSID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func SignSID(sid string, secret []byte) string {
	m := hmac.New(sha256.New, []byte(secret))
	m.Write([]byte(sid))
	return base64.URLEncoding.EncodeToString(m.Sum(nil))
}

func VerifySID(sid, sig string, secret []byte) bool {
	want := SignSID(sid, secret)
	return hmac.Equal([]byte(want), []byte(sig))
}

func EncodeCookie(sid string, secret []byte) string {
	return sid + "|" + SignSID(sid, secret)
}

func DecodeCookie(cookie string) (sid, sig string, err error) {
	parts := strings.Split(cookie, "|")
	if len(parts) != 2 {
		return "", "", errors.New("invalid cookie format")
	}
	return parts[0], parts[1], nil
}

type Cookie struct {
	Name     string
	Path     string
	MaxAge   int
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
	Domain   string
}

func SetCookie(w http.ResponseWriter, value string, c Cookie) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.Name,
		Value:    value,
		Path:     c.Path,
		MaxAge:   c.MaxAge,
		HttpOnly: c.HttpOnly,
		Secure:   c.Secure,
		SameSite: c.SameSite,
		Domain:   c.Domain,
	})
}

func ClearCookie(w http.ResponseWriter, c Cookie) {
	http.SetCookie(w, &http.Cookie{
		Name:     c.Name,
		Value:    "",
		Path:     c.Path,
		MaxAge:   0,
		HttpOnly: c.HttpOnly,
		Secure:   c.Secure,
		SameSite: c.SameSite,
		Domain:   c.Domain,
	})
}
