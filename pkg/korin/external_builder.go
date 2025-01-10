package korin

import (
	"bytes"
	"fmt"
	"github.com/ShindouMihou/go-little-utils/slices"
	"github.com/ShindouMihou/korin/internal/kcomp"
	"github.com/ShindouMihou/korin/pkg/kplugins"
	"github.com/ttacon/chalk"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	NoOpLogger Logger = func(args ...any) {}
)

type Logger func(args ...any)
type Korin struct {
	BuildDirectory string
	Plugins        []kplugins.Plugin
	BuildCommand   string
	Logger         Logger
}

// Process processes the directory specified in the `dir` parameter.
// Unlike `Build`, this function returns the `errors` that occurred during the processing, and won't
// print them out.
func (ko Korin) Process(dir string) []error {
	return kcomp.Process(&kcomp.Configuration{
		BuildDirectory: ko.BuildDirectory,
		Plugins:        ko.Plugins,
		BuildCommand:   ko.BuildCommand,
		Logger:         ko.Logger,
	}, dir)
}

func (ko Korin) handle(errs []error) bool {
	if len(errs) > 0 {
		fmt.Println()
		fmt.Println(
			chalk.Bold.
				NewStyle().
				WithForeground(chalk.Red).
				WithBackground(chalk.ResetColor).
				Style("----------------------"),
			"processing of files failed with the following errors:",
		)

		for i, err := range errs {
			if i > 0 {
				fmt.Println()
			}
			fmt.Println(
				chalk.Bold.
					NewStyle().
					WithForeground(chalk.Red).
					WithBackground(chalk.ResetColor).
					Style("err:"),
				err.Error(),
			)
		}
		return false
	}
	return true
}

// Build builds the directory specified in the `dir` parameter.
// Unlike process, this will actually handle the errors and prints them out.
func (ko Korin) Build(dir string) {
	errs := ko.Process(dir)
	if !ko.handle(errs) {
		return
	}
}

// DockerBuildStep is a synonymous function to `Build` but is specifically for Docker builds,
// as this will exit the program with `os.Exit(1)` which will signify to Docker that the
// build step failed.
func (ko Korin) DockerBuildStep(dir string) {
	errs := ko.Process(dir)
	if !ko.handle(errs) {
		os.Exit(1)
	}
}

// BuildRootDirectory builds the root directory (which is "." directory of the project)
func (ko Korin) BuildRootDirectory() {
	ko.Build(".")
}

// BuildWorkingDirectory builds the current working directory of the project, based on `os.Getwd()`
func (ko Korin) BuildWorkingDirectory() {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	ko.Build(currentWorkingDir)
}

// Run builds the root directory, and then runs the specified path, it will be automatically
// be proxied to the build directory, so the path, such as "cmd/main.go" will be automatically proxied to
// ".build/cmd/main.go" instead.
func (ko Korin) Run(path string) {
	errs := ko.Process(".")
	if !ko.handle(errs) {
		return
	}

	ko.Logger()
	ko.Logger(
		chalk.Bold.
			NewStyle().
			WithForeground(chalk.Green).
			WithBackground(chalk.ResetColor).
			Style("running:"),
		path,
	)
	ko.Logger(
		chalk.Bold.
			NewStyle().
			WithForeground(chalk.Yellow).
			WithBackground(chalk.ResetColor).
			Style("-----------------------------"),
		path,
	)
	ko.Logger()

	var stdoutBuf, stderrBuf, stdinBuf bytes.Buffer
	cmd := exec.Command("go", "run", filepath.Join(".build/", path))
	cmd.Stdin = io.MultiReader(os.Stdin, &stdinBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}
}

// Plugin adds a plugin to the list of plugins that will be used during the processing of the files.
func (ko Korin) Plugin(plugin kplugins.Plugin) {
	if slices.Filter(ko.Plugins, func(b kplugins.Plugin) bool {
		return b.Name() == plugin.Name() && b.Group() == b.Group()
	}) != nil {
		return
	}
	ko.Plugins = append(ko.Plugins, plugin)
}

// New creates a new instance of Korin with the default values.
func New() *Korin {
	return &Korin{
		BuildDirectory: ".build/",
		Plugins: []kplugins.Plugin{
			kplugins.NewErrorPropogationPlugin(),
			kplugins.NewPrintLinePlugin(),
			kplugins.NewPluginSerializerAnnotations(),
			kplugins.NewPluginEnvironmentKey(),
		},
		BuildCommand: "go build {$BUILD_FOLDER}/{$FILE_NAME}.go",
		Logger: func(args ...any) {
			fmt.Println(args...)
		},
	}
}
