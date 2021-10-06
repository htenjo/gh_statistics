package repository

type User struct {
	Id          int64
	SessionId   string
	AccessToken string
	Email       string
	Username    string
	Repos       string
}