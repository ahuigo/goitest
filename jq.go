package goitest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/itchyny/gojq"

	"github.com/stretchr/testify/assert"
)

func jqRun(expr string, body []byte) (v any, err error) {
	// expr = ".foo |.."
	// input = map[string]any{"foo": []any{1, 2, 3}}
	query, err := gojq.Parse(expr)
	if err != nil {
		log.Fatalln(err)
	}
	var input any
	if err = json.Unmarshal(body, &input); err != nil {
		return string(body), err
	}
	iter := query.Run(input) // or query.RunWithContext
	v, ok := iter.Next()
	if !ok {
		return nil, errors.New("no value")
	}
	if err, ok := v.(error); ok {
		if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
			return nil, errors.New("halted")
		}
	}
	return v, nil
}

type AssertError struct {
	Expected any
	Actual any
	Err    error
}

func (e *AssertError) Error() string {
	// es := 	spew.Dump(e.Expect)
	// as := 	spew.Dump(e.Actual)

	return fmt.Sprintf("assert equal failed:\nexpected: %#v,\nactual: %#v,\n err: %v",
		sdump(e.Expected),
		sdump(e.Actual),
		e.Err)
}

func jqEqual(expr string, expectStr string,  body []byte) (err error) {
	var expect any
	if err := json.Unmarshal([]byte(expectStr), &expect); err != nil {
		return fmt.Errorf("json.Unmarshal data error: %s", expectStr)
	}
	actual, err := jqRun(expr, body)
	if err != nil {
		return err
	}
	// assert.EqualValues(t, expect, actual)
	if assert.ObjectsAreEqualValues(expect, actual) {
		return nil
	}
	return &AssertError{Expected: expect, Actual: actual, Err: errors.New("not equal")}
}

func jqObjToString(obj any) string{
	switch obj.(type) {
	case string:
		return obj.(string)
	default:
		buf, _ := json.Marshal(obj)
		return string(buf)
	}
}
func jqEqualRegex(expr string, expectStr string,  body []byte) (err error) {
	actual, err := jqRun(expr, body)
	if err != nil {
		return err
	}
	actualStr := jqObjToString(actual)
	if regexp.MustCompile(expectStr).MatchString(actualStr) {
		return nil
	}

	return &AssertError{Expected: expectStr, Actual: actual, Err: errors.New("not equal")}
}
func jqEqualWildcard(expr string, expect string,  body []byte) (err error) {
	actual, err := jqRun(expr, body)
	if err != nil {
		return err
	}
	actualStr := jqObjToString(actual)
	expect = strings.Replace(expect, "*", ".*", -1)
	if regexp.MustCompile(expect).MatchString(actualStr) {
		return nil
	}
	return &AssertError{Expected: expect, Actual: actual, Err: errors.New("not equal")}
}

// refer: github.com/stretchr/testify/assert.EqualValues
func sdump(expected any) (e string ){
	et := reflect.TypeOf(expected)
	switch et {
	case reflect.TypeOf(""):
		e = reflect.ValueOf(expected).String()
	// case reflect.TypeOf(time.Time{}):
	// 	e = spewConfigStringerEnabled.Sdump(expected)
	// 	a = spewConfigStringerEnabled.Sdump(actual)
	default:
		e = spew.Sdump(expected)
	}
	return e
}
