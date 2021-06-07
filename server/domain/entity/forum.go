package entity

type Forum struct {
	ForumID      int    `json:"-"`
	Title        string `json:"title"`
	Creator      string `json:"user"`
	Forumname    string `json:"slug"`
	PostsCount   int    `json:"posts"`
	ThreadsCount int    `json:"threads"`
}

type ForumCreateInput struct {
	Title     string `json:"title"`
	Creator   string `json:"user"`
	Forumname string `json:"slug"`
}
