package filter

import (
	"github.com/gato/aof"
	"regexp"
	"testing"
)

func TestCommandMatch(t *testing.T) {
	var op aof.Operation
	var ftr Filter

	op.Command = "SELECT"

	// simple match (exact)
	ftr.Command = regexp.MustCompile("SELECT")
	if !Match(op, ftr, false) {
		t.Errorf("Op.Command '%s' should match '%s'", op.Command, ftr.Command)
		return
	}
	// inverse match
	if Match(op, ftr, true) {
		t.Errorf("inverse of a matching filter should return false")
		return
	}

	// Regexp matches
	ftr.Command = regexp.MustCompile("SEL.*")
	if !Match(op, ftr, false) {
		t.Errorf("Op.Command '%s' should match '%s'", op.Command, ftr.Command)
		return
	}
	ftr.Command = regexp.MustCompile(".ELECT")
	if !Match(op, ftr, false) {
		t.Errorf("Op.Command '%s' should match '%s'", op.Command, ftr.Command)
		return
	}

	// no match
	ftr.Command = regexp.MustCompile("-ELECT")
	if Match(op, ftr, false) {
		t.Errorf("Op.Command '%s' shouldn't match '%s'", op.Command, ftr.Command)
		return
	}
	// inverse match
	if !Match(op, ftr, true) {
		t.Errorf("inverse of a non matching filter should return true")
		return
	}

}
