package main

import (
	"fmt"
	"strings"

	"github.com/Tomoka64/syllabus_line_bot/data_scraper/sils"
	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func (b *mybot) EventRouter(eve []*linebot.Event, ctx context.Context) {
	for _, event := range eve {
		switch event.Type {
		case linebot.EventTypeMessage:

			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				b.handleText(message, event.ReplyToken, ctx)
				fmt.Print(message)
			}

		case linebot.EventTypePostback:
			//r.handledata(event.Postback.Data, event.ReplyToken, event.Source.UserID)
			b.searchTemplate(event.Postback.Data, event.ReplyToken, ctx)
		}

	}
}

func (b *mybot) handleText(message *linebot.TextMessage, replyToken string, ctx context.Context) {
	switch message.Text {
	case "あ":
		err := b.ReplyTemplate(replyToken, b.SendSilsTemplate("検索方法"), ctx)
		if err != nil {
			log.Infof(ctx, "could not send the template あ:", err)
		}
	case "い":
		b.SendSpseTemplate(replyToken)
	}
	b.FirstPageTemplate(replyToken)
}

func (b *mybot) SendSilsTemplate(message string) *linebot.ButtonsTemplate {
	return linebot.NewButtonsTemplate(
		"", message, "曜日検索",
		linebot.NewPostbackTemplateAction("曜日で検索", "period:", "", ""),
		linebot.NewPostbackTemplateAction("レベルで検索", "level", "", ""),
		linebot.NewPostbackTemplateAction("終わる", "j", "終わる", ""),
	)
}
func (b *mybot) SendSpseTemplate(message string) *linebot.ButtonsTemplate {
	return nil
}
func (b *mybot) FirstPageTemplate(message string) *linebot.ButtonsTemplate {
	return nil
}

func (b *mybot) searchTemplate(message, replyToken string, ctx context.Context) {
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

func (r *mybot) SendCarouselTemplate(ctx context.Context, replyToken, altText string,
	columns ...*linebot.CarouselColumn) error {

	return r.Reply(ctx, replyToken,
		linebot.NewTemplateMessage(altText,
			linebot.NewCarouselTemplate(columns...)))
}

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

	templates := make([]*linebot.CarouselColumn, len(datas[:10]))

	for i, data := range datas[:10] {
		url := linebot.NewURITemplateAction(data.CourseTitle[:15]+"...", data.URL)
		templates[i] = linebot.NewCarouselColumn(
			"", data.Day+data.Period, "選んでください",
			url,
		)
	}

	return templates
}
