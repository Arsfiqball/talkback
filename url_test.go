package talkback

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromQueryString(t *testing.T) {
	type scenarioT struct {
		query string
		out   Query
	}

	scenarios := []scenarioT{
		{
			query: "field1_eq=value1&field2_ne=value2&field3_isnull=true",
			out: Query{
				Conditions: []Condition{
					{"field1", "eq", []string{"value1"}},
					{"field2", "ne", []string{"value2"}},
					{"field3", "isnull", []string{"true"}},
				},
			},
		},
		{
			query: "field1_gt=value1&field2_lt=value2&field3_gte=value3&field4_lte=value4",
			out: Query{
				Conditions: []Condition{
					{"field1", "gt", []string{"value1"}},
					{"field2", "lt", []string{"value2"}},
					{"field3", "gte", []string{"value3"}},
					{"field4", "lte", []string{"value4"}},
				},
			},
		},
		{
			query: "field1_contain=value1&field2_ncontain=value2&field3_contains=value3&field4_ncontains=value4",
			out: Query{
				Conditions: []Condition{
					{"field1", "contain", []string{"value1"}},
					{"field2", "ncontain", []string{"value2"}},
					{"field3", "contains", []string{"value3"}},
					{"field4", "ncontains", []string{"value4"}},
				},
			},
		},
		{
			query: "field1_in=value1&field2_nin=value2",
			out: Query{
				Conditions: []Condition{
					{"field1", "in", []string{"value1"}},
					{"field2", "nin", []string{"value2"}},
				},
			},
		},
		{
			query: "sort=-field1&sort=field2",
			out: Query{
				Sort: []Sort{
					{"field1", true},
					{"field2", false},
				},
			},
		},
		{
			query: "with=field1&with=field2",
			out: Query{
				With: []string{"field1", "field2"},
			},
		},
		{
			query: "group=field1&group=field2",
			out: Query{
				Group: []string{"field1", "field2"},
			},
		},
		{
			query: "accumulator=field1&accumulator=field2",
			out: Query{
				Accumulator: []string{"field1", "field2"},
			},
		},
		{
			query: "limit=200&skip=20",
			out: Query{
				Limit: 200,
				Skip:  20,
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.query, func(t *testing.T) {
			out, err := FromQueryString(scenario.query)

			assert.NoError(t, err, "error should be nil")
			assert.ElementsMatch(t, scenario.out.Conditions, out.Conditions, "conditions should match")
		})
	}
}
