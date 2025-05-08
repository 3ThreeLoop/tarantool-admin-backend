package sqlx

import (
	"api-mini-shop/pkg/share"
	"fmt"
	"os"
	"strings"
	"time"
)

// built sql filter
func BuildSQLFilter(filters []share.Filter) (string, []interface{}) {
	var clauses []string
	var params []interface{}
	placeholder := 1

	// Get current OS time
	app_timezone := os.Getenv("APP_TIMEZONE")
	if app_timezone == "" {
		app_timezone = "Asia/Phnom_Penh"
	}
	location, err := time.LoadLocation(app_timezone)
	if err != nil {
		return "", nil
	}

	// convert value types
	convertToTime := func(v interface{}) (time.Time, bool) {
		switch val := v.(type) {
		case string:
			parsed, err := time.ParseInLocation("2006-01-02", val, location)
			if err == nil {
				return parsed, true
			}
		case time.Time:
			return val, true
		}
		return time.Time{}, false
	}

	convertToBool := func(v interface{}) (bool, bool) {
		switch val := v.(type) {
		case bool:
			return val, true
		case string:
			if val == "true" {
				return true, true
			} else if val == "false" {
				return false, true
			}
		}
		return false, false
	}

	for _, f := range filters {
		field := f.Property
		op := strings.ToLower(f.Operator)

		switch op {
		case "eq", "neq", "lt", "lte", "gt", "gte":
			sqlOp := map[string]string{
				"eq":  "=",
				"neq": "!=",
				"lt":  "<",
				"lte": "<=",
				"gt":  ">",
				"gte": ">=",
			}[op]

			// try to convert to time
			if t, ok := convertToTime(f.Value); ok {
				f.Value = t
			} else if b, ok := convertToBool(f.Value); ok {
				f.Value = b
			}

			clauses = append(clauses, fmt.Sprintf("%s %s $%d", field, sqlOp, placeholder))
			params = append(params, f.Value)
			placeholder++

		case "like":
			clauses = append(clauses, fmt.Sprintf("%s LIKE $%d", field, placeholder))
			params = append(params, f.Value)
			placeholder++

		case "in":
			vals, ok := f.Value.([]interface{})
			if !ok || len(vals) == 0 {
				continue
			}
			var ph []string
			for _, v := range vals {
				ph = append(ph, fmt.Sprintf("$%d", placeholder))
				params = append(params, v)
				placeholder++
			}
			clauses = append(clauses, fmt.Sprintf("%s IN (%s)", field, strings.Join(ph, ", ")))

		case "between":
			vals, ok := f.Value.([]interface{})
			if !ok || len(vals) != 2 {
				continue
			}

			start, ok1 := convertToTime(vals[0])
			end, ok2 := convertToTime(vals[1])

			if ok1 && ok2 {
				// Make end date inclusive to the end of the day
				end = end.Add(24 * time.Hour).Add(-time.Second)
				clauses = append(clauses, fmt.Sprintf("%s BETWEEN $%d AND $%d", field, placeholder, placeholder+1))
				params = append(params, start, end)
				placeholder += 2
			} else {
				// Fallback: use original values
				clauses = append(clauses, fmt.Sprintf("%s BETWEEN $%d AND $%d", field, placeholder, placeholder+1))
				params = append(params, vals[0], vals[1])
				placeholder += 2
			}
		}
	}

	if len(clauses) == 0 {
		return "", nil
	}
	return strings.Join(clauses, " AND "), params
}

// built sql sort
func BuildSort(sorts []share.Sort) (string, []interface{}) {
	var orderClauses []string
	var params []interface{}

	for _, sort := range sorts {
		field := sort.Property
		direction := strings.ToUpper(sort.Direction)

		// ensure the direction is either ASC or DESC
		if direction != "ASC" && direction != "DESC" {
			direction = "ASC"
		}

		// add the sort clause to the list of order clauses
		orderClauses = append(orderClauses, fmt.Sprintf("%s %s", field, direction))
	}

	if len(orderClauses) == 0 {
		return "", nil
	}

	// join the clauses with commas and return the final order by string
	return "ORDER BY " + strings.Join(orderClauses, ", "), params
}

// built sql paging
func BuildPaging(page int, perPage int) string {
	var params []interface{}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	limit := perPage

	params = append(params, offset, limit)

	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}
