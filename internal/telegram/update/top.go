package update

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/benkenobi3/dick-and-dot/internal/database/repository"
)

const TopLimit = 10

func (h *handler) topCommand(ctx context.Context, chatID int64) (string, error) {
	topDicks, err := h.dicks.GetTopDicksByChatID(ctx, chatID, TopLimit)
	if err != nil {
		return "", fmt.Errorf("cannot get dicks for /top command: %w", err)
	}

	if len(topDicks) == 0 {
		return "У вас нет писек", nil
	}

	finalText := fmt.Sprintf("Топ %d членов этого чата: \n\n", TopLimit)
	topPos := topPositions(topDicks)

	for idx, dick := range topDicks {
		config := tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				ChatID: chatID,
				UserID: dick.UserID,
			},
		}

		chatMember, err := h.bot.GetChatMember(config)
		if err != nil {
			finalText += fmt.Sprintf("%d | [ДАННЫЕ УДАЛЕНЫ] - писька равна %d см \n", topPos[idx], dick.Length)
			continue
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
