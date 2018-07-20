package statuserror_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

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

	assert.Equal(t, "@StatusErr[UnknownError][500000000][unknown error]", summary)

	statusErr, err := statuserror.ParseStatusErrSummary(summary)
	assert.NoError(t, err)

	assert.Equal(t, statuserror.NewUnknownErr(), statusErr)

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
		statuserror.NewErrorField("query", "key", "missing"),
		statuserror.NewErrorField("header", "Authorization", "missing"),
	))
}
