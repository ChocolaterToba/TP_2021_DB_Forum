package entity

type StatsInfo struct {
	UsersCount   int `json:"user"`
	ForumsCount  int `json:"forum"`
	ThreadsCount int `json:"thread"`
	PostsCount   int `json:"post"`
}
