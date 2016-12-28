package sflags

import (
	"fmt"
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

func TestValidateValue_IsBoolFlag(t *testing.T) {
	boolV := true
	v := &validateValue{Value: newBoolValue(&boolV)}
	assert.True(t, v.IsBoolFlag())

	v = &validateValue{Value: newStringValue(strP("stringValue"))}
	assert.False(t, v.IsBoolFlag())
}

func TestValidateValue_IsCumulative(t *testing.T) {
	v := &validateValue{Value: newStringValue(strP("stringValue"))}
	assert.False(t, v.IsCumulative())

	v = &validateValue{Value: newStringSliceValue(&[]string{})}
	assert.True(t, v.IsCumulative())
}

func TestValidateValue_String(t *testing.T) {
	v := &validateValue{Value: newStringValue(strP("stringValue"))}
	assert.Equal(t, "stringValue", v.String())

	v = &validateValue{Value: nil}
	assert.Equal(t, "", v.String())
}

func TestValidateValue_Set(t *testing.T) {
	sV := strP("stringValue")
	v := &validateValue{Value: newStringValue(sV)}
	assert.NoError(t, v.Set("newVal"))
	assert.Equal(t, "newVal", *sV)

	v.validateFunc = func(val string) error {
		return nil
	}
	assert.NoError(t, v.Set("newVal"))

	v.validateFunc = func(val string) error {
		return fmt.Errorf("invalid %s", val)
	}
	assert.EqualError(t, v.Set("newVal"), "invalid newVal")
}
