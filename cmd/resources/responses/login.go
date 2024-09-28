package responses

// LoginResponse represents the response when a user successfully logs in
type LoginResponse struct {
	// Token - the JWT token to be used for authorization
	Token string `json:"token"`
	// Name - the user's name
	Name string `json:"name"`
	// ID - the user's ID
	ID string `json:"id"`
}
