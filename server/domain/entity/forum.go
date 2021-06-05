package entity

type Forum struct {
	ForumID      int    `json:"-"`
	Title        string `json:"title"`
	Creator      string `json:"user"`
	IDString     string `json:"slug"`
	postsCount   int    `json:"posts"`
	threadsCount int    `json:"threads"`
}
