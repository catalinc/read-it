package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/catalinc/readit/lib"
)

const indexHtml = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>HNews</title>
		<style type="text/css">a {text-decoration: none}</style>
	</head>
	<body>
		{{range .}}
		<div>
			<a title="Vote" href="/vote/{{.Id}}">&#9650;</a>&nbsp;
			<a title="{{.Title}}" href="/view/{{.Id}}">{{.Title}}</a>&nbsp;
			<small>{{.Votes}} votes</small>
		</div>
		{{else}}
			<div>
				<strong>No links yet</strong>
			</div>
		{{end}}
		 <form action="/add" method="post">
		 	<label for="title">Title</label>	
		 	<input type="text" name="title">
		 	<label for="url">URL</label>	
		 	<input type="text" name="url">
		 	<input type="submit" value="Add">
		 </form>
	</body>
</html>`

const viewLinkHtml = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>View Link</title>
		<style type="text/css">a {text-decoration: none}</style>
	</head>
	<body>
		<a href="{{.Url}}" title="{{.Title}}" target="_blank">{{.Title}}</a>
		{{range .Comments}}
		<div>
			{{.Text}}&nbsp;	(<small>{{.Added}}</small>)
		</div>
		{{else}}
			<div>
				<strong>No comments yet</strong>
			</div>
		{{end}}
		<a href="/" title="Back">Back</a>
		<form action="/comment/{{.Id}}" method="post">
		 	<label for="title">Title</label>	
		 	<input type="text" name="title">
		 	<label for="text">Text</label>
		 	<input type="text" name="text">		 	
		 	<input type="submit" value="Post">
		</form>
	</body>
</html>`

var indexTemplate = template.Must(template.New("index").Parse(indexHtml))
var viewLinkTemplate = template.Must(template.New("view").Parse(viewLinkHtml))

func index(w http.ResponseWriter, r *http.Request) {
	sort.Sort(readit.ByVotesDesc(readit.Links))

	err := indexTemplate.Execute(w, readit.Links)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}

func addLink(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	url := r.FormValue("url")

	if len(title) == 0 || len(url) == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	readit.AddLink(title, url)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func voteLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Path[len("/vote/"):], 10, 64)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	link := readit.GetLink(id)
	if link == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	link.Votes++
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func viewLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Path[len("/view/"):], 10, 64)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	link := readit.GetLink(id)
	if link == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	err = viewLinkTemplate.Execute(w, link)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}

func commentLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Path[len("/comment/"):], 10, 64)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	title := r.FormValue("title")
	text := r.FormValue("text")
	if len(title) == 0 || len(text) == 0 {
		http.Redirect(w, r, fmt.Sprintf("/link/%d", id), http.StatusSeeOther)
		return
	}

	link := readit.GetLink(id)
	if link == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	readit.AddComment(link, title, text)

	err = viewLinkTemplate.Execute(w, link)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/add", addLink)
	http.HandleFunc("/vote/", voteLink)
	http.HandleFunc("/view/", viewLink)
	http.HandleFunc("/comment/", commentLink)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
