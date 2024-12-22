package main

import "korin/pkg/kbuild"

func main() {
	korin := kbuild.NewKorin()
	korin.Run("examples/sample.go")
}
