package update

import (
	"context"
	"github.com/benkenobi3/dick-and-dot/internal/database/repository"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

type handler struct {
	bot   *tgbotapi.BotAPI
	dicks repository.Dicks
}

type Handler interface {
	Handle(ctx context.Context, update tgbotapi.Update)
}

func NewHandler(db *sqlx.DB, bot *tgbotapi.BotAPI) Handler {
	dicks := repository.NewDicks(db)

	return &handler{
		bot:   bot,
		dicks: dicks,
	}
}

func (h *handler) Handle(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if !update.Message.IsCommand() {
		// ignore any non-command messages
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "start":
		msg.Text = "Пора меряться письками! Чтобы начать, напиши /dick"
	case "help":
		msg.Text = "Бот создан, чтобы ты мог помериться письками в чате /n/n" +
			"Напиши /dick, чтобы отдать свой писюн богу рандома, он изменит его от -5 до +10 см /n" +
			"Напиши /top, чтобы увидеть топ писек в этом чате /n" +
			"Напиши /help, чтобы еще раз прочитать это сообщение /n"
	case "dick":
		msg.Text = ""
	default:
		// ignore unknown command
		return
	}

	_, err := h.bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}
