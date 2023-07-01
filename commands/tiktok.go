package cmd

import (
	"botwa/utils"
	"botwa/utils/scrapper"
	"fmt"
	"strconv"

	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var tiktok = &command.Command{
	Name:         "tiktok",
	Aliases:      []string{"tt"},
	Description:  "Download video from Tiktok",
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

		res, err := scrapper.GetSnaptik(link)
		if err != nil {
			fmt.Println(err)
			return ctx.GenerateReplyMessage("Error")
		}

		text := fmt.Sprintf("Username: %s\nDescription: %s", res.Username, res.Description)
		
		videoMessage, err := ctx.UploadVideoFromUrl(res.VideoUrl[0], text)
		if err != nil {
			fmt.Println(err)
			return ctx.GenerateReplyMessage("Error while uploading video")
		}
		ctx.SendReplyMessage(videoMessage)

		for i, url := range res.ImageUrl {
			if i > 0 {
				text = "Slide: " + strconv.Itoa(i+1)
			}
			imageMessage, err := ctx.UploadImageFromUrl(url, text)
			if err != nil {
				fmt.Println(err)
				return ctx.GenerateReplyMessage("Error while uploading image")
			}
			ctx.SendReplyMessage(imageMessage)
		}

		return nil
	},
}

func init() {
	embed.Commands.Add(tiktok)
}
