package talkback

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToSqlWhere(t *testing.T) {
	sampleDate := "2007-12-14"
	parsedSampleDate, err := time.Parse("2006-01-02", sampleDate)
	if err != nil {
		t.Error(err)
	}

	sampleTime := "12:07:07"
	parsedSampleTime, err := time.Parse("15:04:05", sampleTime)
	if err != nil {
		t.Error(err)
	}

	sampleDateTime := "2007-12-14 12:07:07"
	parsedSampleDateTime, err := time.Parse("2006-01-02 15:04:05", sampleDateTime)
	if err != nil {
		t.Error(err)
	}

	sampleISO8601 := "2007-12-14T12:07:07Z"
	parsedSampleISO8601, err := time.Parse(time.RFC3339, sampleISO8601)
	if err != nil {
		t.Error(err)
	}

	type scenarioT struct {
		query        Query
		translations SqlTranslations
		statement    string
		args         []interface{}
		err          error
	}

	scenarios := []scenarioT{
		{
			query: Query{
				Conditions: []Condition{
					{"field1", "eq", []string{"value1"}},
					{"field2", "ne", []string{sampleTime}},
					{"field3", "isnull", []string{"true"}},
				},
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{
					Column:        "field1",
					TypeConverter: SqlConvertString,
				},
				"field2": SqlFieldTranslation{
					Column:        "field2",
					TypeConverter: SqlConvertTime,
				},
				"field3": SqlFieldTranslation{
					Column: "field3",
				},
			},
			statement: "field1 = ? AND field2 != ? AND field3 IS NULL",
			args:      []interface{}{"value1", parsedSampleTime},
			err:       nil,
		},
		{
			query: Query{
				Conditions: []Condition{
					{"field1", "gt", []string{"12"}},
					{"field2", "lt", []string{sampleDate}},
					{"field3", "gte", []string{"true"}},
					{"field4", "lte", []string{"7.2"}},
				},
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{
					Column:        "field1",
					TypeConverter: SqlConvertInt,
				},
				"field2": SqlFieldTranslation{
					Column:        "field2",
					TypeConverter: SqlConvertDate,
				},
				"field3": SqlFieldTranslation{
					Column:        "field3",
					TypeConverter: SqlConvertBool,
				},
				"field4": SqlFieldTranslation{
					Column:        "field4",
					TypeConverter: SqlConvertFloat,
				},
			},
			statement: "field1 > ? AND field2 < ? AND field3 >= ? AND field4 <= ?",
			args: []interface{}{
				int(12),
				parsedSampleDate,
				true,
				float64(7.2),
			},
			err: nil,
		},
		{
			query: Query{
				Conditions: []Condition{
					{"field1", "contain", []string{"value1"}},
					{"field2", "ncontain", []string{"value2"}},
					{"field3", "contains", []string{"value3"}},
					{"field4", "ncontains", []string{"value4"}},
				},
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{
					Column: "field1",
				},
				"field2": SqlFieldTranslation{
					Column: "field2",
				},
				"field3": SqlFieldTranslation{
					Column: "field3",
				},
				"field4": SqlFieldTranslation{
					Column: "field4",
				},
			},
			statement: "CAST(field1 AS TEXT) ILIKE ? AND CAST(field2 AS TEXT) NOT ILIKE ? AND CAST(field3 AS TEXT) LIKE ? AND CAST(field4 AS TEXT) NOT LIKE ?",
			args: []interface{}{
				"%" + "value1" + "%",
				"%" + "value2" + "%",
				"%" + "value3" + "%",
				"%" + "value4" + "%",
			},
			err: nil,
		},
		{
			query: Query{
				Conditions: []Condition{
					{"field1", "in", []string{sampleDateTime, sampleDateTime}},
					{"field2", "nin", []string{sampleISO8601, sampleISO8601}},
				},
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{
					Column:        "field1",
					TypeConverter: SqlConvertDateTime,
				},
				"field2": SqlFieldTranslation{
					Column:        "field2",
					TypeConverter: SqlConvertISO8601,
				},
			},
			statement: "field1 IN (?) AND field2 NOT IN (?)",
			args: []interface{}{
				[]interface{}{parsedSampleDateTime, parsedSampleDateTime},
				[]interface{}{parsedSampleISO8601, parsedSampleISO8601},
			},
			err: nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.statement, func(t *testing.T) {
			statement, args, err := ToSqlWhere(scenario.query, scenario.translations)

			assert.Equal(t, scenario.statement, statement, "statement should be equal")
			assert.Equal(t, scenario.args, args, "args should be equal")
			assert.Equal(t, scenario.err, err, "err should be equal")
		})
	}
}

func TestToSqlSelect(t *testing.T) {
	type scenarioT struct {
		query        Query
		translations SqlTranslations
		statement    string
		err          error
	}

	scenarios := []scenarioT{
		{
			query: Query{
				Group: []string{"field1", "field2"},
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{},
				"field2": SqlFieldTranslation{},
			},
			statement: "field1, field2",
			err:       nil,
		},
		{
			query: Query{
				Group: []string{"field1", "field2"},
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{
					Column: "ex.field1",
				},
				"field2": SqlFieldTranslation{
					Column: "ex.field2",
				},
			},
			statement: "ex.field1 AS field1, ex.field2 AS field2",
			err:       nil,
		},
		{
			query: Query{
				Accumulator: []string{"field1", "field2"},
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{
					Alias: "alias1",
				},
				"field2": SqlFieldTranslation{
					Alias: "alias2",
				},
			},
			statement: "field1 AS alias1, field2 AS alias2",
			err:       nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.statement, func(t *testing.T) {
			statement, err := ToSqlSelect(scenario.query, scenario.translations)

			assert.Equal(t, scenario.statement, statement, "statement should be equal")
			assert.Equal(t, scenario.err, err, "err should be equal")
		})
	}
}

func TestToSqlGroup(t *testing.T) {
	type scenarioT struct {
		query        Query
		translations SqlTranslations
		statement    string
		err          error
	}

	scenarios := []scenarioT{
		{
			query: Query{
				Group: []string{"field1", "field2"},
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{},
				"field2": SqlFieldTranslation{},
			},
			statement: "field1, field2",
			err:       nil,
		},
		{
			query: Query{
				Group: []string{"field1", "field2"},
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{
					Column: "ex.field1",
				},
				"field2": SqlFieldTranslation{
					Column: "ex.field2",
				},
			},
			statement: "ex.field1, ex.field2",
			err:       nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.statement, func(t *testing.T) {
			statement, err := ToSqlGroup(scenario.query, scenario.translations)

			assert.Equal(t, scenario.statement, statement, "statement should be equal")
			assert.Equal(t, scenario.err, err, "err should be equal")
		})
	}
}

func TestToSqlOrderBy(t *testing.T) {
	type scenarioT struct {
		query        Query
		translations SqlTranslations
		statement    string
		err          error
	}

	scenarios := []scenarioT{
		{
			query: Query{
				Sort: []Sort{
					{"field1", false},
					{"field2", true},
				},
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{},
				"field2": SqlFieldTranslation{},
			},
			statement: "field1 ASC, field2 DESC",
			err:       nil,
		},
		{
			query: Query{
				Sort: []Sort{
					{"field1", false},
					{"field2", true},
				},
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{
					Column: "ex.field1",
				},
				"field2": SqlFieldTranslation{
					Column: "ex.field2",
				},
			},
			statement: "ex.field1 ASC, ex.field2 DESC",
			err:       nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.statement, func(t *testing.T) {
			statement, err := ToSqlOrderBy(scenario.query, scenario.translations)

			assert.Equal(t, scenario.statement, statement, "statement should be equal")
			assert.Equal(t, scenario.err, err, "err should be equal")
		})
	}
}

func TestToSqlLimit(t *testing.T) {
	type scenarioT struct {
		query     Query
		statement int
		err       error
	}

	scenarios := []scenarioT{
		{
			query: Query{
				Limit: 10,
			},
			statement: 10,
			err:       nil,
		},
		{
			query: Query{
				Limit: 0,
			},
			statement: 0,
			err:       nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(strconv.Itoa(scenario.statement), func(t *testing.T) {
			statement, err := ToSqlLimit(scenario.query)

			assert.Equal(t, scenario.statement, statement, "statement should be equal")
			assert.Equal(t, scenario.err, err, "err should be equal")
		})
	}
}

func TestToSqlOffset(t *testing.T) {
	type scenarioT struct {
		query     Query
		statement int
		err       error
	}

	scenarios := []scenarioT{
		{
			query: Query{
				Skip: 10,
			},
			statement: 10,
			err:       nil,
		},
		{
			query: Query{
				Skip: 0,
			},
			statement: 0,
			err:       nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(strconv.Itoa(scenario.statement), func(t *testing.T) {
			statement, err := ToSqlOffset(scenario.query)

			assert.Equal(t, scenario.statement, statement, "statement should be equal")
			assert.Equal(t, scenario.err, err, "err should be equal")
		})
	}
}

func TestToSqlPreload(t *testing.T) {
	type scenarioT struct {
		query       Query
		preloadable SqlPreloadable
		preloads    []string
		err         error
	}

	scenarios := []scenarioT{
		{
			query: Query{
				With: []string{"field1", "field2"},
			},
			preloadable: SqlPreloadable{
				"field1": "Field1",
				"field2": "Field2",
			},
			preloads: []string{"Field1", "Field2"},
			err:      nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(strings.Join(scenario.preloads, ","), func(t *testing.T) {
			preloads, err := ToSqlPreload(scenario.query, scenario.preloadable)

			assert.Equal(t, scenario.preloads, preloads, "preloads should be equal")
			assert.Equal(t, scenario.err, err, "err should be equal")
		})
	}
}

func TestToSql(t *testing.T) {
	type scenarioT struct {
		query        Query
		translations SqlTranslations
		statement    string
		args         []interface{}
		err          error
	}

	scenarios := []scenarioT{
		{
			query: Query{
				Conditions: []Condition{
					{"field1", "eq", []string{"value1"}},
					{"field2", "ne", []string{"value2"}},
				},
				Group: []string{"field1", "field2"},
				Sort: []Sort{
					{"field1", false},
					{"field2", true},
				},
				Limit: 10,
				Skip:  10,
			},
			translations: SqlTranslations{
				"field1": SqlFieldTranslation{
					Alias: "alias1",
				},
				"field2": SqlFieldTranslation{
					Column: "ex.field2",
				},
			},
			statement: "SELECT field1 AS alias1, ex.field2 AS field2 FROM somewhere WHERE field1 = ? AND ex.field2 != ? GROUP BY field1, ex.field2 ORDER BY field1 ASC, ex.field2 DESC LIMIT 10 OFFSET 10",
			args:      []interface{}{"value1", "value2"},
			err:       nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.statement, func(t *testing.T) {
			statement, args, err := ToSql("somewhere", scenario.query, scenario.translations)

			assert.Equal(t, scenario.statement, statement, "statement should be equal")
			assert.Equal(t, scenario.args, args, "args should be equal")
			assert.Equal(t, scenario.err, err, "err should be equal")
		})
	}
}
