package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenSource(t *testing.T) {
	err := os.MkdirAll("./test", os.ModePerm)
	if err != nil {
		t.Fatalf("create test directory: %s", err)
	}
	defer func() {
		err := os.RemoveAll("./test")
		if err != nil {
			t.Errorf("failed to remove test directory: %s", err)
		}
	}()

	absPath, err := filepath.Abs("./test.go")
	if err != nil {
		t.Fatalf("failed to get absolute path: %s", err)
	}

	tests := []struct {
		description string
		path        string
		exists      bool
		content     string
		err         bool
	}{
		{
			description: "normal source -> success",
			path:        "./test.go",
			exists:      true,
			content: `package main

func main() {
	fmt.Println("Hello, World!")
}
`,
		},
		{
			description: "source in directory -> success",
			path:        "./test/main.go",
			exists:      true,
			content: `package main

func main() {
	fmt.Println("Hello, World!")
}
`,
		},
		{
			description: "source(absolute path) -> success",
			path:        absPath,
			exists:      true,
			content: `package main

func main() {
	fmt.Println("Hello, World!")
}
`,
		},
		{
			description: "non-existent source -> error",
			path:        "./test/main.go",
			err:         true,
		},
		{
			description: "empty source -> error",
			path:        "",
			err:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			if test.exists {
				func() {
					f, err := os.Create(test.path)
					if err != nil {
						t.Fatalf("failed to create test file: %s", err)
					}
					defer f.Close()

					_, err = f.WriteString(test.content)
					if err != nil {
						t.Fatalf("failed to write test file: %s", err)
					}
				}()

				defer func() {
					err := os.Remove(test.path)
					if err != nil {
						t.Errorf("failed to remove test file: %s", err)
					}
				}()
			}

			src, err := openSource(test.path)
			if err != nil {
				if !test.err {
					t.Fatalf("unexpected error: %s", err)
				}
				return
			}
			defer src.Close()

			if test.err {
				t.Fatalf("expected error but got none")
			}

			if test.exists {
				sb := strings.Builder{}

				_, err := io.Copy(&sb, src)
				if err != nil {
					t.Fatalf("failed to read source: %s", err)
				}

				actualContent := sb.String()
				if actualContent != test.content {
					t.Fatalf("unexpected content: %s", actualContent)
				}
			}
		})
	}
}

func TestDestinationDir(t *testing.T) {
	err := os.MkdirAll("./test", os.ModePerm)
	if err != nil {
		t.Fatalf("create test directory: %s", err)
	}
	defer func() {
		err := os.RemoveAll("./test")
		if err != nil {
			t.Errorf("failed to remove test directory: %s", err)
		}
	}()

	absPath, err := filepath.Abs("./test2")
	if err != nil {
		t.Fatalf("failed to get absolute path: %s", err)
	}

	tests := []struct {
		description string
		path        string
		exists      bool
		err         bool
	}{
		{
			description: "empty destination -> error",
			path:        "",
			err:         true,
		},
		{
			description: "normal destination -> success",
			path:        "./test2",
			exists:      true,
		},
		{
			description: "non-existent destination -> success",
			path:        "./test2",
		},
		{
			description: "destination in directory -> success",
			path:        "./test/test",
			exists:      true,
		},
		{
			description: "non-existent destination in directory -> success",
			path:        "./test/test",
		},
		{
			description: "destination(absolute path) -> success",
			path:        absPath,
			exists:      true,
		},
		{
			description: "non-existent destination(absolute path) -> success",
			path:        absPath,
			exists:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			if test.exists {
				err := os.MkdirAll(test.path, os.ModePerm)
				if err != nil {
					t.Fatalf("failed to create test file: %s", err)
				}

				defer func() {
					err := os.RemoveAll(test.path)
					if err != nil {
						t.Errorf("failed to remove test file: %s", err)
					}
				}()
			}

			dest, err := destinationDir(test.path)
			if err != nil {
				if !test.err {
					t.Fatalf("unexpected error: %s", err)
				}
				return
			}

			if test.err {
				t.Fatalf("expected error but got none")
			}

			if dest != test.path {
				t.Fatalf("unexpected destination: %s", dest)
			}
		})
	}
}
