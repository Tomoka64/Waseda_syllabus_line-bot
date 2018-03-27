package main

import (
	"fmt"
	"net/http"

	"github.com/Tomoka64/Waseda_syllabus_line-bot/data_scraper/sils"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func init() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/api/sils", getSILS)
	http.Handle("/favicon.ico", http.NotFoundHandler())
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}
func getSILS(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(ctx, "SILS", nil)
	silsData := sils.GetData(ctx)
	for _, data := range silsData.Datas {
		if _, err := datastore.Put(ctx, key, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
