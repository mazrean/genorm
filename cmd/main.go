package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mazrean/genorm/cmd/generator"
)

var (
	// Set at build time.
	version string
	commit  string
	date    string

	// flags
	showVersionInfo bool
	source          string
	destination     string
	packageName     string
	moduleName      string
)

func init() {
	flag.BoolVar(&showVersionInfo, "version", false, "If true, output version information.")
	flag.StringVar(&source, "source", "", "The source file to parse.")
	flag.StringVar(&destination, "destination", "", "The destination file to write.")
	flag.StringVar(&packageName, "package", "", "The root package name to use.")
	flag.StringVar(&moduleName, "module", "", "The root module name to use.")
}

func main() {
	flag.Parse()

	if showVersionInfo {
		err := printVersionInfo(version, commit, date)
		if err != nil {
			panic(err)
		}
	}

	src, err := openSource(source)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := destinationDir(destination)
	if err != nil {
		panic(err)
	}

	if len(packageName) == 0 {
		panic("package name is required")
	}
	if len(moduleName) == 0 {
		panic("module name is required")
	}

	err = generator.Generate(packageName, moduleName, dst, src)
	if err != nil {
		panic(err)
	}
}

func printVersionInfo(version string, commit string, date string) error {
	_, err := io.WriteString(os.Stderr, fmt.Sprintf(`Version: %s
Commit: %s
Date: %s
`, version, commit, date))
	if err != nil {
		return fmt.Errorf("print version info: %w", err)
	}

	return nil
}

func openSource(source string) (io.ReadCloser, error) {
	if len(source) == 0 {
		return nil, errors.New("Source file is required.")
	}

	file, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("open source: %w", err)
	}

	return file, nil
}

func destinationDir(destination string) (string, error) {
	if len(destination) == 0 {
		return "", errors.New("Destination directory path is required.")
	}

	err := os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	return destination, nil
}
