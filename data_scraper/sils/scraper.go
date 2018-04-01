package sils

import (
	"strings"

	"golang.org/x/net/context"

	"github.com/PuerkitoBio/goquery"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
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
	for _, url := range urls {
		datas := Scrape(ctx, url)
		for _, data := range datas.Datas {
			cl.Datas = append(cl.Datas, data)
		}
	}
	return cl
}

func Scrape(ctx context.Context, url string) *Classes {
	client := urlfetch.Client(ctx)
	resp, err := client.Get(url)
	if err != nil {
		log.Infof(ctx, "could not get the url:", err)
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Infof(ctx, "could not create a new doc:", err)
	}

	var data []string
	doc.Find("div.ctable-main > table > tbody > tr > td").Each(func(i int, s *goquery.Selection) {
		s.Each(func(i int, s *goquery.Selection) {
			ss := s.Text()
			data = append(data, ss)
		})
	})
	if len(data) < 8 {
		return &Classes{}
	}
	schedule := strings.Split(data[4], "semester")
	classes := &Classes{}
	if strings.ContainsRune(data[4], 58) {
		classes = makeDuplicateForClassesTwiceAweek(url, data)
		return classes
	}
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
	classes.Datas = append(classes.Datas, class)
	return classes
}

func makeDuplicateForClassesTwiceAweek(url string, data []string) *Classes {
	classes := &Classes{}
	schedule := strings.Split(data[4], "semester")
	term := schedule[0]
	dayAndperiod := strings.Split(schedule[1], "／")
	FirstDayAndperiod := strings.Split(Split(dayAndperiod[0]), ".")
	SeconddayAndperiod := strings.Split(Split(dayAndperiod[1]), ".")
	class1 := &Class{
		School:      strings.TrimSpace(data[1]),
		CourseTitle: strings.TrimSpace(data[2]),
		Instructor:  strings.TrimSpace(data[3]),
		Term:        strings.TrimSpace(term),
		Day:         strings.TrimSpace(FirstDayAndperiod[0]),
		Period:      strings.TrimSpace(FirstDayAndperiod[1]),
		Category:    strings.TrimSpace(data[5]),
		Credit:      strings.TrimSpace(data[7]),
		URL:         url,
	}
	class2 := &Class{
		School:      strings.TrimSpace(data[1]),
		CourseTitle: strings.TrimSpace(data[2]),
		Instructor:  strings.TrimSpace(data[3]),
		Term:        strings.TrimSpace(term),
		Day:         strings.TrimSpace(SeconddayAndperiod[0]),
		Period:      strings.TrimSpace(SeconddayAndperiod[1]),
		Category:    strings.TrimSpace(data[5]),
		Credit:      strings.TrimSpace(data[7]),
		URL:         url,
	}
	classes.Datas = append(classes.Datas, class1, class2)
	return classes
}

func Split(s string) string {
	for i, val := range s {
		if val == rune(':') {
			return s[i+1:]
		}
	}
	return ""
}
