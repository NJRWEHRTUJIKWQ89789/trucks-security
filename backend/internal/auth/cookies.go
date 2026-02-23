package auth

import (
	"net/http"
	"time"
)

const (
	accessCookieName  = "cargomax_token"
	refreshCookieName = "cargomax_refresh"
)

// SetAuthCookies writes the access and refresh tokens as HttpOnly, SameSite=Strict
// cookies on the response.
func SetAuthCookies(w http.ResponseWriter, accessToken, refreshToken, domain string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     accessCookieName,
		Value:    accessToken,
		Path:     "/",
		Domain:   domain,
		MaxAge:   int((15 * time.Minute).Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    refreshToken,
		Path:     "/",
		Domain:   domain,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}

// ClearAuthCookies removes the access and refresh cookies by setting MaxAge to -1.
// The secure flag must match the value used when the cookies were originally set,
// otherwise the browser may not clear them correctly (e.g. over plain HTTP in dev).
func ClearAuthCookies(w http.ResponseWriter, domain string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     accessCookieName,
		Value:    "",
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     "/",
		Domain:   domain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}

// GetAccessToken returns the access token value from the request cookie,
// or an empty string if the cookie is absent.
func GetAccessToken(r *http.Request) string {
	c, err := r.Cookie(accessCookieName)
	if err != nil {
		return ""
	}
	return c.Value
}

// GetRefreshToken returns the refresh token value from the request cookie,
// or an empty string if the cookie is absent.
func GetRefreshToken(r *http.Request) string {
	c, err := r.Cookie(refreshCookieName)
	if err != nil {
		return ""
	}
	return c.Value
}
