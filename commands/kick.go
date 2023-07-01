package cmd

import (
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	waTypes "go.mau.fi/whatsmeow/types"
)

var kick = &command.Command{
	Name:         "kick",
	Aliases:      []string{"remove"},
	Description:  "Remove someone from group",
	GroupOnly:    true,
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		if len(ctx.Arguments) == 0 {
			return ctx.GenerateReplyMessage("Mention member\nEx: /kick @member")
		}
		for _, v := range util.ParseMentionedJid(ctx.Message) {
			jid, _ := waTypes.ParseJID(v)
			Member := map[waTypes.JID]whatsmeow.ParticipantChange{
				jid: whatsmeow.ParticipantChangeRemove,
			}
			ctx.Client.UpdateGroupParticipants(ctx.MessageInfo.Chat, Member)
		}
		return nil
	},
}

func init() {
	embed.Commands.Add(kick)
}
