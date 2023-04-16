package talkback

import (
	"strconv"
	"strings"
	"time"
)

// SqlFieldTranslation is a translation from a field name to a SQL field.
type SqlFieldTranslation struct {
	Condition     string
	Alias         string
	TypeConverter func(value string) (interface{}, error)
}

// SqlTranslations is a map of field names to SQL translations.
type SqlTranslations map[string]SqlFieldTranslation

// ToSqlWhere converts a Query to a SQL WHERE statement.
func ToSqlWhere(query Query, translations SqlTranslations) (string, []interface{}, error) {
	statements := []string{}
	args := []interface{}{}

	for _, cond := range query.Conditions {
		translation, ok := translations[cond.Field]
		if !ok {
			return "", nil, ErrInvalidField
		}

		statement, arg, err := conditionToSql(translation, cond)
		if err != nil {
			return "", nil, err
		}

		statements = append(statements, statement)

		if arg != nil {
			args = append(args, arg)
		}
	}

	return strings.Join(statements, " AND "), args, nil
}

// valueToSql converts a value to a SQL argument.
func valueToSql(translation SqlFieldTranslation, value string) (interface{}, error) {
	var result interface{} = value

	if translation.TypeConverter != nil {
		var err error

		result, err = translation.TypeConverter(value)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// sliceValuesToSql converts a slice of values to a slice of SQL arguments.
func sliceValuesToSql(translation SqlFieldTranslation, values []string) ([]interface{}, error) {
	var result []interface{}

	for _, value := range values {
		r, err := valueToSql(translation, value)
		if err != nil {
			return nil, err
		}

		result = append(result, r)
	}

	return result, nil
}

// conditionToSql converts a Condition to a SQL statement and argument.
func conditionToSql(translation SqlFieldTranslation, cond Condition) (string, interface{}, error) {
	sliceValue, err := sliceValuesToSql(translation, cond.Values)
	if err != nil {
		return "", nil, err
	}

	firstValue := sliceValue[0]
	likeValue := "%" + cond.Values[0] + "%"
	column := translation.Condition

	switch cond.Op {
	case OpIsNull:
		return column + " IS NULL", nil, nil
	case OpEq:
		return column + " = ?", firstValue, nil
	case OpNe:
		return column + " != ?", firstValue, nil
	case OpGt:
		return column + " > ?", firstValue, nil
	case OpGte:
		return column + " >= ?", firstValue, nil
	case OpLt:
		return column + " < ?", firstValue, nil
	case OpLte:
		return column + " <= ?", firstValue, nil
	case OpContain:
		return castAsText(column) + " ILIKE ?", likeValue, nil
	case OpNcontain:
		return castAsText(column) + " NOT ILIKE ?", likeValue, nil
	case OpContains:
		return castAsText(column) + " LIKE ?", likeValue, nil
	case OpNcontains:
		return castAsText(column) + " NOT LIKE ?", likeValue, nil
	case OpIn:
		return translation.Condition + " IN (?)", sliceValue, nil
	case OpNin:
		return translation.Condition + " NOT IN (?)", sliceValue, nil
	default:
		return "", nil, ErrInvalidOp
	}
}

// castAsText casts a field as text.
func castAsText(field string) string {
	return "CAST(" + field + " AS TEXT)"
}

// SqlConvertString is a TypeConverter that converts a string to a string.
func SqlConvertString(value string) (interface{}, error) {
	return value, nil
}

// SqlConvertString is a TypeConverter that converts a string to a string.
func SqlConvertInt(value string) (interface{}, error) {
	return strconv.Atoi(value)
}

// SqlConvertString is a TypeConverter that converts a string to a string.
func SqlConvertBool(value string) (interface{}, error) {
	return strconv.ParseBool(value)
}

// SqlConvertString is a TypeConverter that converts a string to a string.
func SqlConvertFloat(value string) (interface{}, error) {
	return strconv.ParseFloat(value, 64)
}

// SqlConvertString is a TypeConverter that converts a string to a string.
func SqlConvertDate(value string) (interface{}, error) {
	return time.Parse("2006-01-02", value)
}

// SqlConvertString is a TypeConverter that converts a string to a string.
func SqlConvertDateTime(value string) (interface{}, error) {
	return time.Parse("2006-01-02 15:04:05", value)
}

// SqlConvertString is a TypeConverter that converts a string to a string.
func SqlConvertTime(value string) (interface{}, error) {
	return time.Parse("15:04:05", value)
}

// SqlConvertString is a TypeConverter that converts a string to a string.
func SqlConvertISO8601(value string) (interface{}, error) {
	return time.Parse(time.RFC3339, value)
}
