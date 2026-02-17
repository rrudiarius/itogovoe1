package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// 1. Разделить строку на слайс строк
	parts := strings.Split(data, ",")
	
	// 2. Проверить, чтобы длина слайса была равна 2
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных, ожидается: 'шаги,длительность'")
	}
	
	// 3. Преобразовать первый элемент (количество шагов) в int
	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования шагов: %v", err)
	}
	
	// 4. Проверить: количество шагов должно быть больше 0
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0, получено: %d", steps)
	}
	
	// 5. Преобразовать второй элемент в time.Duration
	durationStr := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования длительности: %v", err)
	}
	
	// 6. Проверить, что длительность > 0
	if duration <= 0 {
		return 0, 0, fmt.Errorf("длительность должна быть больше 0, получено: %v", duration)
	}
	
	// 7. Вернуть результат
	return steps, duration, nil
}

// Для вызова WalkingSpentCalories из другого пакета
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Эта функция будет реализована в пакете spentcalories
	// Сейчас вернем заглушку
	return 0, fmt.Errorf("функция WalkingSpentCalories не реализована")
}

func DayActionInfo(data string, weight, height float64) string {
	// 1. Получить данные о количестве шагов и продолжительности
	steps, duration, err := parsePackage(data)
	if err != nil {
		// Вывести ошибку на экран (в реальном приложении лучше использовать log)
		fmt.Printf("Ошибка: %v\n", err)
		return ""
	}
	
	// 2. Проверить, чтобы количество шагов было больше 0
	// (уже проверено в parsePackage, но можно проверить еще раз)
	if steps <= 0 {
		return ""
	}
	
	// 3. Вычислить дистанцию в метрах
	distanceMeters := float64(steps) * stepLength
	
	// 4. Перевести дистанцию в километры
	distanceKm := distanceMeters / mInKm
	
	// 5. Вычислить количество калорий
	// Используем временную реализацию, пока функция не будет реализована в spentcalories
	calories := 0.0
	if walkingCalories, err := WalkingSpentCalories(steps, weight, height, duration); err == nil {
		calories = walkingCalories
	} else {
		// Если функция не реализована, используем упрощенный расчет
		// Формула: вес * дистанция_в_км * коэффициент 0.78
		calories = weight * distanceKm * 0.78
	}
	
	// 6. Сформировать строку результата
	result := fmt.Sprintf("Количество шагов: %d.\n", steps)
	result += fmt.Sprintf("Дистанция составила %.2f км.\n", distanceKm)
	result += fmt.Sprintf("Вы сожгли %.2f ккал.", calories)
	
	return result
}
