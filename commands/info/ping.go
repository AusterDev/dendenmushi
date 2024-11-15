package commands

import (
	"fmt"
	"dendenmushi/handlers"

	"github.com/bwmarrin/discordgo"
)

type PingCmd struct{}

func (c PingCmd) GetMeta() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "ping",
		Description: "Shows the ping of the connection",
	}
}

func (c PingCmd) Run(ctx *handlers.Ctx) error {
	ctx.Reply(fmt.Sprintf("💓 Heartbeat Latency: **`%v`**", ctx.Session.HeartbeatLatency()))
	return nil
}