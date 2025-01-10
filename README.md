# Korin

Korin is an experimental Golang preprocessor designed for learning and fun. It empowers developers to minimize code duplication by enabling preprocessing before compilation. By embedding directly into your codebase, Korin enhances code maintainability with its extensible plugin architecture.

## How Korin Works

Korin processes your Golang code line-by-line, allowing plugins to analyze, modify, or extend the code dynamically. Each plugin can:
- Add new lines
- Modify existing lines
- Transform specific code patterns

### Example: `k:float` Plugin
The `k:float` plugin automates error propagation:
```go
stats, _ := json.Marshal(statistics)            // +k:float
```
Becomes:
```go
stats, err := json.Marshal(statistics)
if err != nil {
    return
}
```
This transformation ensures proper error handling, returning the zero value based on the function's result type (e.g., `0` for `int`, `""` for `string`). The source code for this plugin can be found in [`plugin_error_propagation.go`](pkg/korin/plugin_error_propagation.go).

### Example: `k:named` Plugin
The `k:named` plugin generates field annotations based on field names. For instance:
```go
type Test struct {
    NameCharacters string // +k:named(json,yaml,bson)
}
```
```go
// +k:named(json,yaml,bson)
type TestWithMultiline struct {
    NameCharacters string 
    Age int
}
```
Transforms to:
```go
type Test struct {
    NameCharacters string `json:"name_characters" yaml:"name_characters" bson:"name_characters"`
}
```
```go
type TestWithMultiline struct {
    NameCharacters string `json:"name_characters" yaml:"name_characters" bson:"name_characters"`
    Age int `json:"age" yaml:"age" bson:"age"`
}
```
You can customize the case format:
- **Camel Case:**
  ```go
  type Test struct {
      NameCharacters string // +k:named(camelCase,json,yaml,bson)
  }
  ```
- **Snake Case:** Default behavior.
- **Original Name:**
  ```go
  type Test struct {
      NameCharacters string // +k:named(original,json,yaml,bson)
  }
  ```
  
### Example: `k:env` Plugin

The `k:env` plugin automatically sets the value of a variable to its corresponding environment value. For instance:
```go
var Port = "{$ENV:PORT}" // +k:env
```

Transforms to:
```go
var Port = "8080"
````

You can also specify a type, and Korin will try to automatically convert it to that value on the fly. For instance:
```go
const Port = "{$ENV:PORT}" // +k:env(int)
```

When using the types `bool`, any `int`, any `float`, or `rune`, Korin will automatically drop the value of the environment variable 
as is, for instance, with the above code, it becomes:
```go
const Port int = 8080
```

Otherwise, if it's a `string` or an unknown type, such as a typealias, or a custom type, Korin will automatically convert the value to a string, for instance:
```go
const Port = "{$ENV:PORT}" // +k:env(PortType)
```

Transforms to:
```go
const Port PortType = "8080"
```

## Supported Syntaxes
Korin supports preprocessing for:
- **Type declarations**
- **Field declarations**
- **Variable declarations**
- **Function declarations** (including anonymous functions)
- **Import declarations**
- **Package declarations**
- **Return statements**
- **Comment declarations**
- **Const declarations** (including multiple constants)
- **Var declarations** (including multiple variables)

## Usage

Korin intercepts the entry point to preprocess your codebase. To use Korin, create a `cmd/korin.go` file and add the following:
```go
package main

import "github.com/ShindouMihou/korin/pkg/korin"

func main() {
    korin := korin.New()

    // Uncomment to silence the logger (error logs during preprocessing remain unaffected)
    // korin.Logger = kbuild.NoOpLogger
	
    korin.Run("cmd/app.go")  // Replace with your actual entry point
}
```
Korin scans the codebase (excluding files listed in `.korignore`), preprocesses required files, and outputs them to the configured build directory (default: `.build`).

### Using Korin as a Docker Buildstep
We recommend using Korin as a build step in Dockerfile to ensure that the code is preprocessed before building. To use Korin as a Docker build step, we recommend creating a 
file  called `cmd/korin_build.go` and adding the following:
```go
package main

import "github.com/ShindouMihou/korin/pkg/korin"

func main() {
    korin := korin.New()

    // Uncomment to silence the logger (error logs during preprocessing remain unaffected)
    // korin.Logger = kbuild.NoOpLogger
	
    korin.DockerBuildStep(".")  // Replace with your path to source files if it isn't "."
}
```

Once done, you can add `RUN go run cmd/korin_build.go` as a Docker step and it should preprocess the files beforehand. It is recommended to use this method over `Build` and 
related since it will exit the program with an exit code when it fails to preprocess a specific file.

### Plugins
To create a plugin, you can read about it in the [Creating a Plugin](docs/creating_a_plugin.md) documentation, otherwise, 
you can check out the [native plugins](pkg/kplugins) that are already available, these should have the `plugin_` prefix on their filenames, by 
default, all these plugins are automatically included when you create a new `korin` instance.

### Performance
Korin processes a 2,000-4,000 line codebase in under 5-10 milliseconds on average, though performance may vary depending on project complexity.

## License

Korin is licensed under the [MIT License](LICENSE). You are free to use, modify, and distribute this project in compliance with the license terms.

---
Elevate your Golang projects with Korin's seamless preprocessing capabilities!
