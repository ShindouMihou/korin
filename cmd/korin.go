package main

import "github.com/ShindouMihou/korin/pkg/korin"

func main() {
	korin := korin.New()
	korin.Run("examples/sample.go")
}
