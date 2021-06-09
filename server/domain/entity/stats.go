package entity

type StatsInfo struct {
	UsersCount   int `json:"users"`
	ForumsCount  int `json:"forums"`
	ThreadsCount int `json:"threads"`
	PostsCount   int `json:"posts"`
}
