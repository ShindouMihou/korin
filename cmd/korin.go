package main

import "github.com/ShindouMihou/korin/pkg/kbuild"

func main() {
	korin := kbuild.NewKorin()
	korin.Run("examples/sample.go")
}
