package talkback

import (
	"strconv"
	"strings"
	"time"
)

// SqlFieldTranslation is a translation from a field name to a SQL field.
type SqlFieldTranslation struct {
	Column        string
	Alias         string
	TypeConverter func(value string) (interface{}, error)
}

// SqlTranslations is a map of field names to SQL translations.
type SqlTranslations map[string]SqlFieldTranslation

// ToSqlWhere converts a Query to a SQL WHERE statement.
func ToSqlWhere(query Query, translations SqlTranslations) (string, []interface{}, error) {
	statements := []string{}
	args := []interface{}{}

	translations = sanitizeSqlTranslation(translations)

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

// sanitizeSqlTranslation sanitizes a SqlTranslations map.
func sanitizeSqlTranslation(translations SqlTranslations) SqlTranslations {
	result := SqlTranslations{}

	for field, translation := range translations {
		if translation.Column == "" {
			translation.Column = field
		}

		if translation.Alias == "" {
			translation.Alias = field
		}

		result[field] = translation
	}

	return result
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
	column := translation.Column

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
		return translation.Column + " IN (?)", sliceValue, nil
	case OpNin:
		return translation.Column + " NOT IN (?)", sliceValue, nil
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

// ToSqlSelect converts a Query to a SQL SELECT statement.
func ToSqlSelect(query Query, translations SqlTranslations) (string, error) {
	fields, err := ToSqlSelectSlice(query, translations)
	if err != nil {
		return "", err
	}

	return strings.Join(fields, ", "), nil
}

// ToSqlSelectSlice converts a Query to a slice of SQL SELECT statements.
func ToSqlSelectSlice(query Query, translations SqlTranslations) ([]string, error) {
	fields := []string{}
	qSelects := query.Group
	qSelects = append(qSelects, query.Accumulator...)

	translations = sanitizeSqlTranslation(translations)

	for _, field := range qSelects {
		translation, ok := translations[field]
		if !ok {
			return nil, ErrInvalidField
		}

		col := translation.Column
		if translation.Alias != translation.Column {
			col = col + " AS " + translation.Alias
		}

		fields = append(fields, col)
	}

	return fields, nil
}

// ToSqlGroup converts a Query to a SQL GROUP BY statement.
func ToSqlGroup(query Query, translations SqlTranslations) (string, error) {
	fields, err := ToSqlGroupSlice(query, translations)
	if err != nil {
		return "", err
	}

	return strings.Join(fields, ", "), nil
}

// ToSqlGroupSlice converts a Query to a slice of SQL GROUP BY statements.
func ToSqlGroupSlice(query Query, translations SqlTranslations) ([]string, error) {
	fields := []string{}

	translations = sanitizeSqlTranslation(translations)

	for _, field := range query.Group {
		translation, ok := translations[field]
		if !ok {
			return nil, ErrInvalidField
		}

		fields = append(fields, translation.Column)
	}

	return fields, nil
}

// ToSqlOrderBy converts a Query to a SQL ORDER BY statement.
func ToSqlOrderBy(query Query, translations SqlTranslations) (string, error) {
	fields, err := ToSqlOrderBySlice(query, translations)
	if err != nil {
		return "", err
	}

	return strings.Join(fields, ", "), nil
}

// ToSqlOrderBySlice converts a Query to a slice of SQL ORDER BY statements.
func ToSqlOrderBySlice(query Query, translations SqlTranslations) ([]string, error) {
	fields := []string{}

	translations = sanitizeSqlTranslation(translations)

	for _, field := range query.Sort {
		translation, ok := translations[field.Field]
		if !ok {
			return nil, ErrInvalidField
		}

		col := translation.Column
		if field.Reverse {
			col = col + " DESC"
		} else {
			col = col + " ASC"
		}

		fields = append(fields, col)
	}

	return fields, nil
}

// ToSqlLimit converts a Query to a SQL LIMIT statement.
func ToSqlLimit(query Query) (int, error) {
	return query.Limit, nil
}

// ToSqlOffset converts a Query to a SQL OFFSET statement.
func ToSqlOffset(query Query) (int, error) {
	return query.Skip, nil
}

// SqlPreloadable is a map of preloads (key) and their corresponding model (value).
type SqlPreloadable map[string]string

// ToSqlPreload converts a Query to a SQL preload statement.
func ToSqlPreload(query Query, preloadable SqlPreloadable) ([]string, error) {
	preloads := []string{}

	for _, preload := range query.With {
		model, ok := preloadable[preload]
		if !ok {
			return nil, ErrInvalidPreload
		}

		preloads = append(preloads, model)
	}

	return preloads, nil
}

func ToSql(table string, query Query, translations SqlTranslations) (string, []interface{}, error) {
	cselect, err := ToSqlSelect(query, translations)
	if err != nil {
		return "", nil, err
	}

	cwhere, cwhereargs, err := ToSqlWhere(query, translations)
	if err != nil {
		return "", nil, err
	}

	cgroup, err := ToSqlGroup(query, translations)
	if err != nil {
		return "", nil, err
	}

	corder, err := ToSqlOrderBy(query, translations)
	if err != nil {
		return "", nil, err
	}

	climit, err := ToSqlLimit(query)
	if err != nil {
		return "", nil, err
	}

	coffset, err := ToSqlOffset(query)
	if err != nil {
		return "", nil, err
	}

	sql := "SELECT " + cselect + " FROM " + table +
		" WHERE " + cwhere +
		" GROUP BY " + cgroup +
		" ORDER BY " + corder +
		" LIMIT " + strconv.Itoa(climit) +
		" OFFSET " + strconv.Itoa(coffset)

	return sql, cwhereargs, nil
}

// SqlPlan is a plan for executing a query.
type SqlPlan struct {
	Select    string
	Where     string
	WhereArgs []interface{}
	Group     string
	Order     string
	Limit     int
	Offset    int
	Preload   []string
}

// ToSqlPlan converts a Query to a SqlPlan.
func ToSqlPlan(query Query, translations SqlTranslations, preloadable SqlPreloadable) (SqlPlan, error) {
	cselect, err := ToSqlSelect(query, translations)
	if err != nil {
		return SqlPlan{}, err
	}

	cwhere, cwhereargs, err := ToSqlWhere(query, translations)
	if err != nil {
		return SqlPlan{}, err
	}

	cgroup, err := ToSqlGroup(query, translations)
	if err != nil {
		return SqlPlan{}, err
	}

	corder, err := ToSqlOrderBy(query, translations)
	if err != nil {
		return SqlPlan{}, err
	}

	climit, err := ToSqlLimit(query)
	if err != nil {
		return SqlPlan{}, err
	}

	coffset, err := ToSqlOffset(query)
	if err != nil {
		return SqlPlan{}, err
	}

	cpreload, err := ToSqlPreload(query, preloadable)
	if err != nil {
		return SqlPlan{}, err
	}

	return SqlPlan{
		Select:    cselect,
		Where:     cwhere,
		WhereArgs: cwhereargs,
		Group:     cgroup,
		Order:     corder,
		Limit:     climit,
		Offset:    coffset,
		Preload:   cpreload,
	}, nil
}
