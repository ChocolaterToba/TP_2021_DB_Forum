package entity

type User struct {
	UserID      int    `json:"-"`
	Username    string `json:"nickname"`
	FullName    string `json:"fullname"`
	Description string `json:"about"`
	EMail       string `json:"email"`
}

//easyjson:json
type Users []*User

type UserEditInput struct {
	FullName    string `json:"fullname"`
	Description string `json:"about"`
	EMail       string `json:"email"`
}

type VoteInput struct {
	Username string `json:"nickname"`
	Vote     int    `json:"voice"`
}
