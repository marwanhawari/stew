package prompt_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/marwanhawari/stew/constants"
	"github.com/marwanhawari/stew/lib/ui/prompt"
)

func TestExitUserSelectionError_Error(t *testing.T) {
	type TestCase struct {
		name string
		err  error
		want string
	}
	cases := []TestCase{{
		name: "test1",
		err:  errors.New("testErr"),
		want: fmt.Sprintf("%v Exited from user selection: %v",
			constants.RedColor("Error:"),
			constants.RedColor(errors.New("testErr"))),
	}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			e := prompt.ExitUserSelectionError{
				Err: tt.err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ExitUserSelectionError.Error()\n"+
					" got = %+v\n"+
					"want = %+v", got, tt.want)
			}
		})
	}
}
