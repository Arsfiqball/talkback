package talkback

import (
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// MongoFieldTranslation is a translation of a field to a MongoDB field.
type MongoFieldTranslation struct {
	Field         string
	Accumulation  bson.M
	TypeConverter func(value string) (interface{}, error)
}

// MongoTranslations is a map of field names to their translations.
type MongoTranslations map[string]MongoFieldTranslation

// ToMongoFilter converts a Query to a MongoDB filter statement.
func ToMongoFilter(query Query, translations MongoTranslations) (bson.D, error) {
	d := bson.D{}

	translations = sanitizeMongoTranslation(translations)

	for _, cond := range query.Conditions {
		translation, ok := translations[cond.Field]
		if !ok {
			return nil, ErrInvalidField
		}

		convertedCond, err := condtionToMongo(translation, cond)
		if err != nil {
			return nil, err
		}

		d = append(d, convertedCond)
	}

	return d, nil
}

func sanitizeMongoTranslation(translations MongoTranslations) MongoTranslations {
	result := MongoTranslations{}

	for field, translation := range translations {
		if translation.Field == "" {
			translation.Field = field
		}

		result[field] = translation
	}

	return result
}

func condtionToMongo(translation MongoFieldTranslation, cond Condition) (bson.E, error) {
	values := []interface{}{}

	for _, value := range cond.Values {
		convertedValue, err := translation.TypeConverter(value)
		if err != nil {
			return bson.E{}, fmt.Errorf("failed to convert value of %s: %w", cond.Field, err)
		}

		values = append(values, convertedValue)
	}

	firstValue := values[0]
	field := translation.Field

	switch cond.Op {
	case OpIsNull:
		return bson.E{Key: field, Value: bson.M{"$exists": false}}, nil
	case OpEq:
		return bson.E{Key: field, Value: firstValue}, nil
	case OpNe:
		return bson.E{Key: field, Value: bson.M{"$ne": firstValue}}, nil
	case OpGt:
		return bson.E{Key: field, Value: bson.M{"$gt": firstValue}}, nil
	case OpGte:
		return bson.E{Key: field, Value: bson.M{"$gte": firstValue}}, nil
	case OpLt:
		return bson.E{Key: field, Value: bson.M{"$lt": firstValue}}, nil
	case OpLte:
		return bson.E{Key: field, Value: bson.M{"$lte": firstValue}}, nil
	case OpContain:
		return bson.E{Key: field, Value: bson.M{"$regex": firstValue, "$options": "i"}}, nil
	case OpNcontain:
		return bson.E{Key: field, Value: bson.M{"$not": bson.M{"$regex": firstValue, "$options": "i"}}}, nil
	case OpContains:
		return bson.E{Key: field, Value: bson.M{"$regex": firstValue}}, nil
	case OpNcontains:
		return bson.E{Key: field, Value: bson.M{"$not": bson.M{"$regex": firstValue}}}, nil
	case OpIn:
		return bson.E{Key: field, Value: bson.M{"$in": values}}, nil
	case OpNin:
		return bson.E{Key: field, Value: bson.M{"$nin": values}}, nil
	default:
		return bson.E{}, ErrInvalidOp
	}
}

// MongoConvertString is a TypeConverter that converts a string to a string.
func MongoConvertString(value string) (interface{}, error) {
	return value, nil
}

// MongoConvertString is a TypeConverter that converts a string to a string.
func MongoConvertInt(value string) (interface{}, error) {
	return strconv.Atoi(value)
}

// MongoConvertString is a TypeConverter that converts a string to a string.
func MongoConvertBool(value string) (interface{}, error) {
	return strconv.ParseBool(value)
}

// MongoConvertString is a TypeConverter that converts a string to a string.
func MongoConvertFloat(value string) (interface{}, error) {
	return strconv.ParseFloat(value, 64)
}

// MongoConvertString is a TypeConverter that converts a string to a string.
func MongoConvertDate(value string) (interface{}, error) {
	return time.Parse("2006-01-02", value)
}

// MongoConvertString is a TypeConverter that converts a string to a string.
func MongoConvertDateTime(value string) (interface{}, error) {
	return time.Parse("2006-01-02 15:04:05", value)
}

// MongoConvertString is a TypeConverter that converts a string to a string.
func MongoConvertTime(value string) (interface{}, error) {
	return time.Parse("15:04:05", value)
}

// MongoConvertString is a TypeConverter that converts a string to a string.
func MongoConvertISO8601(value string) (interface{}, error) {
	return time.Parse(time.RFC3339, value)
}

// ToMongoProject converts a Query to a MongoDB project statement.
func ToMongoProject(query Query, translations MongoTranslations) (bson.D, error) {
	return bson.D{}, nil // TODO: Implement
}

// ToMongoGroup converts a Query to a MongoDB group statement.
func ToMongoGroup(query Query, translations MongoTranslations) (bson.D, error) {
	id := bson.M{}
	accumulators := bson.D{}
	translations = sanitizeMongoTranslation(translations)

	for _, group := range query.Group {
		translation, ok := translations[group]
		if !ok {
			return nil, ErrInvalidField
		}

		id[group] = "$" + translation.Field
	}

	for _, field := range query.Accumulator {
		translation, ok := translations[field]
		if !ok {
			return nil, ErrInvalidField
		}

		accumulators = append(accumulators, bson.E{Key: field, Value: translation.Accumulation})
	}

	return append(bson.D{{Key: "_id", Value: id}}, accumulators...), nil
}

// ToMongoSort converts a Query to a MongoDB sort statement.
func ToMongoSort(query Query, translations MongoTranslations) (bson.D, error) {
	d := bson.D{}
	translations = sanitizeMongoTranslation(translations)

	for _, field := range query.Sort {
		translation, ok := translations[field.Field]
		if !ok {
			return nil, ErrInvalidField
		}

		direction := 1

		if field.Reverse {
			direction = -1
		}

		d = append(d, bson.E{Key: translation.Field, Value: direction})
	}

	return d, nil
}

// ToMongoLimit converts a Query to a MongoDB limit statement.
func ToMongoLimit(query Query) (int, error) {
	return query.Limit, nil
}

// ToMongoSkip converts a Query to a MongoDB skip statement.
func ToMongoSkip(query Query) (int, error) {
	return query.Skip, nil
}

// MongoPreloadable is a map of preloads (key) and their corresponding model (value).
type MongoPreloadable map[string]string

// ToMongoPreload converts a Query to a SQL preload statement.
func ToMongoPreload(query Query, preloadable MongoPreloadable) ([]string, error) {
	return query.With, nil // TODO: Implement
}

// MongoPlan is a plan for executing a query.
type MongoPlan struct {
	Project bson.D
	Filter  bson.D
	Group   bson.D
	Sort    bson.D
	Limit   int
	Offset  int
	Preload []string
}

// ToMongoPlan converts a Query to a MongoPlan.
func ToMongoPlan(query Query, translations MongoTranslations, preloadable MongoPreloadable) (MongoPlan, error) {
	cproject, err := ToMongoProject(query, translations)
	if err != nil {
		return MongoPlan{}, err
	}

	cfilter, err := ToMongoFilter(query, translations)
	if err != nil {
		return MongoPlan{}, err
	}

	cgroup, err := ToMongoGroup(query, translations)
	if err != nil {
		return MongoPlan{}, err
	}

	csort, err := ToMongoSort(query, translations)
	if err != nil {
		return MongoPlan{}, err
	}

	climit, err := ToMongoLimit(query)
	if err != nil {
		return MongoPlan{}, err
	}

	cskip, err := ToMongoSkip(query)
	if err != nil {
		return MongoPlan{}, err
	}

	cpreload, err := ToMongoPreload(query, preloadable)
	if err != nil {
		return MongoPlan{}, err
	}

	return MongoPlan{
		Project: cproject,
		Filter:  cfilter,
		Group:   cgroup,
		Sort:    csort,
		Limit:   climit,
		Offset:  cskip,
		Preload: cpreload,
	}, nil
}
