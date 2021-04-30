package pushover

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := &Error{
		Err: context.DeadlineExceeded,
	}

	assert.Equal(t, context.DeadlineExceeded.Error(), err.Error())
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
	assert.Equal(t, context.DeadlineExceeded, errors.Unwrap(err))
	assert.Equal(t, true, err.Timeout())
	assert.Equal(t, true, err.Temporary())
}
