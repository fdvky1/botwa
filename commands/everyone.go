package cmd

import (
	"fmt"
	"strings"

	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var everyone = &command.Command{
	Name:         "everyone",
	Aliases:      []string{"mentionall"},
	Description:  "Mention All Members",
	GroupOnly:    true,
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		info, err := ctx.Client.GetGroupInfo(ctx.MessageInfo.Chat)
		if err != nil {
			return ctx.GenerateReplyMessage("Cant get group info")
		}
		members := []string{}
		for _, member := range info.Participants{
			members = append(members, member.JID.String())
		}
		text := "@everyone"
		if len(ctx.Arguments) > 1 {
			text = fmt.Sprintf("%s\n\n%s", strings.Join(ctx.Arguments, " "), text)
		}

		return &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text:        &text,
				ContextInfo: &waProto.ContextInfo {
				  MentionedJid: members,
				},
			},
		}
	},
}

func init() {
	embed.Commands.Add(everyone)
}
