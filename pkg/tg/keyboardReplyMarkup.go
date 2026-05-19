package tglib

import tele "gopkg.in/telebot.v4"

func KeyboardReplyBuilder(buttons [][]string) *tele.ReplyMarkup {
	replyMarkup := &tele.ReplyMarkup{}
	var resultRows []tele.Row
	/*
		[
		["hello", "new",],
		[Hystory]
		]
	*/
	for _, list := range buttons {
		var btns []tele.Btn
		for _, btn := range list {
			btns = append(btns, replyMarkup.Text(btn))
		}
		resultRows = append(resultRows, replyMarkup.Row(btns...))
	}
	replyMarkup.Reply(resultRows...)

	// return &tele.SendOptions{
	// 	ReplyMarkup: replyMarkup,
	// }
	return replyMarkup
}
