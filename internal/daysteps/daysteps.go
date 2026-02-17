package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parseDate := strings.Split(data, ",")

	if len(parseDate) != 2 {
		return 0, 0, errors.New("Ошибка")
	}
	stepsStr := strings.TrimSpace(parseDate[0])
	if stepsStr == "" {
		return 0, 0, errors.New("Ошибка")
	}
	durationStr := strings.TrimSpace(parseDate[1])
	if durationStr == "" {
		return 0, 0, errors.New("Ошибка")
	}

	step, err := strconv.Atoi(parseDate[0])
	if err != nil {
		return 0, 0, errors.New("Ошибка")
	}
	if step <= 0 {
		return 0, 0, errors.New("Ошибка")
	}

	duration, err := time.ParseDuration(parseDate[1])
	if err != nil {
		return 0, 0, err
	}
	if duration <= 0 {
		return 0, 0, errors.New("Ошибка")
	}
	return step, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 {
		log.Println(err)
		return ""
	}

	if duration <= 0 {
		log.Println(err)
		return ""
	}

	distance := float64(steps) * stepLength
	distanceKm := distance / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
}
