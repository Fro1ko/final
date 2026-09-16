package api

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()
	if value := r.FormValue("now"); value != "" {
		var err error
		now, err = time.Parse(dateFormat, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte(next))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("repeat rule is empty")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid day rule")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("invalid day interval")
		}
		for {
			start = start.AddDate(0, 0, days)
			if start.After(now) {
				return start.Format(dateFormat), nil
			}
		}
	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("invalid year rule")
		}
		for {
			start = start.AddDate(1, 0, 0)
			if start.After(now) {
				return start.Format(dateFormat), nil
			}
		}
	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid weekday rule")
		}
		weekdays, err := parseWeekdays(parts[1])
		if err != nil {
			return "", err
		}
		for candidate := start; ; candidate = candidate.AddDate(0, 0, 1) {
			weekday := int(candidate.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			if weekdays[weekday] && candidate.After(now) {
				return candidate.Format(dateFormat), nil
			}
		}
	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("invalid month rule")
		}
		days, err := parseMonthDays(parts[1])
		if err != nil {
			return "", err
		}
		months := map[int]bool{}
		if len(parts) == 3 {
			months, err = parseMonths(parts[2])
			if err != nil {
				return "", err
			}
		}
		if !hasValidMonthDate(days, months) {
			return "", fmt.Errorf("month rule has no valid dates")
		}
		for month := firstMonth(start); ; month = month.AddDate(0, 1, 0) {
			if len(months) > 0 && !months[int(month.Month())] {
				continue
			}
			candidates := make([]time.Time, 0, len(days))
			for _, day := range days {
				candidate, ok := monthDay(month, day)
				if ok {
					candidates = append(candidates, candidate)
				}
			}
			sort.Slice(candidates, func(i, j int) bool {
				return candidates[i].Before(candidates[j])
			})
			for _, candidate := range candidates {
				if candidate.After(now) && !candidate.Before(start) {
					return candidate.Format(dateFormat), nil
				}
			}
		}
	default:
		return "", fmt.Errorf("unsupported repeat rule")
	}
}

func parseWeekdays(value string) (map[int]bool, error) {
	result := map[int]bool{}
	for _, item := range strings.Split(value, ",") {
		day, err := strconv.Atoi(item)
		if err != nil || day < 1 || day > 7 {
			return nil, fmt.Errorf("invalid weekday")
		}
		result[day] = true
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("weekday list is empty")
	}
	return result, nil
}

func parseMonthDays(value string) ([]int, error) {
	result := make([]int, 0)
	for _, item := range strings.Split(value, ",") {
		day, err := strconv.Atoi(item)
		if err != nil || (day < 1 && day != -1 && day != -2) || day > 31 {
			return nil, fmt.Errorf("invalid month day")
		}
		result = append(result, day)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("month day list is empty")
	}
	return result, nil
}

func parseMonths(value string) (map[int]bool, error) {
	result := map[int]bool{}
	for _, item := range strings.Split(value, ",") {
		month, err := strconv.Atoi(item)
		if err != nil || month < 1 || month > 12 {
			return nil, fmt.Errorf("invalid month")
		}
		result[month] = true
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("month list is empty")
	}
	return result, nil
}

func firstMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

func monthDay(month time.Time, day int) (time.Time, bool) {
	lastDay := time.Date(month.Year(), month.Month()+1, 0, 0, 0, 0, 0, month.Location()).Day()
	if day < 0 {
		day = lastDay + day + 1
	}
	if day < 1 || day > lastDay {
		return time.Time{}, false
	}
	return time.Date(month.Year(), month.Month(), day, 0, 0, 0, 0, month.Location()), true
}

func hasValidMonthDate(days []int, months map[int]bool) bool {
	for month := 1; month <= 12; month++ {
		if len(months) > 0 && !months[month] {
			continue
		}
		for _, day := range days {
			if _, ok := monthDay(time.Date(2024, time.Month(month), 1, 0, 0, 0, 0, time.UTC), day); ok {
				return true
			}
		}
	}
	return false
}
