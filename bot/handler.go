package main

import (
	"fmt"

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
	case "あ":
		err := b.ReplyTemplate(replyToken, b.TemplateHandler("検索方法"), ctx)
		if err != nil {
			log.Infof(ctx, "could not send the template あ:", err)
		}
	case "い":
		err := b.ReplyTemplate(replyToken, b.TemplateHandler("検索方法"), ctx)
		if err != nil {
			log.Infof(ctx, "could not send the template あ:", err)
		}
	}
	b.FirstPageTemplate(replyToken)
}

func (b *mybot) TemplateHandler(message string) *linebot.ButtonsTemplate {
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
