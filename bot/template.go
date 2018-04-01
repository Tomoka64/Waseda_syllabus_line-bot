package main

import (
	"fmt"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/net/context"
)

func templatesMaker(ctx context.Context, colomuns [][]*linebot.CarouselColumn) []linebot.Template {
	r := []linebot.Template{}
	for _, v := range colomuns {
		r = append(r, linebot.NewCarouselTemplate(v...))
	}

	return r
}
func (b *mybot) NewTemplateMessages(altText string, templates []linebot.Template) []linebot.Message {
	var t []linebot.Message
	for _, template := range templates {
		a := linebot.NewTemplateMessage(altText, template)
		t = append(t, a)
	}
	return t
}

func (b *mybot) SearchTemplate(message string) *linebot.ButtonsTemplate {
	return linebot.NewButtonsTemplate(
		"", message, "検索",
		linebot.NewPostbackTemplateAction("曜日で検索", "period:"+message, "", ""),
		linebot.NewPostbackTemplateAction("レベルで検索", "level:", "", ""),
		linebot.NewPostbackTemplateAction("終わる", "j", "終わる", ""),
	)
}

func facultyTemplates() []*linebot.CarouselColumn {
	faculties := [][]string{
		{"国際教養学部/SILS", "政治経済学部/PSE"},
		{"文化構想学部/CMS", "教育学部/EDU"},
		{"基幹理工学部/FSE", "先進理工学部/ASE"},
		{"人間科学部/HUM", "法学部/LAW"},
		{"文学部/HSS", "商学部/SOC"},
		{"創造理工学部/CSE", "社会科学部/SSS"},
		{"スポーツ科学部/SPS", "やめる/end"},
	}

	templates := make([]*linebot.CarouselColumn, len(faculties))
	for i, faculty := range faculties {
		ret1 := strings.Split(faculty[0], "/")
		ret2 := strings.Split(faculty[1], "/")
		action1 := fmt.Sprintf("search:%s", ret1[1])
		action2 := fmt.Sprintf("search:%s", ret2[1])
		postback1 := linebot.NewPostbackTemplateAction(faculty[0], action1, "", "")
		postback2 := linebot.NewPostbackTemplateAction(faculty[1], action2, "", "")
		templates[i] = linebot.NewCarouselColumn(
			"", "学部を選んでください", "pick your faculty",
			postback1,
			postback2,
		)
	}
	return templates
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
		action := fmt.Sprintf("schedule:%s&%s&%s", datas[0], datas[1], time)
		postback := linebot.NewPostbackTemplateAction(time, action, "", "")
		templates[i] = linebot.NewCarouselColumn(
			"", "何限か", "選んでください",
			postback,
		)
	}

	return templates
}

func periodTemplates(message string) []*linebot.CarouselColumn {
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
		"木曜日": "Thur", "金曜日": "Fri", "土曜日": "Sat",
		"やめる": "end",
	}

	templates := make([]*linebot.CarouselColumn, len(dayOfweek))

	for i, day := range days {
		action := fmt.Sprintf("week:%s&%s", message, dayOfweek[day])
		postback := linebot.NewPostbackTemplateAction(day, action, "", "")
		templates[i] = linebot.NewCarouselColumn(
			"", "曜日選択", "選んでください",
			postback,
		)
	}

	return templates
}
