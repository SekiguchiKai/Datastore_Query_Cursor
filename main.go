package main

import (
	"net/http"
	"html/template"
)

func init() {
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/record", HandleSample)
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./template/index.html"))
	tmpl.Execute(w, nil)
}








