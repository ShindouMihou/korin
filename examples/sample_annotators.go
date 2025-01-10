package main

type Test struct {
	NameCharacters string // +k:named(camelCase,json,yaml,bson)
}

// +k:named(camelCase,json,yaml,bson)
type TestWithStruct struct {
	NameCharacters string
	Test           string
}

type TestWithStruct2 struct { // +k:named(camelCase,json,yaml,bson)
	NameCharacters string
	Test           string
}

// +k:named(camelCase,json,yaml,bson)
type TestWithStructOverride struct {
	NameCharacters string // +k:named(snake_case,json,yaml,bson)
	Test           string
}
