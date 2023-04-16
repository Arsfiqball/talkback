package talkback

import (
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
					Condition:     "field1",
					TypeConverter: SqlConvertString,
				},
				"field2": SqlFieldTranslation{
					Condition:     "field2",
					TypeConverter: SqlConvertTime,
				},
				"field3": SqlFieldTranslation{
					Condition: "field3",
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
					Condition:     "field1",
					TypeConverter: SqlConvertInt,
				},
				"field2": SqlFieldTranslation{
					Condition:     "field2",
					TypeConverter: SqlConvertDate,
				},
				"field3": SqlFieldTranslation{
					Condition:     "field3",
					TypeConverter: SqlConvertBool,
				},
				"field4": SqlFieldTranslation{
					Condition:     "field4",
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
					Condition: "field1",
				},
				"field2": SqlFieldTranslation{
					Condition: "field2",
				},
				"field3": SqlFieldTranslation{
					Condition: "field3",
				},
				"field4": SqlFieldTranslation{
					Condition: "field4",
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
					Condition:     "field1",
					TypeConverter: SqlConvertDateTime,
				},
				"field2": SqlFieldTranslation{
					Condition:     "field2",
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
