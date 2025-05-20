package models

type (
	User struct {
		BaseModel
		Name     string
		Email    string
		Password string
	}
)
