package github

type ghCredentials struct {
	clientId     string
	clientSecret string
}

type OAuthCredentials struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (c *OAuthCredentials) getFullTokenHeader() string {
	return c.TokenType + " " + c.AccessToken
}

type GhUser struct {
	Id        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
}
