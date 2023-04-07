package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/mazrean/genorm/cmd/genorm/generator"
)

var (
	// flags
	showVersionInfo bool
	source          string
	destination     string
	packageName     string
	moduleName      string
	joinNum         int
)

func init() {
	flag.BoolVar(&showVersionInfo, "version", false, "If true, output version information.")
	flag.StringVar(&source, "source", "", "The source file to parse.")
	flag.StringVar(&destination, "destination", "", "The destination file to write.")
	flag.StringVar(&packageName, "package", "", "The root package name to use.")
	flag.StringVar(&moduleName, "module", "", "The root module name to use.")
	flag.IntVar(&joinNum, "join-num", 5, "The number of joins to generate.")
}

func main() {
	flag.Parse()

	if showVersionInfo {
		err := printVersionInfo()
		if err != nil {
			panic(err)
		}
	}

	src, err := openSource(source)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	if len(destination) == 0 {
		panic("Destination directory path is required.")
	}

	if len(packageName) == 0 {
		panic("package name is required")
	}
	if len(moduleName) == 0 {
		panic("module name is required")
	}

	err = generator.Generate(packageName, moduleName, destination, src, generator.Config{
		JoinNum: joinNum,
	})
	if err != nil {
		panic(err)
	}
}

func printVersionInfo() error {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("no build info")
	}

	_, err := io.WriteString(os.Stderr, fmt.Sprintf("Version: %s\n", buildInfo.Main.Version))
	if err != nil {
		return fmt.Errorf("print version info: %w", err)
	}

	return nil
}

func openSource(source string) (io.ReadCloser, error) {
	if len(source) == 0 {
		return nil, errors.New("empty source file")
	}

	file, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("open source: %w", err)
	}

	return file, nil
}
