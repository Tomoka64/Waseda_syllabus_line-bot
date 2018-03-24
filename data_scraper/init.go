package main

import (
	"net/http"

	sils "github.com/Tomoka64/syllabus_line_bot/data_scraper/sils"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func init() {
	http.HandleFunc("/api/sils", getSILS)
	http.Handle("/favicon.ico", http.NotFoundHandler())
}

func getSILS(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(ctx, "sils", nil)
	silsData := sils.GetData(ctx)
	for _, data := range silsData.Datas {
		if _, err := datastore.Put(ctx, key, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
