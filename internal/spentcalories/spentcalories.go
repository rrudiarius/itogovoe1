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
		return 0, "", 0, errors.New("invalid data format: expected 3 fields separated by comma")
	}
	stepsStr := strings.TrimSpace(parseStr[0])
	if stepsStr == "" {
		return 0, "", 0, errors.New("steps field is empty")
	}
	trainingTypeStr := strings.TrimSpace(parseStr[1])
	if trainingTypeStr == "" {
		return 0, "", 0, errors.New("training type field is empty")
	}
	durationStr := strings.TrimSpace(parseStr[2])
	if durationStr == "" {
		return 0, "", 0, errors.New("duration field is empty")
	}

	step, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps value: %w", err)
	}

	if step <= 0 {
		return 0, "", 0, fmt.Errorf("steps must be positive, got: %d", step)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration format: %w", err)
	}

	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("duration must be positive, got: %v", duration)
	}

	return step, trainingTypeStr, duration, nil
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
		log.Printf("error parsing training data: %v", err)
		return "", err
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
		return "", errors.New("unknown training type")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid steps count: %d", steps)
	}
	if height <= 0 {
		return 0, fmt.Errorf("invalid height: %.2f", height)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight: %.2f", weight)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration: %v", duration)
	}

	averageSpeed := meanSpeed(steps, height, duration)
	durationMinute := duration.Minutes()
	calories := (weight * averageSpeed * durationMinute) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid steps count: %d", steps)
	}
	if height <= 0 {
		return 0, fmt.Errorf("invalid height: %.2f", height)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight: %.2f", weight)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration: %v", duration)
	}

	averageSpeed := meanSpeed(steps, height, duration)
	durationMinute := duration.Minutes()
	calories := (weight * averageSpeed * durationMinute) / minInH
	walkingCalories := calories * walkingCaloriesCoefficient

	return walkingCalories, nil
}
