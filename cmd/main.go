package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
)

func init() {
	flag.BoolVar(&showVersionInfo, "version", false, "If true, output version information.")
	flag.StringVar(&source, "source", "", "The source file to parse.")
	flag.StringVar(&destination, "destination", "", "The destination file to write.")
	flag.StringVar(&packageName, "package", "", "The package name to use.")
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

	dst, err := openDestination(destination)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	err = generator.Generate(packageName, src, dst)
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

func openDestination(destination string) (io.WriteCloser, error) {
	if len(destination) == 0 {
		return os.Stdout, nil
	}

	destinationDir := filepath.Dir(destination)
	_, err := os.Stat(destinationDir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("failed to get destination directory info: %w", err)
	}

	if errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(destinationDir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("failed to create destination directory: %w", err)
		}
	}

	file, err := os.Create(destination)
	if err != nil {
		return nil, fmt.Errorf("create destination: %w", err)
	}

	return file, nil
}
