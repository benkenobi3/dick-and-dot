package update

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"

	"github.com/benkenobi3/dick-and-dot/internal/database/repository"
	"github.com/benkenobi3/dick-and-dot/internal/features/random"
)

const ErrorMessage = "Перелом члена. Возникла ошибка при обработке команды"

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
		text, err := h.dickCommand(ctx, update.Message.From.ID, update.Message.Chat.ID)
		msg.Text = text
		if err != nil {
			log.Println(err)
			msg.Text = ErrorMessage
		}
	case "top":
		text, err := h.topCommand(ctx, update.Message.Chat.ID)
		msg.Text = text
		if err != nil {
			log.Println(err)
			msg.Text = ErrorMessage
		}
	default:
		// ignore unknown command
		return
	}

	if msg.Text == "" {
		return
	}

	_, err := h.bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func (h *handler) dickCommand(ctx context.Context, userID, chatID int64) (string, error) {
	allDicks, err := h.dicks.GetDicksByChatId(ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("cannot get dicks for /dick command: %w", err)
	}

	currentDick, exists := allDicks[userID]
	if !exists {
		newLength, wasBlessed := random.GetNewLength(0, true)
		currentDick = repository.Dick{
			UserID: userID,
			ChatID: chatID,
			Length: newLength,
		}
		err = h.dicks.CreateDick(ctx, currentDick)
		if err != nil {
			return "", fmt.Errorf("cannot create new dick: %w", err)
		}

		message := ""
		if wasBlessed {
			message = fmt.Sprintf("🤡🤡🤡🤡🤡🤡🤡🤡\n\n"+
				"!!!!ВАУ!!!\n"+
				"ГОСПОДЬ ПОГЛАДИЛ ТЕБЯ ПО НОВОИСПЕЧЁННОЙ ГОЛОВКЕ\n"+
				"Прими благословение: %d см", random.BlessingSize)
		} else {
			message = fmt.Sprintf("Ты только что получил новый писюн, он равен %d см", currentDick.Length)
		}
		return message, nil
	}

	lastUpd := currentDick.UpdatedAt
	now := time.Now().UTC()
	ableToGrow := now.Day() > lastUpd.Day() || now.Month() > lastUpd.Month() || now.Year() > lastUpd.Year()

	if !ableToGrow {
		return fmt.Sprintf("Как же он наяривает...\nОстынь, писюн будет готов завтра"), nil
	}

	posForDick := getTopPosition(allDicks, userID)
	canBeBlessed := posForDick > 1
	newLength, wasBlessed := random.GetNewLength(currentDick.Length, canBeBlessed)
	diffLength := newLength - currentDick.Length

	currentDick.Length = newLength
	err = h.dicks.UpdateDick(ctx, currentDick)
	if err != nil {
		return "", fmt.Errorf("cannot update dick: %w", err)
	}

	allDicks[userID] = currentDick
	posForDick = getTopPosition(allDicks, userID)

	message := ""
	if wasBlessed {
		message = fmt.Sprintf("🤡🤡🤡🤡🤡🤡🤡🤡\n\n"+
			"!!!!ВАУ!!!\n"+
			"ГОСПОДЬ ПОГЛАДИЛ ТЕБЯ ПО ГОЛОВКЕ\n"+
			"Прими благословение: +%d см\n"+
			"Твой пипидастр размером в %d см взлетел на ангельских крылышках прямиком на %d место",
			random.BlessingSize, currentDick.Length, posForDick)
	} else {
		verb := "вырос"
		if diffLength < 0 {
			verb = "уменьшился"
			diffLength *= -1
		}
		message = fmt.Sprintf("Твой писюн %s на %d см, теперь он равен %d см.\n"+
			"Он разрывает чарты: %d место", verb, diffLength, currentDick.Length, posForDick)
	}

	return message, nil
}

func (h *handler) topCommand(ctx context.Context, chatID int64) (string, error) {

	allDicks, err := h.dicks.GetDicksByChatId(ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("cannot get dicks for /dick command: %w", err)
	}

	sortedDicks := sortDicks(allDicks)
	if len(sortedDicks) > 15 {
		sortedDicks = sortedDicks[:15] // cut to top 15 dicks
	}

	finalText := "Топ 15 членов этого чата: \n\n"
	topPos := topPositions(sortedDicks)
	for idx, dick := range sortedDicks {
		config := tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				ChatID: dick.ChatID,
				UserID: dick.UserID,
			},
		}

		chatMember, err := h.bot.GetChatMember(config)
		if err != nil {
			return "", fmt.Errorf("cannot get chat for /dick command: %w", err)
		}

		userName := chatMember.User.FirstName
		if chatMember.User.LastName != "" {
			userName = fmt.Sprintf("%s %s", userName, chatMember.User.LastName)
		}

		finalText += fmt.Sprintf("%d | %s - писька равна %d см \n", topPos[idx], userName, dick.Length)
	}

	return finalText, nil
}

func topPositions(sortedDicks []repository.Dick) []int {
	if len(sortedDicks) == 0 {
		return []int{}
	}

	topPos := 1
	positions := make([]int, 0, len(sortedDicks))
	positions = append(positions, topPos)

	for idx, dick := range sortedDicks {
		if idx != 0 {
			if dick.Length != sortedDicks[idx-1].Length {
				topPos++
			}
			positions = append(positions, topPos)
		}
	}

	return positions
}

func sortDicks(allDicks map[int64]repository.Dick) []repository.Dick {
	sortedDicks := make([]repository.Dick, 0, len(allDicks))
	for _, dick := range allDicks {
		sortedDicks = append(sortedDicks, dick)
	}

	if len(sortedDicks) < 2 {
		return sortedDicks
	}

	sort.SliceStable(sortedDicks, func(i, j int) bool {
		return sortedDicks[i].Length > sortedDicks[j].Length
	})

	return sortedDicks
}

func getTopPosition(allDicks map[int64]repository.Dick, userID int64) int {
	sortedDicks := sortDicks(allDicks)
	topPos := topPositions(sortedDicks)
	for i, dick := range sortedDicks {
		if dick.UserID == userID {
			return topPos[i]
		}
	}
	return -1
}
