package main

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/line/line-bot-sdk-go/linebot"
)

func (b *mybot) myReplyMessage(replyToken string, textMessage string, ctx context.Context) error {
	_, err := b.client.ReplyMessage(replyToken, linebot.NewTextMessage(textMessage)).WithContext(ctx).Do()
	return err
}

func (b *mybot) ReplyTemplate(replyToken string, template linebot.Template, ctx context.Context) error {
	_, err := b.client.ReplyMessage(replyToken, linebot.NewTemplateMessage("rrr", template)).WithContext(ctx).Do()
	return err
}

// func (r *mybot) Reply(ctx context.Context, replyToken string, message linebot.Message) error {
// 	if _, err := r.client.ReplyMessage(replyToken, message).WithContext(ctx).Do(); err != nil {
// 		fmt.Printf("Reply Error: %v", err)
// 		return err
// 	}
// 	return nil
// }

func (r *mybot) Reply(ctx context.Context, replyToken string, messages ...linebot.Message) error {
	const errtag = "mybot.Reply: "
	if _, err := r.client.ReplyMessage(replyToken, messages...).WithContext(ctx).Do(); err != nil {
		return errors.Wrap(err, errtag)
	}
	return nil
}
