package kcomp

import (
	"errors"
	"fmt"
	"github.com/ShindouMihou/go-little-utils/slices"
	"github.com/ShindouMihou/korin/internal/kproc"
	"github.com/ShindouMihou/korin/internal/kproc/labelers"
	"github.com/ShindouMihou/korin/internal/kstrings"
	"github.com/ShindouMihou/korin/pkg/klabels"
	"github.com/ShindouMihou/korin/pkg/kplugins"
	"github.com/ShindouMihou/siopao/siopao"
	"github.com/inancgumus/screen"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/ttacon/chalk"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Configuration struct {
	BuildDirectory string
	Plugins        []kplugins.Plugin
	BuildCommand   string
	Logger         func(args ...any)
}

func logSettings(config *Configuration) {
	screen.Clear()
	config.Logger(
		chalk.Bold.
			NewStyle().
			WithForeground(chalk.Magenta).
			WithBackground(chalk.ResetColor).
			Style("korin where it matters"),
	)
	config.Logger(
		chalk.Bold.
			NewStyle().
			WithTextStyle(chalk.Italic).
			WithForeground(chalk.ResetColor).
			WithBackground(chalk.ResetColor).
			Style("the simple unnecessary golang preprocessor"),
	)
	config.Logger()
	config.Logger(
		chalk.Bold.
			NewStyle().
			WithForeground(chalk.Yellow).
			WithBackground(chalk.ResetColor).
			Style("build dir:"),
		config.BuildDirectory,
	)
	config.Logger(
		chalk.Bold.
			NewStyle().
			WithForeground(chalk.Yellow).
			WithBackground(chalk.ResetColor).
			Style("build command:"),
		config.BuildCommand,
	)
	config.Logger(
		chalk.Bold.
			NewStyle().
			WithForeground(chalk.Yellow).
			WithBackground(chalk.ResetColor).
			Style("plugins:"),
		kplugins.SyntaxHelper.CommaSeparate(slices.Map(config.Plugins, func(v kplugins.Plugin) string {
			return fmt.Sprintf("%s.%s:%s", v.Group(), v.Name(), v.Version())
		})),
	)
	config.Logger()
}

func readIgnoreFile(file string) []string {
	gitignore := siopao.Open(file)
	reader, err := gitignore.TextReader()
	if err != nil {
		return []string{}
	}

	var patterns []string
	if err = reader.EachLine(func(line string) {
		if kstrings.HasPrefix(line, "#") {
			return
		}
		patterns = append(patterns, line)
	}); err != nil {
		return patterns
	}
	return patterns
}

func getIgnorePatterns(dir string) *ignore.GitIgnore {
	gitignore := filepath.Join(dir, ".gitignore")
	korinignore := filepath.Join(dir, ".korignore")

	reader, err := siopao.Open(korinignore).TextReader()
	if err == nil {
		var patterns []string
		_ = reader.EachLine(func(line string) {
			if line[0] == '#' {
				return
			}

			patterns = append(patterns, line)
		})
		return ignore.CompileIgnoreLines(patterns...)
	}

	reader, err = siopao.Open(gitignore).TextReader()
	if err == nil {
		var patterns []string
		_ = reader.EachLine(func(line string) {
			if line[0] == '#' {
				return
			}

			patterns = append(patterns, line)
		})
		return ignore.CompileIgnoreLines(patterns...)
	}

	return nil
}

func Process(config *Configuration, dir string) []error {
	start := time.Now()

	logSettings(config)
	if err := createBuildDirectory(config); err != nil {
		return []error{err}
	}

	ignorePatterns := getIgnorePatterns(dir)
	source := siopao.Open(dir)

	var errs []error
	var wg sync.WaitGroup

	if err := source.Recurse(true, func(file *siopao.File) {
		if strings.Contains(file.Path(), ".build/") {
			return
		}

		ext := filepath.Ext(file.Path())
		if ext != ".go" {
			return
		}

		if ignorePatterns != nil && ignorePatterns.MatchesPath(file.Path()) {
			config.Logger(
				chalk.Bold.
					NewStyle().
					WithForeground(chalk.Red).
					WithBackground(chalk.ResetColor).
					Style("ignored:"),
				file.Path(),
			)
			return
		}

		wg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					errs = append(
						errs,
						errors.Join(
							fmt.Errorf("failed to process file %s", file.Path()),
							r.(error),
						),
					)
				}
				wg.Done()
			}()
			contents, err := process(config, file)
			if err != nil {
				errs = append(
					errs,
					errors.Join(
						fmt.Errorf("failed to process file %s", file.Path()),
						err,
					),
				)
				return
			}
			path := filepath.Join(".build/", file.Path())
			if err := siopao.Open(path).Overwrite(contents); err != nil {
				errs = append(
					errs,
					errors.Join(
						fmt.Errorf("failed to process file %s", file.Path()),
						err,
					),
				)
				return
			}
			config.Logger(
				chalk.Bold.
					NewStyle().
					WithForeground(chalk.Green).
					WithBackground(chalk.ResetColor).
					Style("processed:"),
				file.Path(),
			)
		}()
	}); err != nil {
		return []error{err}
	}

	wg.Wait()
	config.Logger()
	config.Logger(
		chalk.Bold.
			NewStyle().
			WithForeground(chalk.Yellow).
			WithBackground(chalk.ResetColor).
			Style("processing time:"),
		time.Since(start).Milliseconds(),
		"ms",
	)

	if len(errs) > 0 {
		return errs
	}
	exec.Command("gofmt", ".build/...")
	return nil
}

func createBuildDirectory(config *Configuration) error {
	build := siopao.Open(config.BuildDirectory)
	err := build.DeleteRecursively()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(config.BuildDirectory), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func process(config *Configuration, file *siopao.File) (string, error) {
	reader, err := file.TextReader()
	if err != nil {
		panic(err)
	}
	index := 0
	var stack []klabels.Analysis

	var headers kplugins.Headers
	var contents kplugins.Writer

	isInImportScope := false

	isInTypeDeclaration := false
	isInConstScope, isInVarScope := false, false

	if err := reader.EachLine(func(line string) {
		labels := kproc.LabelLine(index, line)
		if isInTypeDeclaration {
			if line == "}" {
				isInTypeDeclaration = false
				labels.Labels = append(labels.Labels, klabels.Label{Kind: klabels.ScopeEndKind})
			} else {
				labels.Labels = append(labels.Labels, labelers.FieldDeclaration(line))
			}
		}

		if isInConstScope || isInVarScope {
			if line == ")" {
				kind := klabels.VarScopeEndKind
				if isInConstScope {
					kind = klabels.ConstScopeEndKind
				}
				labels.Labels = append(labels.Labels, klabels.Label{Kind: klabels.LabelKind(kind)})
				isInConstScope, isInVarScope = false, false
			} else {
				kind := klabels.VarDeclarationKind
				if isInConstScope {
					kind = klabels.ConstDeclarationKind
				}
				labels.Labels = append(labels.Labels, klabels.Label{Kind: klabels.LabelKind(kind)})
			}
		}

		if len(labels.Labels) > 0 {
			if kplugins.ReadHelper.Get(klabels.ConstScopeBeginKind, labels.Labels) != nil {
				isInConstScope = true
			} else if kplugins.ReadHelper.Get(klabels.VarScopeBeginKind, labels.Labels) != nil {
				isInVarScope = true
			} else if kplugins.ReadHelper.Get(klabels.TypeDeclarationKind, labels.Labels) != nil && kstrings.HasSuffix(line, "{") {
				isInTypeDeclaration = true
			}
		}

		stack = append(stack, labels)
		defer func() {
			index++
		}()

		if kstrings.HasPrefix(line, "import (") {
			isInImportScope = true
			return
		}
		if kstrings.HasPrefix(line, "import \"") {
			pkg := strings.TrimPrefix(line, "import \"")
			pkg = strings.TrimSuffix(pkg, "\"")
			headers.Import(pkg)
			return
		}
		if isInImportScope {
			if line == ")" {
				isInImportScope = false
			} else {
				pkg := strings.TrimSpace(line)
				pkg = line[2:len(pkg)]
				headers.Import(pkg)
			}
			return
		}
		if kstrings.HasPrefix(line, "package ") {
			headers.Package(line[len("package "):])
			return
		}
		hasChanges := false
		for _, plugin := range config.Plugins {
			result, err := plugin.Process(line, index, &headers, stack)
			if err != nil {
				panic(err)
			}
			if result != "" {
				hasChanges = true
				contents.Write(result)
				break
			}
		}
		if !hasChanges {
			contents.Write(line)
		}
		contents.NextLine()
	}); err != nil {
		return "", err
	}
	return headers.Format() + contents.Contents(), nil
}
