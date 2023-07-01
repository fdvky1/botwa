package cmd

import (
	"botwa/utils/scrapper"
	"fmt"

	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var ytmp4 = &command.Command{
	Name:         "ytmp4",
	Aliases:      []string{"mp4"},
	Description:  "Download video from YouTube",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		if len(ctx.Arguments) == 0 {
			return ctx.GenerateReplyMessage("Url?")
		}
		res, err := scrapper.Y2Mate(ctx.Arguments[0])
		if err != nil {
			fmt.Println(err)
			return ctx.GenerateReplyMessage("Error")
		}
		imageMessage, err := ctx.UploadImageFromUrl(fmt.Sprintf("https://i.ytimg.com/vi/%s/0.jpg", res.Vid), fmt.Sprintf("%s\nAuthor: %s", res.Title, res.A))
		if err != nil {
			fmt.Println(err)
			return ctx.GenerateReplyMessage("Error while uploading image")
		}
		ctx.SendReplyMessage(imageMessage)
		//correct me if wrong:v
		for _, v := range res.Links.Mp4{
			convert, err := scrapper.Y2MateDownloadUrl(res.Vid, v.K)
			if err == nil {
				videoMessage, err := ctx.UploadVideoFromUrl(convert.Dlink, "There!")
				if err != nil {
					fmt.Println(err)
					return ctx.GenerateReplyMessage("Error while uploading video")
				}
				ctx.SendReplyMessage(videoMessage)
				break
			}
		}
		return nil
	},
}

func init() {
	embed.Commands.Add(ytmp4)
}
