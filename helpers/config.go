package helpers

type (
	ConfigServerApi struct {
		Request       []string      `json:"request"`
		MiddlewareApi MiddlewareApi `json:"middleware"`
		Response      any           `json:"response"`
	}

	MiddlewareApi struct {
		Auth     string   `json:"auth"`
		Logging  bool     `json:"logging"`
		Security []string `json:"security"`
	}
)
