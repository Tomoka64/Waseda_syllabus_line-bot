package main

import (
	"golang.org/x/net/context"

	"github.com/Tomoka64/Waseda_syllabus_line-bot/data_scraper/sils"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func (b *mybot) Kensaku(ctx context.Context, replyToken, faculty, day, times string) {
	var datas []*sils.Class
	log.Infof(ctx, "%s", faculty)

	q := datastore.NewQuery(faculty).Filter("day =", day).Filter("period =", times)
	num, err := q.Count(ctx)
	if err != nil {
		return
	}
	if num == 0 {
		b.myReplyMessage(replyToken, "見つかりませんでした", ctx)
		return
	}
	for it := q.Run(ctx); ; {
		var post sils.Class

		_, err := it.Next(&post)
		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Infof(ctx, "error during kensaku", err)
			break
		}
		log.Infof(ctx, "%d", len(datas))

		datas = append(datas, &post)

	}

	log.Infof(ctx, "%d", len(datas))
	b.SendCarouselTemplates(ctx, replyToken, "結果", Wrapper(ctx, datas))
}

func Wrapper(ctx context.Context, datas []*sils.Class) [][]*linebot.CarouselColumn {
	var t [][]*linebot.CarouselColumn
	n := len(datas) / 10
	if len(datas)%10 > 0 {
		n++
	}

	ln := 10
	for i := 0; i < n; i++ {
		if len(datas)-i*10 < 10 {
			ln = len(datas)
		}

		t = append(t, ResultTemplate(ctx, datas[i*10:ln]))
		ln += 10
	}

	return t
}

func ResultTemplate(ctx context.Context, datas []*sils.Class) []*linebot.CarouselColumn {
	x := len(datas)
	templates := make([]*linebot.CarouselColumn, x)
	log.Infof(ctx, "--%d--", x)
	for i, data := range datas[:x] {
		log.Infof(ctx, "====================== %d", i)
		url := linebot.NewURITemplateAction(FixLength(data.CourseTitle)+"...", data.URL)
		templates[i] = linebot.NewCarouselColumn(
			"", data.Day+data.Period, data.Instructor,
			url,
		)
	}

	return templates
}

func FixLength(title string) string {
	t := title
	for i := len(t); i < 17; i++ {
		t += "."
	}
	return t[:17]
}
