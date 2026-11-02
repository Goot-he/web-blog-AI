package utils

import (
	"gateway/internal/config"
	"net/http"
	"time"
)

const (
	refreshCookieName = "refresh"
	refreshCookiePath = "/api/v1"
)

func SetRefreshCookie(writer http.ResponseWriter, refreshToken string, cfg config.AuthConf) {
	http.SetCookie(writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    refreshToken,
		Path:     refreshCookiePath,
		MaxAge:   cfg.RefreshExpire,
		HttpOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func RefreshTokenFromCookie(request *http.Request) (string, error) {
	cookie, err := request.Cookie(refreshCookieName)
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", http.ErrNoCookie
	}
	return cookie.Value, nil
}

func ClearRefreshCookie(writer http.ResponseWriter, cfg config.AuthConf) {
	http.SetCookie(writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}
