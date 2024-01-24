package dto

type CreateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserTokenRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserTokenResponse struct {
	AccessToken string `json:"access_token"`
}
