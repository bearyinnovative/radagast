package config

import (
	"fmt"
	"strconv"
	"strings"
)

// Config provides a configuration repository.
//
// Config keys are string and config value can be retrieved by:
//
//      c.Get("a.b.c").Int()  // or other types
type Config map[string]interface{}

// NewFromMap creates a config instance from a map.
func NewFromMap(c map[string]interface{}) Config {
	return Config(c)
}

// Get returns parsed Result from dotted path.
//
//      c.Get("a.b.c")
func (c Config) Get(path string) Result {
	r := makeResult(c)
	for _, key := range strings.Split(path, ".") {
		if r.Type != ResultTypeConfig {
			// not a config map, cannot continue
			return resultNil
		}

		config := r.Config()
		if value, exists := config[key]; !exists {
			return resultNil
		} else {
			r = makeResult(value)
		}
	}
	return r
}

// GetSlice returns slice of Reults from dotted path.
//
//      c.GetSlice("a.b.c")
func (c Config) GetSlice(path string) (results []Result) {
	var iv interface{} = c
	for _, key := range strings.Split(path, ".") {
		config, isConfig := iv.(Config)
		if !isConfig {
			// not a config map, cannot continue
			return
		}

		if v, exists := config[key]; !exists {
			return
		} else {
			iv = v
		}
	}

	switch slice := iv.(type) {
	case []interface{}:
		for _, item := range slice {
			results = append(results, makeResult(item))
		}
	case []bool:
		for _, item := range slice {
			results = append(results, makeResult(item))
		}
	case []int:
		for _, item := range slice {
			results = append(results, makeResult(item))
		}
	case []string:
		for _, item := range slice {
			results = append(results, makeResult(item))
		}
	case []Config:
		for _, item := range slice {
			results = append(results, makeResult(item))
		}
	default:
	}

	return
}

type ResultType int

const (
	ResultTypeUnknown ResultType = iota
	ResultTypeNil
	ResultTypeBool
	ResultTypeInt
	ResultTypeString
	ResultTypeConfig
)

func (t ResultType) String() string {
	switch t {
	case ResultTypeUnknown:
		return "interface{}"
	case ResultTypeNil:
		return "nil"
	case ResultTypeBool:
		return "bool"
	case ResultTypeInt:
		return "int"
	case ResultTypeString:
		return "string"
	case ResultTypeConfig:
		return "Config"
	default:
		return ""
	}
}

// Result wraps a value.
type Result struct {
	Type ResultType

	value       interface{}
	valueBool   *bool
	valueInt    *int
	valueString *string
	valueConfig *Config
}

var resultNil = makeResult(nil)

func makeResult(raw interface{}) (r Result) {
	r.value = raw
	switch v := raw.(type) {
	case nil:
		r.Type = ResultTypeNil
	case bool:
		r.Type = ResultTypeBool
		r.valueBool = &v
	case int:
		r.Type = ResultTypeInt
		r.valueInt = &v
	case string:
		r.Type = ResultTypeString
		r.valueString = &v
	case Config:
		r.Type = ResultTypeConfig
		r.valueConfig = &v
	default:
		r.Type = ResultTypeUnknown
	}

	return
}

func (r Result) Value() interface{} {
	return r.value
}

func (r Result) NilE() (value interface{}, err error) {
	switch r.Type {
	case ResultTypeNil:
		value = nil
	default:
		err = fmt.Errorf("value is not nil, is %s", r.Type)
	}

	return
}

func (r Result) Nil() interface{} {
	v, _ := r.NilE()
	return v
}

func (r Result) BoolE() (value bool, err error) {
	switch r.Type {
	case ResultTypeNil:
		value = false
	case ResultTypeBool:
		value = *r.valueBool
	case ResultTypeInt:
		value = *r.valueInt != 0
	case ResultTypeString:
		value, err = strconv.ParseBool(*r.valueString)
	default:
		err = fmt.Errorf("value is not bool, is %s", r.Type)
	}

	return
}

func (r Result) Bool() bool {
	v, _ := r.BoolE()
	return v
}

func (r Result) IntE() (value int, err error) {
	switch r.Type {
	case ResultTypeInt:
		value = *r.valueInt
	case ResultTypeString:
		v, err := strconv.ParseInt(*r.valueString, 0, 0)
		if err == nil {
			value = int(v)
		}
	default:
		err = fmt.Errorf("value is not int, is %s", r.Type)
	}

	return
}

func (r Result) Int() int {
	v, _ := r.IntE()
	return v
}

func (r Result) StringE() (value string, err error) {
	switch r.Type {
	case ResultTypeNil:
		value = ""
	case ResultTypeBool:
		value = strconv.FormatBool(*r.valueBool)
	case ResultTypeInt:
		value = strconv.FormatInt(int64(*r.valueInt), 10)
	case ResultTypeString:
		value = *r.valueString
	default:
		err = fmt.Errorf("value is not string, is %s", r.Type)
	}

	return
}

func (r Result) String() string {
	v, _ := r.StringE()
	return v
}

func (r Result) ConfigE() (value Config, err error) {
	switch r.Type {
	case ResultTypeConfig:
		value = *r.valueConfig
	default:
		err = fmt.Errorf("value is not Config, is %s", r.Type)
	}

	return
}

func (r Result) Config() Config {
	v, _ := r.ConfigE()
	return v
}
