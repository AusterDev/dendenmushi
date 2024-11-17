package handlers

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

type CommandHandler struct {
	commands map[string]SlashCmd
	s        *discordgo.Session
}

func NewCmdHandler(s *discordgo.Session) *CommandHandler {
	return &CommandHandler{
		commands: make(map[string]SlashCmd),
		s:        s,
	}
}

func (h *CommandHandler) Add(cmds []SlashCmd, sync bool) {
	for _, cmd := range cmds {
		m := cmd.GetMeta()

		if m == nil {
			log.Println("GetMeta retured nil for command")
			continue
		}
		h.commands[m.Name] = cmd

		log.Printf("Command %s has been loaded", m.Name)

		if sync {
			if _, err := h.s.ApplicationCommandCreate(h.s.State.User.ID, "", m); err != nil {
				log.Printf("Failed to create %s command: %v", m.Name, err)
			} else {
				log.Printf("Command %s has been created on discord", m.Name)
			}
		}
	}
}

func (h *CommandHandler) Get(name string) SlashCmd {
	return h.commands[name]
}

func (h *CommandHandler) Handle(i *discordgo.Interaction) {
	cmd := h.Get(i.ApplicationCommandData().Name)
	if cmd != nil {
		ctx := NewCtx(h.s, cmd.GetMeta(), i)
		go func() {
			if err := cmd.Run(ctx); err != nil {
				log.Printf("Command %s failed: %v", cmd.GetMeta().Name, err)
			}
		}()
	}

}

type Ctx struct {
	Session     *discordgo.Session
	Meta        *discordgo.ApplicationCommand
	Interaction *discordgo.Interaction
}

type ComplexMessage struct {
	Content string
	Data    *discordgo.InteractionResponseData
	Defer   bool
}

func NewCtx(s *discordgo.Session, m *discordgo.ApplicationCommand, i *discordgo.Interaction) *Ctx {
	return &Ctx{s, m, i}
}

func (ctx *Ctx) Reply(content string) error {
	err := ctx.Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	return err
}

func (ctx *Ctx) ReplyComplex(msg ComplexMessage) error {
	var msgType discordgo.InteractionResponseType
	if msg.Defer {
		msgType = discordgo.InteractionResponseDeferredChannelMessageWithSource
	} else {
		msgType = discordgo.InteractionResponseChannelMessageWithSource
	}

	err := ctx.Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
		Type: msgType,
		Data: msg.Data,
	})
	return err
}

type SlashCmd interface {
	GetMeta() *discordgo.ApplicationCommand
	Run(ctx *Ctx) error
}
