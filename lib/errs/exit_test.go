package errs_test

import (
	"errors"
	"testing"

	"github.com/marwanhawari/stew/lib/errs"
)

func TestMaybeExit(t *testing.T) {
	t.Parallel()
	tests := []maybeExitTestCase{{
		name: "nil",
		err:  nil,
		want: 0,
	}, {
		name: "error",
		err:  errors.New("test"),
		want: 215,
		exit: true,
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, tt.test)
	}
}

type maybeExitTestCase struct {
	name string
	err  error
	want int
	exit bool
}

func (tc maybeExitTestCase) test(t *testing.T) {
	t.Helper()
	var (
		code    int
		exit    bool
		handler = func(c int) {
			code = c
			exit = true
		}
	)
	errs.WithExitFn(handler, func() {
		errs.MaybeExit(tc.err)
	})
	if exit != tc.exit {
		t.Errorf("exit mismatch:\n"+
			"want = %t\n"+
			" got = %t", tc.exit, exit)
	}
	if tc.want != code {
		t.Errorf("retcode mismatch:\n"+
			"want = %d\n"+
			" got = %d", tc.want, code)
	}
}
