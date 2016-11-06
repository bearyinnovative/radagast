package config

import (
	"testing"
	"time"
)

func TestResult_Value(t *testing.T) {
	things := []interface{}{
		"1",
		1,
		nil,
		true,
		false,
		42.0,
		time.Now(),
	}

	for _, thing := range things {
		r := makeResult(thing)
		if v := r.Value(); v != thing {
			t.Errorf("expected: %+v, got: %+v", thing, v)
		}
	}
}

func TestResult_Type(t *testing.T) {
	cases := []struct {
		rtype ResultType
		value interface{}
	}{
		// no supported types
		{ResultTypeUnknown, time.Now()},
		{ResultTypeUnknown, 42.0},
		{ResultTypeUnknown, int64(42)},

		{ResultTypeNil, nil},

		{ResultTypeBool, true},
		{ResultTypeBool, false},

		{ResultTypeInt, 42},
		{ResultTypeInt, 0},
		{ResultTypeInt, 1000},

		{ResultTypeString, "foo"},
		{ResultTypeString, ""},
		{ResultTypeString, "42"},

		{ResultTypeConfig, Config{}},
	}

	for _, testcase := range cases {
		result := makeResult(testcase.value)
		if result.Type != testcase.rtype {
			t.Errorf("expected %s, got: %s", testcase.rtype, result.Type)
		}
	}
}

func TestResult_Nil(t *testing.T) {
	cases := []struct {
		raw           interface{}
		expectedValue interface{}
		hasError      bool
	}{
		{nil, nil, false},
		{true, nil, true},
		{false, nil, true},
		{42, nil, true},
		{0, nil, true},
		{1000, nil, true},
		{"foo", nil, true},
		{"", nil, true},
		{"42", nil, true},
		{Config{}, nil, true},
		{time.Now(), nil, true},
		{42.0, nil, true},
		{int64(42), nil, true},
	}

	for _, testcase := range cases {
		result := makeResult(testcase.raw)
		value, err := result.NilE()

		if value != testcase.expectedValue {
			t.Errorf(
				"expected %+v, got: %+v",
				testcase.expectedValue,
				value,
			)
		}

		if value != result.Nil() {
			t.Fatalf("should unbox value")
		}

		if hasError := err != nil; hasError != testcase.hasError {
			t.Errorf(
				"expected error: %b, got: %b",
				testcase.hasError,
				hasError,
			)
		}
	}
}

func TestResult_Bool(t *testing.T) {
	cases := []struct {
		raw           interface{}
		expectedValue bool
		hasError      bool
	}{
		{nil, false, false},
		{true, true, false},
		{false, false, false},
		{42, true, false},
		{0, false, false},
		{1000, true, false},
		{"foo", false, true},
		{"1", true, false},
		{"true", true, false},
		{"0", false, false},
		{"false", false, false},
		{Config{}, false, true},
		{time.Now(), false, true},
		{42.0, false, true},
		{int64(42), false, true},
	}

	for _, testcase := range cases {
		result := makeResult(testcase.raw)
		value, err := result.BoolE()

		if value != testcase.expectedValue {
			t.Errorf(
				"expected %+v, got: %+v",
				testcase.expectedValue,
				value,
			)
		}

		if value != result.Bool() {
			t.Fatalf("should unbox value")
		}

		if hasError := err != nil; hasError != testcase.hasError {
			t.Errorf(
				"expected error: %t, got: %t",
				testcase.hasError,
				hasError,
			)
		}
	}
}

func TestResult_Int(t *testing.T) {
	cases := []struct {
		raw           interface{}
		expectedValue int
		hasError      bool
	}{
		{nil, 0, true},
		{true, 0, true},
		{false, 0, true},
		{42, 42, false},
		{0, 0, false},
		{1000, 1000, false},
		{"foo", 0, false},
		{"0", 0, false},
		{"42", 42, false},
		{Config{}, 0, true},
		{time.Now(), 0, true},
		{42.0, 0, true},
		{int64(42), 0, true},
	}

	for _, testcase := range cases {
		result := makeResult(testcase.raw)
		value, err := result.IntE()

		if value != testcase.expectedValue {
			t.Errorf(
				"expected %+v, got: %+v",
				testcase.expectedValue,
				value,
			)
		}

		if value != result.Int() {
			t.Fatalf("should unbox value")
		}

		if hasError := err != nil; hasError != testcase.hasError {
			t.Errorf(
				"expected error: %t, got: %t",
				testcase.raw,
				testcase.hasError,
				hasError,
			)
		}
	}
}

func TestResult_String(t *testing.T) {
	cases := []struct {
		raw           interface{}
		expectedValue string
		hasError      bool
	}{
		{nil, "", false},
		{true, "true", false},
		{false, "false", false},
		{42, "42", false},
		{0, "0", false},
		{1000, "1000", false},
		{"foo", "foo", false},
		{"0", "0", false},
		{"42", "42", false},
		{Config{}, "", true},
		{time.Now(), "", true},
		{42.0, "", true},
		{int64(42), "", true},
	}

	for _, testcase := range cases {
		result := makeResult(testcase.raw)
		value, err := result.StringE()

		if value != testcase.expectedValue {
			t.Errorf(
				"expected %+v, got: %+v",
				testcase.expectedValue,
				value,
			)
		}
		if value != result.String() {
			t.Fatalf("should unbox value")
		}

		if hasError := err != nil; hasError != testcase.hasError {
			t.Errorf(
				"expected error: %t, got: %t",
				testcase.raw,
				testcase.hasError,
				hasError,
			)
		}
	}
}

func TestResult_Config(t *testing.T) {
	cases := []struct {
		raw      interface{}
		hasError bool
	}{
		{nil, true},
		{true, true},
		{false, true},
		{42, true},
		{0, true},
		{"foo", true},
		{"42", true},
		{Config{}, false},
		{time.Now(), true},
		{42.0, true},
		{int64(42), true},
	}

	for _, testcase := range cases {
		result := makeResult(testcase.raw)
		_, err := result.ConfigE()

		if hasError := err != nil; hasError != testcase.hasError {
			t.Errorf(
				"expected error: %t, got: %t",
				testcase.raw,
				testcase.hasError,
				hasError,
			)
		}
	}
}

func TestConfig_NewFromMap(t *testing.T) {
	c := NewFromMap(map[string]interface{}{"foo": "bar"})
	if c["foo"] != "bar" {
		t.Errorf("NewFromMap failed: %+v", c)
	}
}

func TestConfig_Get(t *testing.T) {
	c := Config{
		"nil":       nil,
		"boolTrue":  true,
		"boolFalse": false,
		"int":       42,
		"string":    "foo",
		"config": Config{
			"nil":       nil,
			"boolTrue":  true,
			"boolFalse": false,
			"int":       42,
			"string":    "foo",
		},
	}

	if v := c.Get("nil").Nil(); v != nil {
		t.Errorf("expected nil, got: %+v", v)
	}

	if v := c.Get("boolTrue").Bool(); v != true {
		t.Errorf("expected true, got: %+v", v)
	}

	if v := c.Get("boolFalse").Bool(); v != false {
		t.Errorf("expected false, got: %+v", v)
	}

	if v := c.Get("int").Int(); v != 42 {
		t.Errorf("expected 42, got: %+v", v)
	}

	if v := c.Get("string").String(); v != "foo" {
		t.Errorf("expected foo, got: %+v", v)
	}

	if v := c.Get("config.nil").Nil(); v != nil {
		t.Errorf("expected nil, got: %+v", v)
	}

	if v := c.Get("config.boolTrue").Bool(); v != true {
		t.Errorf("expected true, got: %+v", v)
	}

	if v := c.Get("config.boolFalse").Bool(); v != false {
		t.Errorf("expected false, got: %+v", v)
	}

	if v := c.Get("config.int").Int(); v != 42 {
		t.Errorf("expected 42, got: %+v", v)
	}

	if v := c.Get("config.string").String(); v != "foo" {
		t.Errorf("expected foo, got: %+v", v)
	}
}

func TestConfig_GetSlice(t *testing.T) {
	c := Config{
		"value": "foo",
		"slice": []string{"1"},
		"config": Config{
			"value": "foo",
			"slice": []string{"1"},
		},
	}

	if v := c.GetSlice("value"); v != nil {
		t.Errorf("unexpected slice: %+v", v)
	}

	if v := c.GetSlice("config.value"); v != nil {
		t.Errorf("unexpected slice: %+v", v)
	}

	if v := c.GetSlice("slice"); len(v) != 1 {
		t.Errorf("should get slice: %+v", v)
	} else if s := v[0].String(); s != "1" {
		t.Errorf("unexpected slice: %+v", v)
	}

	if v := c.GetSlice("config.slice"); len(v) != 1 {
		t.Errorf("should get slice")
	} else if s := v[0].String(); s != "1" {
		t.Errorf("unexpected slice: %+v", v)
	}
}
