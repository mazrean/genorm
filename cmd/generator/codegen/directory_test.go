package codegen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDestination(t *testing.T) {
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
		packageName string
		modulePath  string
		exists      bool
		err         bool
	}{
		{
			description: "normal destination -> success",
			path:        "./test2",
			packageName: "test2",
			modulePath:  "github/mazrean/genorm/test2",
			exists:      true,
		},
		{
			description: "non-existent destination -> success",
			path:        "./test2",
			packageName: "test2",
			modulePath:  "github/mazrean/genorm/test2",
		},
		{
			description: "destination in directory -> success",
			path:        "./test/test",
			packageName: "test",
			modulePath:  "github/mazrean/genorm/test",
			exists:      true,
		},
		{
			description: "non-existent destination in directory -> success",
			path:        "./test/test",
			packageName: "test",
			modulePath:  "github/mazrean/genorm/test",
		},
		{
			description: "destination(absolute path) -> success",
			path:        absPath,
			packageName: "test2",
			modulePath:  "github/mazrean/genorm/test2",
			exists:      true,
		},
		{
			description: "non-existent destination(absolute path) -> success",
			path:        absPath,
			packageName: "test2",
			modulePath:  "github/mazrean/genorm/test2",
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

			dir, err := newDirectory(test.path, test.packageName, test.modulePath)
			if err != nil {
				if !test.err {
					t.Fatalf("unexpected error: %s", err)
				}
				return
			}

			if test.err {
				t.Fatalf("expected error but got none")
			}

			assert.Equal(t, test.path, dir.path)
			assert.Equal(t, test.packageName, dir.packageName)
			assert.Equal(t, test.modulePath, dir.modulePath)
			assert.DirExists(t, dir.path)
		})
	}
}
