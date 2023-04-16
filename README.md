# Talkback Lancer

Simplify the way to implement query object pattern to your Go app.


Install:
```sh
go get -u github.com/Arsfiqball/talkback-lancer
```

Use:
```go
func main() {
	urlQS := "field1_eq=value1&field2_ne=value2&field3_isnull=true&group=field1&group=field2&sort=field1&sort=-field2&limit=10&skip=10"

	query, err := FromQueryString(urlQS)
	if err != nil {
		panic(err)
	}

	translations := SqlTranslations{
		"field1": SqlFieldTranslation{
			Alias: "alias1",
		},
		"field2": SqlFieldTranslation{
			Column: "ex.field2",
		},
		"field3": SqlFieldTranslation{},
	}

	statement, args, err := ToSql("somewhere", query, translations)
	if err != nil {
		panic(err)
	}

	fmt.Println(statement)
	fmt.Println(args)
}
```

See more in [API Docs](/api.md)
