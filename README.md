# sflags - Generate flags by parsing structures. [![GoDoc](https://godoc.org/github.com/octago/sflags?status.svg)](http://godoc.org/github.com/octago/sflags) [![Build Status](https://travis-ci.org/octago/sflags.svg?branch=master)](https://travis-ci.org/octago/sflags)  [![codecov](https://codecov.io/gh/octago/sflags/branch/master/graph/badge.svg)](https://codecov.io/gh/octago/sflags)
 [![Go Report Card](https://goreportcard.com/badge/github.com/octago/sflags)](https://goreportcard.com/report/github.com/octago/sflags)


Look at the examples in examples folder for different flag libraries.

## Options for flag tag

The flag default key string is the struct field name but can be specified in the struct field's tag value.
The "flag" key in the struct field's tag value is the key name, followed by an optional comma and options. Examples:
```
// Field is ignored by this package.
Field int `flag:"-"`

// Field appears in flags as "myName".
Field int `flag:"myName"`

// If this field is from nested struct, prefix from parent struct will be ingored.
Field int `flag:"~myName"`

// You can set short name for flags by providing it's value after a space
// Prefixes will not be applied for short names.
Field int `flag:"myName a"`

// this field will be removed from generated help text.
Field int `flag:",hidden"`

// this field will be marked as deprecated in generated help text
Field int `flag:",deprecated"`
```

## Options for desc tag
If you specify description in description tag (`desc` by default) it will be used in USAGE section.

```
Addr string `desc:"HTTP address"`
```
this description produces something like:
```
  -addr value
    	HTTP host (default 127.0.0.1)
```

## Options for env tag


## Options for Parse function:

```
// DescTag sets custom description tag. It is "desc" by default.
func DescTag(val string)

// FlagTag sets custom flag tag. It is "flag" be default.
func FlagTag(val string)

// Prefix sets prefix that will be applied for all flags (if they are not marked as ~).
func Prefix(val string)

// EnvPrefix sets prefix that will be applied for all environment variables (if they are not marked as ~).
func EnvPrefix(val string)

// FlagDivider sets custom divider for flags. It is dash by default. e.g. "flag-name".
func FlagDivider(val string)

// EnvDivider sets custom divider for environment variables.
// It is underscore by default. e.g. "ENV_NAME".
func EnvDivider(val string)
```



Features:

 - [ ] Set environment name
 - [ ] Set usage
 - [ ] Long and short forms
 - [ ] Skip field
 - [ ] Required
 - [ ] Placeholders (by `name`)
 - [ ] Deprecated and hidden options
 - [ ] Multiple ENV names
 - [ ] Interface for user types.


Supported types in structures:

 - [x] `int`, `int8`, `int16`, `int32`, `int64`
 - [x] `uint`, `uint8`, `uint16`, `uint32`, `uint64`
 - [x] `float32`, `float64`
 - [x] slices for all previous numeric types (e.g. `[]int`, `[]float64`)
 - [x] `bool`
 - [x] `[]bool`
 - [x] `string`
 - [x] `[]string`
 - [x] nested structures
 - [x] net.TCPAddr
 - [x] net.IP
 - [x] time.Duration
 - [x] regexp.Regexp


Custom types:
 - [x] HexBytes

 - [x] count
 - [ ] ipmask
 - [ ] enum values
 - [ ] enum list values
 - [ ] file
 - [ ] file list
 - [ ] url
 - [ ] url list
 - [ ] units (bytes 1kb = 1024b, speed, etc)

Supported flags and cli libraries:

 - [ ] [flag](https://golang.org/pkg/flag/)
 - [ ] [spf13/pflag](https://github.com/spf13/pflag)
 - [ ] [spf13/cobra](https://github.com/spf13/cobra)
 - [ ] [urfave/cli](https://github.com/urfave/cli)

Matrix:

| Name | Hidden | Deprecated | Short | Env |
|----|----|----|----|
| flag | - | - | - | - |
| pflag | [ ] | [ ] | [ ] | [ ] |
| kingpin | [ ] | [ ] | [ ] | [ ] |
| urfave | [ ] | [ ] | [ ] | [ ] |