package main

import (
	"fmt"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func (r *mybot) SendCarouselTemplate(ctx context.Context, replyToken, altText string,
	columns ...*linebot.CarouselColumn) error {

	return r.Reply(ctx, replyToken,
		linebot.NewTemplateMessage(altText,
			linebot.NewCarouselTemplate(columns...)))
}

func timeTemplates(datas []string) []*linebot.CarouselColumn {
	schedule := []string{
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"やめる",
	}

	templates := make([]*linebot.CarouselColumn, len(schedule))

	for i, time := range schedule {
		action := fmt.Sprintf("schedule:%s&%s", datas[0], time)
		postback := linebot.NewPostbackTemplateAction(time, action, "", "")
		templates[i] = linebot.NewCarouselColumn(
			"", "何限か", "選んでください",
			postback,
		)
	}

	return templates
}

func periodTemplates() []*linebot.CarouselColumn {
	days := []string{
		"月曜日",
		"火曜日",
		"水曜日",
		"木曜日",
		"金曜日",
		"土曜日",
		"やめる",
	}

	dayOfweek := map[string]string{
		"月曜日": "Mon", "火曜日": "Tues", "水曜日": "Wed",
		"木曜日": "Thu", "金曜日": "Fri", "土曜日": "Sat",
		"やめる": "end",
	}

	templates := make([]*linebot.CarouselColumn, len(dayOfweek))

	for i, day := range days {
		action := fmt.Sprintf("week:%s", dayOfweek[day])
		postback := linebot.NewPostbackTemplateAction(day, action, "", "")
		templates[i] = linebot.NewCarouselColumn(
			"", "曜日選択", "選んでください",
			postback,
		)
	}

	return templates
}

func (b *mybot) serveTemplate(message, replyToken string, ctx context.Context) {
	d := strings.Split(message, ":")
	e := strings.Split(d[1], "&")
	log.Infof(ctx, "======================")
	log.Infof(ctx, "len of d: %d\n len of e: %d\n", len(d), len(e))
	switch d[0] {
	case "period":
		fmt.Println(b.SendCarouselTemplate(ctx, replyToken, "period", periodTemplates()...))
	case "week":
		fmt.Println(b.SendCarouselTemplate(ctx, replyToken, "week", timeTemplates(d[1:])...))
	case "schedule":
		b.Kensaku(ctx, replyToken, e[0], e[1])
	}
}
