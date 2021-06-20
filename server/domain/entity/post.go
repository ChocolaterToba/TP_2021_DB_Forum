package entity

import (
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type Post struct {
	PostID    int       `json:"id"`
	ParentID  int       `json:"parent"`
	Creator   string    `json:"author"`
	Message   string    `json:"message"`
	IsEdited  bool      `json:"isEdited"`
	Forumname string    `json:"forum"`
	ThreadID  int       `json:"thread"`
	Created   time.Time `json:"created"`
}

//easyjson:json
type Posts []*Post

type PostEditInput struct {
	Message string `json:"message"`
}

type RelatedObjectsInput struct {
	RelatedObjects map[string]bool `json:"related"`
}

func QueryToRelatedObjectsInput(query *fasthttp.Args) (*RelatedObjectsInput, error) {
	postInput := new(RelatedObjectsInput)
	postInput.RelatedObjects = make(map[string]bool)

	relatedObjectsBytes := query.Peek("related")
	relatedObjectStrings := string(relatedObjectsBytes)
	for _, relatedObjectString := range strings.Split(relatedObjectStrings, ",") {
		switch relatedObjectString {
		case "user", "forum", "thread":
			postInput.RelatedObjects[relatedObjectString] = true
		case "":
			continue
		default:
			return nil, UnsupportedRelatedObjectError
		}
	}
	return postInput, nil
}

type PostFullOutput struct {
	PostOutput   *Post   `json:"post"`
	ThreadOutput *Thread `json:"thread,omitempty"`
	ForumOutput  *Forum  `json:"forum,omitempty"`
	UserOutput   *User   `json:"author,omitempty"`
}
