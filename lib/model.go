package readit

import "time"

type Comment struct {
	Title string
	Text  string
	Added time.Time
}

type Link struct {
	Id       uint64
	Title    string
	Url      string
	Votes    int64
	Added    time.Time
	Comments []*Comment
}

func NewLink(id uint64, title string, url string) *Link {
	return &Link{
		Id:       id,
		Title:    title,
		Url:      url,
		Votes:    0,
		Added:    time.Now(),
		Comments: make([]*Comment, 0)}
}

func NewComment(title string, text string) *Comment {
	return &Comment{
		Title: title,
		Text:  text}
}

type ByVotesDesc []*Link

func (a ByVotesDesc) Len() int           { return len(a) }
func (a ByVotesDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVotesDesc) Less(i, j int) bool { return a[i].Votes > a[j].Votes }
