package entity

type User struct {
	UserID      int    `json:"-"`
	Username    string `json:"nickname"`
	FullName    string `json:"fullname"`
	Description string `json:"about"`
	EMail       string `json:"email"`
}
