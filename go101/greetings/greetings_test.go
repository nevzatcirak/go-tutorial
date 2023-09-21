package greetings

import (
	"regexp"
	"testing"
)

func TestName(t *testing.T) {
	name := "Gladys"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Name(name)
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Name("%v") = %q, %v, want match for %#q, nil`, name, msg, err, want)
	}
}
