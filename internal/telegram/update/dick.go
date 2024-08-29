package update

import (
	"context"
	"fmt"
	"time"

	"github.com/benkenobi3/dick-and-dot/internal/database/repository"
	"github.com/benkenobi3/dick-and-dot/internal/features/random"
)

func abs(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}

func (h *handler) dickCommand(ctx context.Context, userID, chatID int64) (string, error) {
	dick, err := h.dicks.GetDick(ctx, chatID, userID)
	if err != nil {
		return "", fmt.Errorf("error in /dick command: %w", err)
	}

	if dick.UpdatedAt == nil {
		newLength, wasBlessed := random.GetLengthToAdd(true)

		// первый раз прощаем и не выдаем сходу отрицательный член
		newLength = abs(newLength)

		newDick := repository.Dick{
			UserID: userID,
			ChatID: chatID,
			Length: newLength,
		}
		err = h.dicks.AddDick(ctx, newDick)
		if err != nil {
			return "", fmt.Errorf("cannot add new dick: %w", err)
		}

		if wasBlessed {
			return fmt.Sprintf("🤡🤡🤡🤡🤡🤡🤡🤡\n\n"+
				"!!!!ВАУ!!!\n"+
				"ГОСПОДЬ ПОГЛАДИЛ ТЕБЯ ПО НОВОИСПЕЧЁННОЙ ГОЛОВКЕ\n"+
				"Прими благословение: %d см", random.BlessingSize), nil
		}

		return fmt.Sprintf("Ты только что получил новый писюн, он равен %d см", newLength), nil
	}

	now := time.Now().UTC()
	lastUpdated := dick.UpdatedAt.UTC()

	ableToGrow := now.Day() > lastUpdated.Day() || now.Month() > lastUpdated.Month() || now.Year() > lastUpdated.Year()

	if !ableToGrow {
		return fmt.Sprintf("Как же он наяривает...\nОстынь, писюн будет готов завтра"), nil
	}

	lengthToAdd, wasBlessed := random.GetLengthToAdd(true)
	newDick := repository.Dick{
		UserID: userID,
		ChatID: chatID,
		Length: lengthToAdd,
	}
	err = h.dicks.AddDick(ctx, newDick)
	if err != nil {
		return "", fmt.Errorf("cannot add new dick: %w", err)
	}

	if wasBlessed {
		return fmt.Sprintf("🤡🤡🤡🤡🤡🤡🤡🤡\n\n"+
			"!!!!ВАУ!!!\n"+
			"ГОСПОДЬ ПОГЛАДИЛ ТЕБЯ ПО ГОЛОВКЕ\n"+
			"Прими благословение: +%d см\n"+
			"Твой пипидастр размером в %d см взлетел на ангельских крылышках",
			random.BlessingSize, dick.Length+random.BlessingSize), nil
	}

	verb := "вырос"
	if lengthToAdd < 0 {
		verb = "уменьшился"
	}

	return fmt.Sprintf("Твой писюн %s на %d см, теперь он равен %d см",
		verb,
		abs(lengthToAdd),
		dick.Length+lengthToAdd), nil
}
