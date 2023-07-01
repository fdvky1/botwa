package cmd

import (
	"botwa/utils/scrapper"
	"fmt"

	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var tiktok = &command.Command{
	Name:         "tiktok",
	Aliases:      []string{"tt"},
	Description:  "Download video from Tiktok",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		if len(ctx.Arguments) == 0 {
			return ctx.GenerateReplyMessage("Url?")
		}
		res, err := scrapper.GetSnaptik(ctx.Arguments[0])
		if err != nil {
			fmt.Println(err)
			return ctx.GenerateReplyMessage("Error")
		}
		
		for _, result := range res.VideoUrl{
			videoMessage, err := ctx.UploadVideoFromUrl(result, "")
			if err != nil {
				fmt.Println(err)
				return ctx.GenerateReplyMessage("Error while uploading video")
			}
			ctx.SendReplyMessage(videoMessage)
		}

		return nil
	},
}

func init() {
	embed.Commands.Add(tiktok)
}
