package readit

import "sync/atomic"

var (
	idCounter uint64  = 0
	Links     []*Link = make([]*Link, 0)
)

func GetLink(id uint64) *Link {
	for _, link := range Links {
		if link.Id == id {
			return link
		}
	}
	return nil
}

func AddLink(title string, url string) *Link {
	atomic.AddUint64(&idCounter, 1)

	link := NewLink(idCounter, title, url)
	Links = append(Links, link)

	return link
}

func AddComment(link *Link, title string, text string) *Comment {
	comment := NewComment(title, text)
	link.Comments = append(link.Comments, comment)

	return comment
}
