package github

import "time"

const (
	Green  PrReviewFlag = "green"
	Yellow PrReviewFlag = "yellow"
	Red    PrReviewFlag = "red"
)

type PrReviewFlag string

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

type PullRequestDetail struct {
	Url        string       `json:"url"`
	HtmlUrl    string       `json:"html_url"`
	Title      string       `json:"title"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	ReviewFlag PrReviewFlag `json:"review_flag"`
}

type RepoPR struct {
	RepositoryName string
	RepositoryURL  string
	Prs            []PullRequestDetail
}
