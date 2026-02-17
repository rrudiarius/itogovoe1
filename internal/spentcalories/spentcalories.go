package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// 1. Разделить строку на слайс строк
	parts := strings.Split(data, ",")
	
	// 2. Проверить, чтобы длина слайса была равна 3
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных, ожидается: 'шаги,вид_активности,длительность'")
	}
	
	// 3. Преобразовать первый элемент (количество шагов) в int
	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования шагов '%s': %v", stepsStr, err)
	}
	
	// 4. Получить вид активности
	activityType := strings.TrimSpace(parts[1])
	
	// 5. Преобразовать третий элемент в time.Duration
	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования длительности '%s': %v", durationStr, err)
	}
	
	// 6. Проверить корректность значений
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0, получено: %d", steps)
	}
	
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("длительность должна быть больше 0, получено: %v", duration)
	}
	
	// 7. Вернуть результат
	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	// 1. Рассчитать длину шага
	stepLength := height * stepLengthCoefficient
	
	// 2. Умножить количество шагов на длину шага
	distanceMeters := float64(steps) * stepLength
	
	// 3. Разделить на количество метров в километре
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// 1. Проверить, что продолжительность больше 0
	if duration <= 0 {
		return 0
	}
	
	// 2. Вычислить дистанцию
	dist := distance(steps, height)
	
	// 3. Вычислить продолжительность в часах
	durationHours := duration.Hours()
	
	// 4. Проверить деление на 0
	if durationHours == 0 {
		return 0
	}
	
	// 5. Вычислить среднюю скорость
	return dist / durationHours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. Проверить входные параметры на корректность
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	
	if duration <= 0 {
		return 0, fmt.Errorf("длительность должна быть больше 0")
	}
	
	// 2. Рассчитать среднюю скорость
	meanSpeedValue := meanSpeed(steps, height, duration)
	
	// 3. Рассчитать и вернуть количество калорий
	// Перевести продолжительность в минуты
	durationMinutes := duration.Minutes()
	
	// Формула: (вес * средняя_скорость * длительность_в_минутах) / минут_в_часе
	calories := (weight * meanSpeedValue * durationMinutes) / minInH
	
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. Проверить входные параметры на корректность
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	
	if duration <= 0 {
		return 0, fmt.Errorf("длительность должна быть больше 0")
	}
	
	// 2. Рассчитать среднюю скорость
	meanSpeedValue := meanSpeed(steps, height, duration)
	
	// 3. Рассчитать количество калорий
	// Перевести продолжительность в минуты
	durationMinutes := duration.Minutes()
	
	// Базовая формула
	calories := (weight * meanSpeedValue * durationMinutes) / minInH
	
	// 4. Применить корректирующий коэффициент
	calories *= walkingCaloriesCoefficient
	
	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// 1. Получить значения из строки данных
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка парсинга:", err)
		return "", err
	}
	
	// 2. Вычислить дистанцию
	dist := distance(steps, height)
	
	// 3. Вычислить среднюю скорость
	speed := meanSpeed(steps, height, duration)
	
	// 4. В зависимости от типа активности рассчитать калории
	var calories float64
	var caloriesErr error
	
	// Приводим тип активности к нижнему регистру для удобства сравнения
	activityLower := strings.ToLower(strings.TrimSpace(activityType))
	
	switch activityLower {
	case "бег", "running", "run":
		calories, caloriesErr = RunningSpentCalories(steps, weight, height, duration)
	case "ходьба", "walking", "walk":
		calories, caloriesErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activityType)
	}
	
	if caloriesErr != nil {
		log.Println("Ошибка вычисления калорий:", caloriesErr)
		return "", caloriesErr
	}
	
	// 5. Сформировать строку результата
	result := fmt.Sprintf("Тип тренировки: %s\n", activityType)
	result += fmt.Sprintf("Длительность: %.2f ч.\n", duration.Hours())
	result += fmt.Sprintf("Дистанция: %.2f км.\n", dist)
	result += fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
	result += fmt.Sprintf("Сожгли калорий: %.2f", calories)
	
	return result, nil
}
