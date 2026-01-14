package models

import "strconv"

type (
	Auth struct{}
)

func (a Auth) GetBearerTokenExpiration() int {
	v, _v := env.GetEnv("bearer_token_expiration")

	if v {
		time, _ := strconv.Atoi(_v)
		return time
	} else {
		return 60 // default seconds
	}

}

func (a Auth) GetJwtExpiration() int {
	v, _v := env.GetEnv("jwt_expiration")

	if v {
		time, _ := strconv.Atoi(_v)
		return time
	} else {
		return 60 // default seconds
	}

}

func (a Auth) GetJwtSecretKey() string {
	v, _v := env.GetEnv("jwt_secret_key")

	if v {
		return _v
	} else {
		return "my_secret_key"
	}
}
