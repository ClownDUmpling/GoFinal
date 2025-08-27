package api

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

// Определяем високосный ли год
func isLeap(year int) bool {
	return year%400 == 0 || (year%100 != 0 && year%4 == 0)
}

func NextDate(now time.Time, dateStr, repeat string) (string, error) {

	//Проверка на пустое правило
	if repeat == "" {
		return "", ErrInvalidRule
	}

	//Парсим начальную дату
	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", ErrInvalidDate
	}

	//Разбиваем правило на кусочки
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", ErrInvalidRule
	}

	//Определяем тип правила
	switch parts[0] {
	case "d":
		return handleDaily(now, date, parts)
	case "y":
		return handleYearly(now, date)
	case "w":
		return handleWeekly(now, date, parts)
	case "m":
		return handleMonthly(now, date, parts)
	default:
		return "", ErrUnsupported
	}
}

// Работаем с правилом d
func handleDaily(now, date time.Time, parts []string) (string, error) {
	// Проверка на наличие количества дней
	if len(parts) != 2 {
		return "", ErrInvalidRule
	}
	// Проверка на правильное количество
	days, err := strconv.Atoi(parts[1])
	if err != nil || days <= 0 || days > 400 {
		return "", ErrInvalidDayInterval
	}

	// Вычисляем следующую дату
	nextDate := date.AddDate(0, 0, days)
	for !afterNow(nextDate, now) {
		nextDate = nextDate.AddDate(0, 0, days)
	}

	return nextDate.Format(DateFormat), nil
}

// Работаем с правилом y
func handleYearly(now, date time.Time) (string, error) {
	// Запоминаем была ли изначальная дата високосным годом
	isLeapDay := date.Month() == time.February && date.Day() == 29
	nextDate := date.AddDate(1, 0, 0)

	// Вычисляем следующую дату
	for !afterNow(nextDate, now) {
		nextDate = nextDate.AddDate(1, 0, 0)

		if isLeapDay {
			if isLeap(nextDate.Year()) {
				// Вручную устанавливаем 29 февраля для високосных
				nextDate = time.Date(nextDate.Year(), time.February, 29, 0, 0, 0, 0, time.UTC)
			} else {
				// Вручную устанавливаем 1 марта для не високосных
				nextDate = time.Date(nextDate.Year(), time.March, 1, 0, 0, 0, 0, time.UTC)
			}
		}
	}

	return nextDate.Format(DateFormat), nil
}

// Работаем с правилом w
func handleWeekly(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", ErrInvalidRule
	}

	// Парсим дни недели
	daysStr := strings.Split(parts[1], ",")
	var weekdays [8]bool // 1-7 (пн-вс), 0 не используется

	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 7 {
			return "", ErrInvalidDay
		}
		weekdays[day] = true
	}

	// Начинаем с текущей даты или переданной даты (если она в будущем)
	current := date
	for i := 0; i < 730; i++ {
		weekday := int(current.Weekday())

		// Преобразуем воскресенье (0) в 7
		if weekday == 0 {
			weekday = 7
		}

		// Проверяем, подходит ли день недели и дата в будущем
		if weekdays[weekday] {

			formattedDate := current.Format(DateFormat)
			parsedDate, err := time.Parse(DateFormat, formattedDate)
			if err != nil {
				return "", err
			}

			if afterNow(parsedDate, now) {
				return formattedDate, nil
			}
		}

		// Переходим к следующему дню
		current = current.AddDate(0, 0, 1)
	}
	return "", ErrInvalidDate
}

// Работаем с правилом m
func handleMonthly(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", ErrInvalidRule
	}

	// Парсим дни месяца
	days, err := parseMonthDays(parts[1])
	if err != nil {
		return "", err
	}

	// Парсим месяцы (если указаны)
	var months []int
	if len(parts) >= 3 {
		months, err = parseMonths(parts[2])
		if err != nil {
			return "", err
		}
	}

	// Создаем массивы допустимых значений
	var validDays [32]bool
	for _, day := range days {
		if day > 0 {
			validDays[day] = true
		}
	}

	var validMonths [13]bool
	if len(months) > 0 {
		for _, month := range months {
			validMonths[month] = true
		}
	} else {
		for i := 1; i <= 12; i++ {
			validMonths[i] = true
		}
	}

	// Начинаем со следующего дня после переданной даты
	current := date

	// Перебираем дни по одному
	for i := 0; i < 1000; i++ {
		current = current.AddDate(0, 0, 1)
		currentDay := current.Day()
		currentMonth := int(current.Month())

		// Проверяем специальные дни
		lastDay := time.Date(current.Year(), current.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		isLastDay := currentDay == lastDay
		isPreLastDay := currentDay == lastDay-1

		// Проверяем условия
		dayValid := validDays[currentDay] ||
			(contains(days, -1) && isLastDay) ||
			(contains(days, -2) && isPreLastDay)

		monthValid := validMonths[currentMonth]

		if dayValid && monthValid && afterNow(current, now) {
			return current.Format(DateFormat), nil
		}
	}

	return "", ErrInvalidDate
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Парсим дни
func parseMonthDays(daysStr string) ([]int, error) {
	parts := strings.Split(daysStr, ",")
	days := make([]int, 0, len(parts))

	for _, part := range parts {
		day, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, ErrInvalidDay
		}

		// Проверяем допустимость дня
		if (day < -2 || day == 0 || day > 31) && !(day == -1 || day == -2) {
			return nil, ErrInvalidDay
		}
		days = append(days, day)
	}
	sort.Ints(days)
	return days, nil
}

// Парсим месяцы
func parseMonths(monthsStr string) ([]int, error) {
	parts := strings.Split(monthsStr, ",")
	months := make([]int, 0, len(parts))

	for _, part := range parts {
		month, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, ErrInvalidMonth
		}
		if month < 1 || month > 12 {
			return nil, ErrInvalidMonth
		}
		months = append(months, month)
	}
	sort.Ints(months)
	return months, nil
}

// afterNow проверяет, что первая дата (без времени) больше второй даты (без времени)
func afterNow(date, now time.Time) bool {
	dateDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	result := dateDate.After(nowDate)
	return result
}
