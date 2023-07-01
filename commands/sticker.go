package cmd

import (
	"botwa/utils"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/util"
	"github.com/itzngga/Roxy/util/cli"

	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var sticker = &command.Command{
	Name:         "sticker",
	Aliases:      []string{"s"},
	Description:  "Create whatsapp sticker",
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		fileName := time.Now().String()
		exifFile := "./file.exif"
		_, err := ctx.DownloadMessageToFile(ctx.Message, util.ParseQuotedRemoteJid(ctx.Message) != nil, fileName)
		if err != nil {
			return ctx.GenerateReplyMessage("Invalid")
		}
		if(ctx.Message.GetImageMessage() != nil || util.ParseQuotedMessage(ctx.Message).GetImageMessage() != nil){
			ok := Image2Webp(fileName)
			if !ok {
				return ctx.GenerateReplyMessage("Error while converting image")
			}
		}else if(ctx.Message.GetVideoMessage() != nil || util.ParseQuotedMessage(ctx.Message).GetVideoMessage() != nil){
			ok := Video2Webp(fileName)
			if !ok {
				return ctx.GenerateReplyMessage("Error while converting video")
			}
		} else {
			os.Remove(fileName)
			return ctx.GenerateReplyMessage("Invalid")
		}
		splited := strings.Split(strings.Join(ctx.Arguments, " "), "|")
		if len(splited) >= 2{
			utils.CreateExif(fileName+".exif", splited[0], splited[1])
			exifFile = fileName+".exif"
		}
		os.Remove(fileName)
		cli.ExecPipeline("webpmux", []byte{}, "-set", "exif", exifFile, fileName+".webp", "-o", fileName+".webp")
		if exifFile != "file.exif" {
			os.Remove(exifFile)
		}
		proto, err := ctx.UploadStickerMessageFromPath(fileName+".webp")
		if err != nil {
			fmt.Println(err)
			return nil
		}
		os.Remove(fileName+".webp")
		return &waProto.Message{
			StickerMessage: proto,
		}
	},
}

func Image2Webp(fileName string) bool {
	_, err := cli.ExecPipeline("ffmpeg", []byte{}, "-i", fileName, "-y", "-vcodec", "libwebp", "-vf", "scale='min(320,iw)':min'(320,ih)':force_original_aspect_ratio=decrease,fps=15, pad=320:320:-1:-1:color=white@0.0, split [a][b]; [a] palettegen=reserve_transparent=on:transparency_color=ffffff [p]; [b][p] paletteuse", "-f", "webp", fileName+".webp")
	if err != nil{
		fmt.Println(err)
		return false;
	}
	return true
}

func Video2Webp(fileName string) bool {
	_, err := cli.ExecPipeline("ffmpeg", []byte{}, "-f", "mp4", "-i", fileName, "-y", "-vcodec", "libwebp", "-vf", "scale='min(320,iw)':min'(320,ih)':force_original_aspect_ratio=decrease,fps=15, pad=320:320:-1:-1:color=white@0.0, split [a][b]; [a] palettegen=reserve_transparent=on:transparency_color=ffffff [p]; [b][p] paletteuse", "-f", "webp", fileName+".webp")
	if err != nil{
		fmt.Println(err)
		return false;
	}
	return true
}


func init() {
	embed.Commands.Add(sticker)
}