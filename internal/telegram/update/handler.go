package update

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"

	"github.com/benkenobi3/dick-and-dot/internal/database/repository"
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
		msg.Text = "Бот создан, чтобы ты мог помериться письками в чате \n\n" +
			"/dick - отдать свой писюн богу рандома, он изменит его от -5 до +10 см \n" +
			"/top - увидеть топ писек в этом чате \n" +
			"/help - еще раз прочитать это сообщение \n\n" +
			"Добавляй бот в свои чаты!"
	case "dick":
		msg.Text = h.dickCommand(ctx, update.Message.From.ID, update.Message.Chat.ID)
	case "top":
		msg.Text = h.topCommand(ctx, update.Message.Chat.ID)
	default:
		// ignore unknown command
		return
	}

	_, err := h.bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func (h *handler) dickCommand(ctx context.Context, userID, chatID int64) string {
	return ""
}

func (h *handler) topCommand(ctx context.Context, chatID int64) string {
	return ""
}
