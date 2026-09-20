package dto

type LoginInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutputDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}