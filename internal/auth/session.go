package auth

import "time"

const tokenExpirySkew = 30 * time.Second

type Session struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	UserID       string    `json:"user_id,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
}

func (s Session) HasAccessToken() bool {
	return s.AccessToken != ""
}

func (s Session) AccessTokenExpired(now time.Time) bool {
	if s.ExpiresAt.IsZero() {
		return false
	}

	return !now.Before(s.ExpiresAt.Add(-tokenExpirySkew))
}
