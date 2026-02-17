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
		return 0, 0, errors.New("invalid data format: expected 2 fields separated by comma")
	}
	stepsStr := strings.TrimSpace(parseDate[0])
	if stepsStr == "" {
		return 0, 0, errors.New("steps field is empty")
	}
	durationStr := strings.TrimSpace(parseDate[1])
	if durationStr == "" {
		return 0, 0, errors.New("duration field is empty")
	}

	step, err := strconv.Atoi(parseDate[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid steps value: %w", err)
	}
	if step <= 0 {
		return 0, 0, fmt.Errorf("steps must be positive, got: %d", step)
	}

	duration, err := time.ParseDuration(parseDate[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid duration format: %w", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("duration must be positive, got: %v", duration)
	}
	return step, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("error parsing package: %v", err)
		return ""
	}

	distance := float64(steps) * stepLength
	distanceKm := distance / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("error calculating calories: %v", err)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)
}
