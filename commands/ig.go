package cmd

import (
	"botwa/utils"
	"botwa/utils/scrapper"
	"fmt"

	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/util"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var ig = &command.Command{
	Name:         "ig",
	Aliases:      []string{"instagram", "insta"},
	Description:  "Download post from Instagram",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		var link string
		command.NewUserQuestion(ctx).
			SetQuestion("Please send media url link", &link).
			WithLikeEmoji().
			ExecWithParser()

		if link != "" {
			if !utils.ParseURL(link) {
				return ctx.GenerateReplyMessage("errors: invalid url scheme")
			}
		}
		
		res, err := scrapper.GetSnapInsta(link)
		if err != nil {
			fmt.Println(err)
			return ctx.GenerateReplyMessage("Error")
		}
		
		for _, result := range res.ResultMedia{
			bytes, err := util.DoHTTPRequest("GET", result)
			if err != nil {
				fmt.Println(err)
				return ctx.GenerateReplyMessage("Error while downloading media")
			}
			ctx.SendMessage(bytes)
		}

		return nil
	},
}

func init() {
	embed.Commands.Add(ig)
}
