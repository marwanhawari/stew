package cmd_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/marwanhawari/stew/lib/cmd"
	"github.com/marwanhawari/stew/lib/testsupport"
)

func TestInstallBatch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	tests := []installBatchTestCase{{
		name: "ppath as foo",
		args: []string{"marwanhawari/ppath!!foo"},
		want: "foo",
	}, {
		name:      "ppath with batch mode",
		args:      []string{"marwanhawari/ppath@v0.0.2"},
		want:      "ppath",
		batchMode: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

type installBatchTestCase struct {
	name      string
	args      []string
	batchMode bool
	wantErr   error
	want      string
}

func (tc installBatchTestCase) test(t *testing.T) {
	rt := testsupport.NewDefaultRuntime(t)
	if err := cmd.Install(rt, tc.args, tc.batchMode); !errors.Is(err, tc.wantErr) {
		t.Errorf("error mismatch:\n"+
			" got = %+v\n"+
			"want = %+v", err, tc.wantErr)
		t.FailNow()
	}
	binary := path.Join(rt.StewBinPath, tc.want)

	fi, err := os.Stat(binary)
	if err != nil {
		t.Errorf("unexpected error: %+v", err)
		t.FailNow()
	}
	if fi.Size() < 1_000_000 {
		t.Errorf("unexpected binary size: %d", fi.Size())
	}
}
