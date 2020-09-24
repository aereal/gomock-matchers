package matchers

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_mapMatcher_Matches(t *testing.T) {
	type m map[string]interface{}
	type v interface{}
	cases := []struct {
		name    string
		matcher gomock.Matcher
		yes, no []v
	}{
		{"nil", Map(nil), []v{m{}}, []v{0, m{"a": 1}}},
		{"exactly", Map(nil).Add("a", gomock.Eq(1)), []v{m{"a": 1}}, []v{m{"b": 2}}},
		{"subset", Map(nil).Add("a", gomock.Eq(1)), []v{m{"a": 1, "b": 2}}, nil},
		{"superset", Map(nil).Add("a", gomock.Eq(1)).Add("b", gomock.Eq(2)), []v{m{"a": 1, "b": 2}}, []v{m{"a": 1}}},
		{"nested", Map(nil).Add("parent", Map(nil).Add("child", Map(nil).Add("a", gomock.Eq(1)))), []v{m{"parent": m{"child": m{"a": 1}}}}, []v{m{"b": 2}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, x := range c.yes {
				if !c.matcher.Matches(x) {
					t.Errorf("%v %s: expected matched but not", x, c.matcher)
				}
			}
			for _, x := range c.no {
				if c.matcher.Matches(x) {
					t.Errorf("%v %s: expected NOT matched but matched", x, c.matcher)
				}
			}
		})
	}
}
