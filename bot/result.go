package main

import (
	"golang.org/x/net/context"

	"github.com/Tomoka64/Waseda_syllabus_line-bot/data_scraper/sils"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func (b *mybot) Kensaku(ctx context.Context, replyToken, day, times string) {
	var datas []*sils.Class
	log.Infof(ctx, "%s", times)

	q := datastore.NewQuery("SILS").Filter("day =", day).Filter("period =", times)
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
	b.SendCarouselTemplate(ctx, replyToken, "結果", ResultTemplate(datas)...)
}

func ResultTemplate(datas []*sils.Class) []*linebot.CarouselColumn {
	x := len(datas)
	switch {
	case x < 10:
		templates := make([]*linebot.CarouselColumn, x)
		for i, data := range datas {
			url := linebot.NewURITemplateAction(data.CourseTitle[:17]+"...", data.URL)
			templates[i] = linebot.NewCarouselColumn(
				"", data.Day+data.Period, data.Instructor,
				url,
			)
			return templates
		}
	case x >= 10:
		templates := make([]*linebot.CarouselColumn, len(datas[:10]))
		for i, data := range datas[:10] {
			url := linebot.NewURITemplateAction(data.CourseTitle[:17]+"...", data.URL)
			templates[i] = linebot.NewCarouselColumn(
				"", data.Day+data.Period+"\t"+data.Category, data.Instructor,
				url,
			)
		}
		return templates
	}

	return nil
}
