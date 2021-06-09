package types

type (
	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	Token struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	JWTContext struct {
		User        string   `json:"user"`
		DisplayName string   `json:"display_name"`
		Roles       []string `json:"roles"`
	}
	TokenDetail struct {
		Token           Token       `json:"token"`
		AccessUUid      string      `json:"access_uuid"`
		RefreshUUid     string      `json:"refresh_uuid"`
		AccessTokenExp  int64       `json:"access_token_exp"`
		RefreshTokenExp int64       `json:"refresh_token_exp"`
		Context         interface{} `json:"context"`
	}
	RefrestToken struct {
		Token string `json:"refresh_token" form:"refresh_token"`
	}
)
