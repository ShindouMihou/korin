# Korin

An experimental, created for learning and fun, Golang preprocessor. Korin is a preprocessor that is embedded with the 
codebase enabling Golang codebases to minimize code duplication on some parts with the use of preprocessing before 
compilation.

Korin works by reading each Golang line-by-line, processing and gaining an understanding of each line to allow "plugins" 
to read the lines and decide whether to preprocess the line to add changes, such as adding new lines after, or 
modifying the line entirely.

As an example, one of the default plugins is `k:float` which propogates errors back up to the function, for example:
```go
stats, _ := json.Marshal(statistics)            // +k:float
```

The above line will be transformed into:
```go
stats, err := json.Marshal(statistics)
if err != nil {
	return 
}
```

It returns based on what the function's result type's zero-value is, for example, a function with a result type of 
`int` will return `0` while a function with a result type of `string` will return `""`.

You can find the source code of the error propagation plugin under [`plugin_error_propogation.go`](pkg/korin/plugin_error_propgation.go).

Another interesting plugin is the `k:named` which will produce annotations for field declaration using the name of the field, for example, 
the following:
```go
type Test struct {
	NameCharacters string // +k:named(json,yaml,bson)
}
```

The above will be transformed into:
```go
type Test struct {
	NameCharacters string `json:"name_characters" yaml:"name_characters" bson:"name_characters"`
}
```

You can also specify whether to use camel case or snake case by adding the `snake_case` or `camelCase` argument to the plugin, **it should be the 
first argument**. For example:
```go
type Test struct {
    NameCharacters string // +k:named(camelCase,json,yaml,bson)
}
````

Additionally, you can also specify it to keep the original name by adding the `original` argument to the plugin, **it should be the first argument**.
For example:
```go
type Test struct {
    NameCharacters string // +k:named(original,json,yaml,bson)
}
```

## Usage

Korin is designed in a way that it takes over the entrypoint to proxy to the actual entrypoint which will run the application, such that, 
you will run `go run cmd/korin.go` instead of your usual `go run cmd/app.go`. To begin, create a file under your `cmd` folder called `korin.go` and 
copy the lines:
```go
package main

import "korin/pkg/kbuild"

func main() {
	korin := kbuild.NewKorin()
	korin.Run("cmd/app.go")  // Replace this value with your entry point if it isn't cmd/app.go
}
```

Korin will scan through the entire codebase, respecting the `.korignore` file, to analyze and check for any files that needs preprocessing, 
and will preprocess any found. All processed files will be under the configured build directory (`.build`, by default). Generally, a codebase of roughly 
2,000-4,000 lines should take under 5-10 milliseconds to preprocess, but this also depends on some factors, but generally, it should be fast.

## License

Korin is licensed under the MIT license. you are free to redistribute, use and do other related actions under the MIT license. 
this is permanent and irrevocable.