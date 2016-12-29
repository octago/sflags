package sflags

import (
	"errors"
	"net"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func strP(value string) *string {
	return &value
}

type simple struct {
	Name string
}

func TestParseStruct(t *testing.T) {
	simpleCfg := &struct {
		Name  string `desc:"name description" env:"-"`
		Name2 string `flag:"name_two t,hidden,deprecated"`
		Name3 string `env:"NAME_THREE"`
		Name4 *string
		Name5 string `flag:"-"`
		name6 string

		Addr *net.TCPAddr
	}{
		Name:  "name_value",
		Name2: "name2_value",
		Name4: strP("name_value4"),
		Addr: &net.TCPAddr{
			IP: net.ParseIP("127.0.0.1"),
		},
	}
	diffTypesCfg := &struct {
		StringValue      string
		ByteValue        byte
		StringSliceValue []string
		BoolSliceValue   []bool
		CounterValue     Counter
		RegexpValue      *regexp.Regexp
		FuncValue        func() // will be ignored
	}{
		StringValue:      "string",
		ByteValue:        10,
		StringSliceValue: []string{},
		BoolSliceValue:   []bool{},
		CounterValue:     10,
		RegexpValue:      &regexp.Regexp{},
	}
	nestedCfg := &struct {
		Sub struct {
			Name  string `desc:"name description"`
			Name2 string `env:"NAME_TWO"`
			Name3 string `flag:"~name3" env:"~NAME_THREE"`
			SUB2  *struct {
				Name4 string
				Name5 string `env:"name_five"`
			}
		}
	}{
		Sub: struct {
			Name  string `desc:"name description"`
			Name2 string `env:"NAME_TWO"`
			Name3 string `flag:"~name3" env:"~NAME_THREE"`
			SUB2  *struct {
				Name4 string
				Name5 string `env:"name_five"`
			}
		}{
			Name:  "name_value",
			Name2: "name2_value",
			SUB2: &struct {
				Name4 string
				Name5 string `env:"name_five"`
			}{
				Name4: "name4_value",
			},
		},
	}
	descCfg := &struct {
		Name  string `desc:"name description"`
		Name2 string `description:"name2 description"`
	}{}
	anonymousCfg := &struct {
		Name1 string
		simple
	}{
		simple: simple{
			Name: "name_value",
		},
	}

	tt := []struct {
		name string

		cfg        interface{}
		optFuncs   []OptFunc
		expFlagSet []*Flag
		expErr     error
	}{
		{
			name: "SimpleCfg test",
			cfg:  simpleCfg,
			expFlagSet: []*Flag{
				{
					Name:     "name",
					EnvName:  "",
					DefValue: "name_value",
					Value:    newStringValue(&simpleCfg.Name),
					Usage:    "name description",
				},
				{
					Name:       "name_two",
					Short:      "t",
					EnvName:    "NAME_TWO",
					DefValue:   "name2_value",
					Value:      newStringValue(&simpleCfg.Name2),
					Hidden:     true,
					Deprecated: true,
				},
				{
					Name:     "name3",
					EnvName:  "NAME_THREE",
					DefValue: "",
					Value:    newStringValue(&simpleCfg.Name3),
				},
				{
					Name:     "name4",
					EnvName:  "NAME4",
					DefValue: "name_value4",
					Value:    newStringValue(simpleCfg.Name4),
				},
				{
					Name:     "addr",
					EnvName:  "ADDR",
					DefValue: "127.0.0.1:0",
					Value:    newTCPAddrValue(simpleCfg.Addr),
				},
			},
		},
		{
			name:     "SimpleCfg test with custom env_prefix and divider",
			cfg:      simpleCfg,
			optFuncs: []OptFunc{EnvPrefix("PP|"), EnvDivider("|")},
			expFlagSet: []*Flag{
				{
					Name:     "name",
					EnvName:  "",
					DefValue: "name_value",
					Value:    newStringValue(&simpleCfg.Name),
					Usage:    "name description",
				},
				{
					Name:       "name_two",
					Short:      "t",
					EnvName:    "PP|NAME_TWO",
					DefValue:   "name2_value",
					Value:      newStringValue(&simpleCfg.Name2),
					Hidden:     true,
					Deprecated: true,
				},
				{
					Name:     "name3",
					EnvName:  "PP|NAME_THREE",
					DefValue: "",
					Value:    newStringValue(&simpleCfg.Name3),
				},
				{
					Name:     "name4",
					EnvName:  "PP|NAME4",
					DefValue: "name_value4",
					Value:    newStringValue(simpleCfg.Name4),
				},
				{
					Name:     "addr",
					EnvName:  "PP|ADDR",
					DefValue: "127.0.0.1:0",
					Value:    newTCPAddrValue(simpleCfg.Addr),
				},
			},
			expErr: nil,
		},
		{
			name: "DifferentTypesCfg",
			cfg:  diffTypesCfg,
			expFlagSet: []*Flag{
				{
					Name:     "string-value",
					EnvName:  "STRING_VALUE",
					DefValue: "string",
					Value:    newStringValue(&diffTypesCfg.StringValue),
					Usage:    "",
				},
				{
					Name:     "byte-value",
					EnvName:  "BYTE_VALUE",
					DefValue: "10",
					Value:    newUint8Value(&diffTypesCfg.ByteValue),
					Usage:    "",
				},
				{
					Name:     "string-slice-value",
					EnvName:  "STRING_SLICE_VALUE",
					DefValue: "[]",
					Value:    newStringSliceValue(&diffTypesCfg.StringSliceValue),
					Usage:    "",
				},
				{
					Name:     "bool-slice-value",
					EnvName:  "BOOL_SLICE_VALUE",
					DefValue: "[]",
					Value:    newBoolSliceValue(&diffTypesCfg.BoolSliceValue),
					Usage:    "",
				},
				{
					Name:     "counter-value",
					EnvName:  "COUNTER_VALUE",
					DefValue: "10",
					Value:    &diffTypesCfg.CounterValue,
					Usage:    "",
				},
				{
					Name:     "regexp-value",
					EnvName:  "REGEXP_VALUE",
					DefValue: "",
					Value:    newRegexpValue(&diffTypesCfg.RegexpValue),
					Usage:    "",
				},
			},
		},
		{
			name: "NestedCfg",
			cfg:  nestedCfg,
			expFlagSet: []*Flag{
				{
					Name:     "sub-name",
					EnvName:  "SUB_NAME",
					DefValue: "name_value",
					Value:    newStringValue(&nestedCfg.Sub.Name),
					Usage:    "name description",
				},
				{
					Name:     "sub-name2",
					EnvName:  "SUB_NAME_TWO",
					DefValue: "name2_value",
					Value:    newStringValue(&nestedCfg.Sub.Name2),
				},
				{
					Name:     "name3",
					EnvName:  "NAME_THREE",
					DefValue: "",
					Value:    newStringValue(&nestedCfg.Sub.Name3),
				},
				{
					Name:     "sub-sub2-name4",
					EnvName:  "SUB_SUB2_NAME4",
					DefValue: "name4_value",
					Value:    newStringValue(&nestedCfg.Sub.SUB2.Name4),
				},
				{
					Name:     "sub-sub2-name5",
					EnvName:  "SUB_SUB2_name_five",
					DefValue: "",
					Value:    newStringValue(&nestedCfg.Sub.SUB2.Name5),
				},
			},
			expErr: nil,
		},
		{
			name:     "DescCfg with custom desc tag",
			cfg:      descCfg,
			optFuncs: []OptFunc{DescTag("description")},
			expFlagSet: []*Flag{
				{
					Name:    "name",
					EnvName: "NAME",
					Value:   newStringValue(&descCfg.Name),
				},
				{
					Name:    "name2",
					EnvName: "NAME2",
					Value:   newStringValue(&descCfg.Name2),
					Usage:   "name2 description",
				},
			},
		},
		{
			name: "Anonymoust cfg with disabled flatten",
			cfg:  anonymousCfg,
			expFlagSet: []*Flag{
				{
					Name:    "name1",
					EnvName: "NAME1",
					Value:   newStringValue(&anonymousCfg.Name1),
				},
				{
					Name:     "name",
					EnvName:  "NAME",
					DefValue: "name_value",
					Value:    newStringValue(&anonymousCfg.Name),
				},
			},
		},
		{
			name:     "Anonymoust cfg with enabled flatten",
			cfg:      anonymousCfg,
			optFuncs: []OptFunc{Flatten(false)},
			expFlagSet: []*Flag{
				{
					Name:    "name1",
					EnvName: "NAME1",
					Value:   newStringValue(&anonymousCfg.Name1),
				},
				{
					Name:     "simple-name",
					EnvName:  "SIMPLE_NAME",
					DefValue: "name_value",
					Value:    newStringValue(&anonymousCfg.Name),
				},
			},
		},
		{
			name: "We need pointer to structure",
			cfg: struct {
			}{},
			expErr: errors.New("object must be a pointer to struct or interface"),
		},
		{
			name:   "We need pointer to structure 2",
			cfg:    strP("something"),
			expErr: errors.New("object must be a pointer to struct or interface"),
		},
		{
			name:   "We need non nil object",
			cfg:    nil,
			expErr: errors.New("object cannot be nil"),
		},
		{
			name:   "We need non nil value",
			cfg:    (*simple)(nil),
			expErr: errors.New("object cannot be nil"),
		},
	}
	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			flagSet, err := ParseStruct(test.cfg, test.optFuncs...)
			if test.expErr == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, test.expErr, err)
			}
			assert.Equal(t, test.expFlagSet, flagSet)
		})
	}
}

func TestParseStruct_NilValue(t *testing.T) {
	name2Value := "name2_value"
	cfg := struct {
		Name1  *string
		Name2  *string
		Regexp *regexp.Regexp
	}{
		Name2: &name2Value,
	}
	assert.Nil(t, cfg.Name1)
	assert.Nil(t, cfg.Regexp)
	assert.NotNil(t, cfg.Name2)

	flags, err := ParseStruct(&cfg)
	require.NoError(t, err)
	require.Equal(t, 3, len(flags))
	assert.NotNil(t, cfg.Name1)
	assert.NotNil(t, cfg.Name2)
	assert.NotNil(t, cfg.Regexp)
	assert.Equal(t, name2Value, flags[1].Value.(Getter).Get())

	err = flags[0].Value.Set("name1value")
	require.NoError(t, err)
	assert.Equal(t, "name1value", *cfg.Name1)

	err = flags[2].Value.Set("aabbcc")
	require.NoError(t, err)
	assert.Equal(t, "aabbcc", cfg.Regexp.String())
}

func TestFlagDivider(t *testing.T) {
	opt := opts{
		flagDivider: "-",
	}
	FlagDivider("_")(&opt)
	assert.Equal(t, "_", opt.flagDivider)
}

func TestFlagTag(t *testing.T) {
	opt := opts{
		flagTag: "flags",
	}
	FlagTag("superflag")(&opt)
	assert.Equal(t, "superflag", opt.flagTag)
}

func TestValidator(t *testing.T) {
	var vf ValidateFunc
	opt := opts{
		validator: nil,
	}
	Validator(vf)(&opt)
	assert.Equal(t, vf, opt.validator)
}

func TestFlatten(t *testing.T) {
	opt := opts{
		flatten: true,
	}
	Flatten(false)(&opt)
	assert.Equal(t, false, opt.flatten)
}
