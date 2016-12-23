package sflags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter_Set(t *testing.T) {
	var err error
	initial := 0
	counter := (*Counter)(&initial)

	assert.Equal(t, 0, initial)
	assert.Equal(t, "0", counter.String())
	assert.Equal(t, 0, counter.Get())
	assert.Equal(t, "count", counter.Type())
	assert.Equal(t, true, counter.IsBoolFlag())
	assert.Equal(t, true, counter.IsCumulative())

	err = counter.Set("")
	assert.NoError(t, err)
	assert.Equal(t, 1, initial)
	assert.Equal(t, "1", counter.String())

	err = counter.Set("10")
	assert.NoError(t, err)
	assert.Equal(t, 10, initial)
	assert.Equal(t, "10", counter.String())

	err = counter.Set("-1")
	assert.NoError(t, err)
	assert.Equal(t, 11, initial)
	assert.Equal(t, "11", counter.String())

	err = counter.Set("b")
	assert.Error(t, err, "strconv.ParseInt: parsing \"b\": invalid syntax")
	assert.Equal(t, 11, initial)
	assert.Equal(t, "11", counter.String())
}

func TestBoolValue_IsBoolFlag(t *testing.T) {
	b := &boolValue{}
	assert.True(t, b.IsBoolFlag())
}
