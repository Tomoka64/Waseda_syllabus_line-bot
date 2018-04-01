package main

import (
	"fmt"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/net/context"
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
			b.serveTemplate(event.Postback.Data, event.ReplyToken, ctx)
		}

	}
}

func (b *mybot) handleText(message *linebot.TextMessage, replyToken string, ctx context.Context) {
	switch message.Text {
	case "学部選択":
		err := b.SendCarouselTemplate(ctx, replyToken, "学部選択", facultyTemplates()...)
		if err != nil {
			log.Errorf(ctx, "could not send the template 学部選択:", err)
		}
	}
}

func (b *mybot) serveTemplate(message, replyToken string, ctx context.Context) {
	d := strings.Split(message, ":")
	e := strings.Split(d[1], "&")
	log.Infof(ctx, "len of d: %d\n len of e: %d\n", len(d), len(e))
	switch d[0] {
	case "search":
		err := b.ReplyTemplate(replyToken, b.SearchTemplate(d[1]), ctx)
		if err != nil {
			log.Errorf(ctx, "could not send the template あ:", err)
		}
	case "period":
		fmt.Println(b.SendCarouselTemplate(ctx, replyToken, "period", periodTemplates(d[1])...))
	case "week":
		fmt.Println(b.SendCarouselTemplate(ctx, replyToken, "week", timeTemplates(e)...))
	case "schedule":
		b.Kensaku(ctx, replyToken, e[0], e[1], e[2])
	}
}
