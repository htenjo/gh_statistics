package storage

type User struct {
	Id          int64
	SessionId   string
	AccessToken string
	Email       string
	Username    string
}
