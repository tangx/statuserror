package status_error_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-courier/status_error"
	"github.com/go-courier/status_error/examples"
)

func ExampleStatusErr() {
	fmt.Println(examples.Unauthorized)
	fmt.Println(status_error.FromErr(nil))
	fmt.Println(status_error.FromErr(fmt.Errorf("unknown")))
	// Output:
	//[]@StatusErr[Unauthorized][401999001][Unauthorized]!
	//<nil>
	//[]@StatusErr[UnknownError][500000000][unknown error] unknown
	//
}

func TestStatusErr(t *testing.T) {
	summary := status_error.NewUnknownErr().Summary()

	assert.Equal(t, "@StatusErr[UnknownError][500000000][unknown error]", summary)

	statusErr, err := status_error.ParseStatusErrSummary(summary)
	assert.NoError(t, err)

	assert.Equal(t, status_error.NewUnknownErr(), statusErr)

	assert.Equal(t, "@StatusErr[Unauthorized][401999001][Unauthorized]!", examples.Unauthorized.StatusErr().Summary())
	assert.Equal(t, "@StatusErr[InternalServerError][500999001][InternalServerError]", examples.InternalServerError.StatusErr().Summary())

	assert.Equal(t, 401, examples.Unauthorized.StatusCode())
	assert.Equal(t, 401, examples.Unauthorized.StatusErr().StatusCode())

	assert.True(t, examples.Unauthorized.StatusErr().Is(examples.Unauthorized))
	assert.True(t, examples.Unauthorized.StatusErr().Is(examples.Unauthorized.StatusErr()))
}

func TestStatusErrBuilders(t *testing.T) {
	t.Log(examples.Unauthorized.StatusErr().WithMsg("msg overwrite"))
	t.Log(examples.Unauthorized.StatusErr().WithDesc("desc overwrite"))
	t.Log(examples.Unauthorized.StatusErr().DisableErrTalk().EnableErrTalk())
	t.Log(examples.Unauthorized.StatusErr().WithID("111"))
	t.Log(examples.Unauthorized.StatusErr().AppendSource("service-abc"))
	t.Log(examples.Unauthorized.StatusErr().AppendErrorField("header", "Authorization", "missing"))
	t.Log(examples.Unauthorized.StatusErr().AppendErrorFields(
		status_error.NewErrorField("query", "key", "missing"),
		status_error.NewErrorField("header", "Authorization", "missing"),
	))
}
