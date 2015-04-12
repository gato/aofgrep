package filter

import (
	"github.com/gato/aof"
	"regexp"
)

type Filter struct {
	Command   *regexp.Regexp
	SubOp     *regexp.Regexp
	Key       *regexp.Regexp
	Parameter *regexp.Regexp
}

func Match(op aof.Operation, filter Filter, inverse bool) bool {
	rCode := false
	if inverse {
		rCode = true
	}
	if filter.Command != nil && filter.Command.FindStringIndex(op.Command) == nil {
		return rCode
	}
	if filter.SubOp != nil && filter.SubOp.FindStringIndex(op.SubOp) == nil {
		return rCode
	}
	if filter.Key != nil && filter.Key.FindStringIndex(op.Key) == nil {
		return rCode
	}
	if filter.Parameter != nil {
		for _, p := range op.Arguments {
			if filter.Parameter.FindStringIndex(p) != nil {
				break
			}
		}
	}
	return !rCode
}
