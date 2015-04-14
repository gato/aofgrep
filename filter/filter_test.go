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
	// insensitive case test
	op.Command = "select"
	ftr.Command = regexp.MustCompile("SELECT")
	if !Match(op, ftr, false) {
		t.Errorf("Op.Command '%s' should match '%s'", op.Command, ftr.Command)
		return
	}
	op.Command = "Select"
	if !Match(op, ftr, false) {
		t.Errorf("Op.Command '%s' should match '%s'", op.Command, ftr.Command)
		return
	}
	op.Command = "SeLeCt"
	if !Match(op, ftr, false) {
		t.Errorf("Op.Command '%s' should match '%s'", op.Command, ftr.Command)
		return
	}

}

func TestNoFilterMatch(t *testing.T) {
	var op aof.Operation
	var ftr Filter
	op.Command = "SELECT"

	// empty filter match
	if !Match(op, ftr, false) {
		t.Errorf("Op.Command '%s' should match '%s'", op.Command, ftr.Command)
		return
	}
	// inverse match => false
	if Match(op, ftr, true) {
		t.Errorf("inverse of a matching filter should return false")
		return
	}
}

func TestKeyMatch(t *testing.T) {
	var op aof.Operation
	var ftr Filter

	op.Key = "K1"

	// simple match (exact)
	ftr.Key = regexp.MustCompile("K1")
	if !Match(op, ftr, false) {
		t.Errorf("Op.Key '%s' should match '%s'", op.Key, ftr.Key)
		return
	}
	// inverse match
	if Match(op, ftr, true) {
		t.Errorf("inverse of a matching filter should return false")
		return
	}

	// Regexp matches
	ftr.Key = regexp.MustCompile("K.*")
	if !Match(op, ftr, false) {
		t.Errorf("Op.Key '%s' should match '%s'", op.Key, ftr.Key)
		return
	}
	ftr.Key = regexp.MustCompile(".1")
	if !Match(op, ftr, false) {
		t.Errorf("Op.Key '%s' should match '%s'", op.Key, ftr.Key)
		return
	}

	// no match
	ftr.Key = regexp.MustCompile("K2")
	if Match(op, ftr, false) {
		t.Errorf("Op.Key '%s' shouldn't match '%s'", op.Key, ftr.Key)
		return
	}
	// inverse match
	if !Match(op, ftr, true) {
		t.Errorf("inverse of a non matching filter should return true")
		return
	}
	// insensitive case test
	ftr.Key = regexp.MustCompile("k1")
	if Match(op, ftr, false) {
		t.Errorf("Op.Key '%s' shouldn't match '%s'", op.Key, ftr.Key)
		return
	}

	op.Key = "k1"
	ftr.Key = regexp.MustCompile("K1")
	if Match(op, ftr, false) {
		t.Errorf("Op.Key '%s' shouldn't match '%s'", op.Key, ftr.Key)
		return
	}
}
