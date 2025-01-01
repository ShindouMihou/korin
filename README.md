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
Transforms to:
```go
type Test struct {
    NameCharacters string `json:"name_characters" yaml:"name_characters" bson:"name_characters"`
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

## Upcoming Features
- **Const declarations**: Planned for future releases.
- **Var declarations**: Coming soon.

## Usage

Korin intercepts the entry point to preprocess your codebase. To use Korin, create a `cmd/korin.go` file and add the following:
```go
package main

import "github.com/ShindouMihou/korin/pkg/kbuild"

func main() {
    korin := kbuild.NewKorin()

    // Uncomment to silence the logger (error logs during preprocessing remain unaffected)
    // korin.Logger = kbuild.NoOpLogger

    korin.Run("cmd/app.go")  // Replace with your actual entry point
}
```
Korin scans the codebase (excluding files listed in `.korignore`), preprocesses required files, and outputs them to the configured build directory (default: `.build`).

### Performance
Korin processes a 2,000-4,000 line codebase in under 5-10 milliseconds on average, though performance may vary depending on project complexity.

## License

Korin is licensed under the [MIT License](LICENSE). You are free to use, modify, and distribute this project in compliance with the license terms.

---
Elevate your Golang projects with Korin's seamless preprocessing capabilities!