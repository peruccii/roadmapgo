package dtos

type AuthInputDTO struct {
	Email    string
	Password string
}

type AuthOutputDTO struct {
	AccessToken string `json:"access_token"`
}
