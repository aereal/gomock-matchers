package matchers

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/mock/gomock"
)

var (
	ErrNotStruct = errors.New("type is not struct")
)

type any = interface{}

// Map returns a gomock.Matcher that compares to other map.
func Map(im map[any]gomock.Matcher) *mapMatcher {
	if im == nil {
		im = map[any]gomock.Matcher{}
	}
	return &mapMatcher{ms: im}
}

type mapMatcher struct {
	ms map[any]gomock.Matcher
}

func (m *mapMatcher) Add(key any, expect gomock.Matcher) *mapMatcher {
	m.ms[key] = expect
	return m
}

func (m *mapMatcher) Matches(x interface{}) bool {
	t := reflect.ValueOf(x)
	if t.Kind() != reflect.Map {
		return false
	}
	if len(m.ms) == 0 && len(t.MapKeys()) > 0 {
		return false
	}
	for k, vm := range m.ms {
		kt := reflect.ValueOf(k)
		key := t.MapIndex(kt)
		if key == zeroValue { // key not fonud
			return false
		}
		if !vm.Matches(key.Interface()) {
			return false
		}
	}
	return true
}

func (m *mapMatcher) String() string {
	return fmt.Sprintf("%#v", m.ms)
}

// MustStruct returns a new struct matcher.
// If expected is not a struct type, will be panic.
func MustStruct(expected interface{}) *structMatcher {
	m, err := Struct(expected)
	if err != nil {
		panic(err)
	}
	return m
}

// Struct returns a gomock.Matcher that compares to other struct fields.
func Struct(expected interface{}) (*structMatcher, error) {
	v := reflect.ValueOf(expected)
	if v.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	m := &structMatcher{expected: expected, expectedValue: v, fields: map[string]gomock.Matcher{}}
	max := v.NumField()
	vt := v.Type()
	for i := 0; i < max; i++ {
		fv := v.Field(i)
		if !fv.CanInterface() {
			continue
		}
		ft := vt.Field(i)
		m.Field(ft.Name, gomock.Eq(fv.Interface()))
	}
	return m, nil
}

type structMatcher struct {
	expected      interface{}
	expectedValue reflect.Value
	fields        map[string]gomock.Matcher
}

func (m *structMatcher) Field(name string, value gomock.Matcher) *structMatcher {
	m.fields[name] = value
	return m
}

func (m *structMatcher) Matches(x interface{}) bool {
	ev := reflect.ValueOf(x)
	if ev.Kind() != reflect.Struct {
		return false
	}
	if len(m.fields) == 0 && ev.NumField() > 0 {
		return false
	}
	for name, fm := range m.fields {
		f := ev.FieldByName(name)
		if f == zeroValue { // field not found
			return false
		}
		if !fm.Matches(f.Interface()) {
			return false
		}
	}
	return true
}

func (m *structMatcher) String() string {
	tn := m.expectedValue.Type().Name()
	fs := []string{}
	for k, v := range m.fields {
		fs = append(fs, fmt.Sprintf("%s=%s", k, v.String()))
	}
	return fmt.Sprintf("Struct:%s(%s)", tn, strings.Join(fs, ", "))
}

var (
	zeroValue = reflect.Value{}
)
