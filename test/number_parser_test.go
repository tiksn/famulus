package phone

import (
	"testing"

	"github.com/tiksn/famulus/internal/pkg/phone"
)

func TestParse(t *testing.T) {
	nums, err := phone.Parse("0689546321", "UA")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log(nums)
}
