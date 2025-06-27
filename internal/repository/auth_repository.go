package repository

type authParams struct {
	Email    string
	Password string
}

type AuthRepository interface {
	AuthUser(params authParams)
}

func AuthUser(params *authParams) (string, error) {
	return "aa", nil
}
