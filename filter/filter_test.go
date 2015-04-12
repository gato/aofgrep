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
	ftr.Command = regexp.MustCompile("SELECT")
	if Match(op, ftr, false) == false {
		t.Errorf("Op.Command '%s' should match '%s'", op.Command, ftr.Command)
		return
	}
	if Match(op, ftr, true) == true {
		t.Errorf("inverse of a matching filter should return false")
		return
	}

}
