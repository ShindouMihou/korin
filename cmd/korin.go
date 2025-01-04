package main

import "github.com/ShindouMihou/korin/pkg/korin"

func main() {
	korin := korin.New()

	// You should use this when building a Docker image.
	// korin.DockerBuildStep("examples/")

	korin.Run("examples/sample.go")
}
