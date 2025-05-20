package middleware

type (
	AuthMiddleware struct{}
)

func (auth *AuthMiddleware) Jwt() {

}

func (auth *AuthMiddleware) BasicAuth() {

}

func (auth *AuthMiddleware) BearerToken() {

}
