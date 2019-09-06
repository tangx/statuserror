package statuserror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-courier/statuserror"
	"github.com/go-courier/statuserror/__examples__"
)

func ExampleStatusErr() {
	fmt.Println(examples.Unauthorized)
	fmt.Println(statuserror.FromErr(nil))
	fmt.Println(statuserror.FromErr(fmt.Errorf("unknown")))
	// Output:
	//[]@StatusErr[Unauthorized][401999001][Unauthorized]!
	//<nil>
	//[]@StatusErr[UnknownError][500000000][unknown error] unknown
	//
}

func TestStatusErr(t *testing.T) {
	summary := statuserror.NewUnknownErr().Summary()

	require.Equal(t, "@StatusErr[UnknownError][500000000][unknown error]", summary)

	statusErr, err := statuserror.ParseStatusErrSummary(summary)
	require.NoError(t, err)

	require.Equal(t, statuserror.NewUnknownErr(), statusErr)

	require.Equal(t, "@StatusErr[Unauthorized][401999001][Unauthorized]!", examples.Unauthorized.StatusErr().Summary())
	require.Equal(t, "@StatusErr[InternalServerError][500999001][InternalServerError]", examples.InternalServerError.StatusErr().Summary())

	require.Equal(t, 401, examples.Unauthorized.StatusCode())
	require.Equal(t, 401, examples.Unauthorized.StatusErr().StatusCode())

	require.True(t, errors.Is(examples.Unauthorized.StatusErr(), examples.Unauthorized))
	require.True(t, errors.Is(examples.Unauthorized.StatusErr(), examples.Unauthorized.StatusErr()))
}

func TestStatusErrBuilders(t *testing.T) {
	t.Log(examples.Unauthorized.StatusErr().WithMsg("msg overwrite"))
	t.Log(examples.Unauthorized.StatusErr().WithDesc("desc overwrite"))
	t.Log(examples.Unauthorized.StatusErr().DisableErrTalk().EnableErrTalk())
	t.Log(examples.Unauthorized.StatusErr().WithID("111"))
	t.Log(examples.Unauthorized.StatusErr().AppendSource("service-abc"))
	t.Log(examples.Unauthorized.StatusErr().AppendErrorField("header", "Authorization", "missing"))
	t.Log(examples.Unauthorized.StatusErr().AppendErrorFields(
		statuserror.NewErrorField("query", "key", "missing"),
		statuserror.NewErrorField("header", "Authorization", "missing"),
	))
}
