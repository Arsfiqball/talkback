package talkback

import (
	"net/url"
	"strconv"
	"strings"
)

// FromQueryString returns a Query from a query string.
func FromQueryString(qs string) (Query, error) {
	params, err := url.ParseQuery(qs)
	if err != nil {
		return Query{}, err
	}

	return FromURLValues(params)
}

// FromQueryString returns a Query from a query string.
func FromURLValues(params url.Values) (Query, error) {
	query := Query{}

	for key, values := range params {
		spliten := strings.Split(key, "_")

		if len(spliten) < 2 {
			continue
		}

		field := strings.Join(spliten[:len(spliten)-1], "_")
		op := spliten[len(spliten)-1]

		cond := Condition{
			Field:  field,
			Op:     op,
			Values: values,
		}

		if cond.Valid() {
			query.Conditions = append(query.Conditions, cond)
		}
	}

	if with, ok := params["with"]; ok {
		query.With = with
	}

	if group, ok := params["group"]; ok {
		query.Group = group
	}

	if accumulator, ok := params["accumulator"]; ok {
		query.Accumulator = accumulator
	}

	if sort, ok := params["sort"]; ok {
		for _, s := range sort {
			field := s
			reverse := false

			if strings.HasPrefix(field, "-") {
				field = field[1:]
				reverse = true
			}

			query.Sort = append(query.Sort, Sort{
				Field:   field,
				Reverse: reverse,
			})
		}
	}

	if limit := params.Get("limit"); limit != "" {
		var err error

		query.Limit, err = strconv.Atoi(limit)
		if err != nil {
			return Query{}, err
		}
	}

	if skip := params.Get("skip"); skip != "" {
		var err error

		query.Skip, err = strconv.Atoi(skip)
		if err != nil {
			return Query{}, err
		}
	}

	return query, nil
}
