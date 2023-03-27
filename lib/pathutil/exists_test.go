package pathutil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marwanhawari/stew/lib/pathutil"
)

func TestPathExists(t *testing.T) {
	type TestCase struct {
		name string
		path string
		want bool
	}
	cases := []TestCase{{
		name: "test1",
		path: "testFile",
		want: true,
	}, {
		name: "test2",
		path: "noFile",
	}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testFilePath := filepath.Join(tempDir, tt.path)
			if tt.want {
				os.WriteFile(testFilePath, []byte("A test file"), 0o644)
			}

			got, err := pathutil.Exists(testFilePath)
			if err != nil {
				t.Fatalf("pathutil.Exists() unexpected error = %+v\n", err)
			}
			if got != tt.want {
				t.Errorf("pathutil.Exists()\n"+
					" got = %v"+
					"want = %v", got, tt.want)
			}
		})
	}
}
