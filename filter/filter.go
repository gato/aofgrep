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
	if filter.Command != nil && filter.Command.FindStringIndex(op.Command) == nil && !inverse {
		return false
	}
	if filter.SubOp != nil && filter.SubOp.FindStringIndex(op.SubOp) == nil && !inverse {
		return false
	}
	if filter.Key != nil && filter.Key.FindStringIndex(op.Key) == nil && !inverse {
		return false
	}
	if filter.Parameter != nil {
		for _, p := range op.Arguments {
			if filter.Parameter.FindStringIndex(p) != nil {
				if !inverse {
					return true
				}
			}
		}
		return false
	}
	if inverse {
		return false
	}
	return true
}
