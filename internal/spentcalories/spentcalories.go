package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	stepLengthCoefficient      = 0.45
	mInKm                      = 1000
	minInH                     = 60
	walkingCaloriesCoefficient = 0.5
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parseStr := strings.Split(data, ",")
	if len(parseStr) != 3 {
		return 0, "", 0, fmt.Errorf("Ошибка")
	}
	stepsStr := strings.TrimSpace(parseStr[0])
	if stepsStr == "" {
		return 0, "", 0, errors.New("Ошибка")
	}
	trainingTypeStr := strings.TrimSpace(parseStr[1])
	if trainingTypeStr == "" {
		return 0, "", 0, errors.New("Ошибка")
	}
	durationStr := strings.TrimSpace(parseStr[2])
	if durationStr == "" {
		return 0, "", 0, errors.New("Ошибка")
	}

	step, err := strconv.Atoi(parseStr[0])
	if err != nil {
		return 0, "", 0, errors.New("Ошибка")
	}

	if step <= 0 {
		return 0, "", 0, errors.New("Ошибка")
	}

	duration, err := time.ParseDuration(parseStr[2])
	if err != nil {
		return 0, "", 0, errors.New("Ошибка")
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("Ошибка")
	}

	return step, parseStr[1], duration, nil
}

func distance(steps int, height float64) float64 {
	if steps <= 0 {
		return 0
	}
	if height <= 0 {
		return 0
	}
	res := height * stepLengthCoefficient
	stepsMulLenSteps := float64(steps) * res
	mKm := stepsMulLenSteps / mInKm
	return mKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 {
		return 0
	}
	if height <= 0 {
		return 0
	}
	if duration <= 0 {
		return 0
	}
	res := distance(steps, height)
	averageSpeed := res / duration.Hours()
	return averageSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)

	if err != nil {
		log.Println(err)
	}
	if steps <= 0 {
		log.Println(err)
	}
	if len(trainingType) <= 0 {
		log.Println(err)
	}
	if duration <= 0 {
		log.Println(err)
	}

	switch trainingType {
	case "Бег":
		cal, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			trainingType, duration.Hours(), distance(steps, height), meanSpeed(steps, height, duration), cal), nil
	case "Ходьба":
		cal, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			trainingType, duration.Hours(), distance(steps, height), meanSpeed(steps, height, duration), cal), nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("не удалось определить колличество шагов %d", steps)
	}
	if height <= 0 {
		return 0, fmt.Errorf("не удалось определить рост %.2f", height)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("не удалось определить вес %.2f", weight)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("не удалось определить продолжительность %d", duration)
	}

	averageSpeed := meanSpeed(steps, height, duration)
	durationMinute := duration.Minutes()
	calories := (weight * averageSpeed * durationMinute) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("не удалось определить колличество шагов %d", steps)
	}
	if height <= 0 {
		return 0, fmt.Errorf("не удалось определить рост %.2f", height)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("не удалось определить вес %.2f", weight)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("не удалось определить продолжительность %d", duration)
	}

	averageSpeed := meanSpeed(steps, height, duration)
	durationMinute := duration.Minutes()
	calories := (weight * averageSpeed * durationMinute) / minInH
	walkingCalories := calories * walkingCaloriesCoefficient

	return walkingCalories, nil
}
