package cmd

import (
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
		if len(ctx.Arguments) == 0 {
			return ctx.GenerateReplyMessage("Url?")
		}
		res, err := scrapper.GetSnapInsta(ctx.Arguments[0])
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
