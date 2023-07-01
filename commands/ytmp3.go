package cmd

import (
	"botwa/utils"
	"botwa/utils/scrapper"
	"fmt"

	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var ytmp3 = &command.Command{
	Name:         "ytmp3",
	Aliases:      []string{"mp3"},
	Description:  "Download audio from YouTube",
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

		res, err := scrapper.Y2Mate(link)
		if err != nil {
			fmt.Println(err)
			return ctx.GenerateReplyMessage("Error")
		}
		imageMessage, err := ctx.UploadImageFromUrl(fmt.Sprintf("https://i.ytimg.com/vi/%s/0.jpg", res.Vid), fmt.Sprintf("%s\n\nChannel: %s", res.Title, res.A))
		if err != nil {
			fmt.Println(err)
			return ctx.GenerateReplyMessage("Error while uploading image")
		}
		ctx.SendReplyMessage(imageMessage)
		//correct me if wrong:v
		for _, v := range res.Links.Mp3{
			convert, err := scrapper.Y2MateDownloadUrl(res.Vid, v.K)
			if err == nil {
				audioMessage, err := ctx.UploadAudioFromUrl(convert.Dlink)
				if err != nil {
					fmt.Println(err)
					return ctx.GenerateReplyMessage("Error while uploading audio")
				}
				ctx.SendReplyMessage(audioMessage)
				break
			}
		}
		return nil
	},
}

func init() {
	embed.Commands.Add(ytmp3)
}
