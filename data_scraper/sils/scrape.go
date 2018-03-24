package sils

import (
	"log"
	"strings"

	"context"

	"github.com/PuerkitoBio/goquery"
)

const (
	filePath = "syllabus_url.txt"
)

type Class struct {
	School      string `datastore:"school"`
	CourseTitle string `datastore:"coursetitle"`
	Instructor  string `datastore:"instructor"`
	Term        string `datastore:"term"`
	Day         string `datastore:"day"`
	Period      string `datastore:"period"`
	Category    string `datastore:"category"`
	Credit      string `datastore:"credit"`
	URL         string `datastore:"url"`
}

type Classes struct {
	Datas []*Class
}

func GetData(ctx context.Context) *Classes {
	urls := FromFile(filePath)
	cl := &Classes{}
	for _, v := range urls {
		cl.Datas = append(cl.Datas, Scrape(ctx, v))
	}

	return cl
}

func Scrape(ctx context.Context, url string) *Class {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	var data []string
	doc.Find("div.ctable-main > table > tbody > tr > td").Each(func(i int, s *goquery.Selection) {
		s.Each(func(i int, s *goquery.Selection) {
			ss := s.Text()
			data = append(data, ss)
		})
	})
	schedule := strings.Split(data[4], "semester")
	term := schedule[0]
	dayAndperiod := strings.Split(schedule[1], ".")
	day := dayAndperiod[0]
	period := dayAndperiod[1]
	class := &Class{
		School:      strings.TrimSpace(data[1]),
		CourseTitle: strings.TrimSpace(data[2]),
		Instructor:  strings.TrimSpace(data[3]),
		Term:        strings.TrimSpace(term),
		Day:         strings.TrimSpace(day),
		Period:      strings.TrimSpace(period),
		Category:    strings.TrimSpace(data[5]),
		Credit:      strings.TrimSpace(data[7]),
		URL:         url,
	}
	return class
}