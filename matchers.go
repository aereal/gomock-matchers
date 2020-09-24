package matchers

import (
	"fmt"
	"reflect"

	"github.com/golang/mock/gomock"
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

var (
	zeroValue = reflect.Value{}
)
