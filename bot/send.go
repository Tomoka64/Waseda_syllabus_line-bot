package main

import (
	"context"

	"github.com/line/line-bot-sdk-go/linebot"
)

func (r *mybot) SendCarouselTemplate(ctx context.Context, replyToken, altText string,
	columns ...*linebot.CarouselColumn) error {

	return r.Reply(ctx, replyToken,
		linebot.NewTemplateMessage(altText,
			linebot.NewCarouselTemplate(columns...)))
}
func (r *mybot) SendCarouselTemplates(ctx context.Context, replyToken, altText string,
	columns [][]*linebot.CarouselColumn) error {
	p := r.NewTemplateMessages(altText,
		templatesMaker(ctx, columns))

	return r.Reply(ctx, replyToken, p...)
}
