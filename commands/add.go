package cmd

import (
	"github.com/itzngga/Roxy/command"
	"github.com/itzngga/Roxy/embed"
	"github.com/itzngga/Roxy/util"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

var add = &command.Command{
	Name:         "add",
	Aliases:      []string{"invite"},
	Description:  "Add someone to group",
	GroupOnly:    true,
	RunFunc: func(ctx *command.RunFuncContext) *waProto.Message {
		if len(ctx.Arguments) == 0 {
			return ctx.GenerateReplyMessage("Invalid number\nEx: /add 628...")
		}
		for _, v := range ctx.Arguments[0:] {
			jid, ok := util.ParseJID(v)
			if ok {
				found, err := ctx.Client.IsOnWhatsApp([]string{jid.User})
				if err == nil && len(found) != 0 {
					Member := map[types.JID]whatsmeow.ParticipantChange{
						jid: whatsmeow.ParticipantChangeAdd,
					}
					ctx.Client.UpdateGroupParticipants(ctx.MessageInfo.Chat, Member)
				}
			}
		}
		return nil
	},
}

func init() {
	embed.Commands.Add(add)
}
