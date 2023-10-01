package update

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type handler struct {
	tgBot *tgbotapi.BotAPI
}

type Handler interface {
	Handle(update tgbotapi.Update)
}

func NewHandler(tgBot *tgbotapi.BotAPI) Handler {
	return &handler{
		tgBot: tgBot,
	}
}

func (h *handler) Handle(update tgbotapi.Update) {

}
