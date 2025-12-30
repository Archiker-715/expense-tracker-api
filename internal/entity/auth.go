package entity

type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}
