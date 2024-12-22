package kbuild

import (
	"bytes"
	"fmt"
	"github.com/ShindouMihou/go-little-utils/slices"
	"github.com/ttacon/chalk"
	"io"
	"korin/internal/kcomp"
	"korin/pkg/korin"
	"os"
	"os/exec"
	"path/filepath"
)

type Korin struct {
	BuildDirectory string
	Plugins        []korin.Plugin
	BuildCommand   string
}

func (ko Korin) Process(dir string) []error {
	return kcomp.Process(&kcomp.Configuration{
		BuildDirectory: ko.BuildDirectory,
		Plugins:        ko.Plugins,
		BuildCommand:   ko.BuildCommand,
	}, dir)
}

func handle(errs []error) bool {
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

func (ko Korin) Build(dir string) {
	errs := ko.Process(".")
	if !handle(errs) {
		return
	}
}

func (ko Korin) Run(path string) {
	errs := ko.Process(".")
	if !handle(errs) {
		return
	}

	fmt.Println()
	fmt.Println(
		chalk.Bold.
			NewStyle().
			WithForeground(chalk.Green).
			WithBackground(chalk.ResetColor).
			Style("running:"),
		path,
	)
	fmt.Println(
		chalk.Bold.
			NewStyle().
			WithForeground(chalk.Yellow).
			WithBackground(chalk.ResetColor).
			Style("-----------------------------"),
		path,
	)
	fmt.Println()

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

func (ko Korin) Plugin(plugin korin.Plugin) {
	if slices.Filter(ko.Plugins, func(b korin.Plugin) bool {
		return b.Name() == plugin.Name() && b.Group() == b.Group()
	}) != nil {
		return
	}
	ko.Plugins = append(ko.Plugins, plugin)
}

func NewKorin() *Korin {
	return &Korin{
		BuildDirectory: ".build/",
		Plugins: []korin.Plugin{
			korin.ErrorPropogationPlugin{},
			korin.PrintLinePlugin{},
		},
		BuildCommand: "go build {$BUILD_FOLDER}/{$FILE_NAME}.go",
	}
}
