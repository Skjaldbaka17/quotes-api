package structs

type UserApiModel struct {
	ApiKey               string `json:"apiKey"`
	Id                   int    `json:"id"`
	Email                string `json:"email"`
	Name                 string `json:"name"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
	Tier                 string `json:"tier"`
}

type UserDBModel struct {
	Id           int    `json:"id"`
	ApiKey       string `json:"api_key"`
	Message      string `json:"message"`
	PasswordHash string `json:"password_hash"`
	Tier         string `json:"tier"`
	Name         string `json:"name"`
	Email        string `json:"email"`
}
