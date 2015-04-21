package main

import (
	"bufio"
	"github.com/gato/aofgrep/filter"
	"regexp"
	"strings"
	"testing"
)

// should move this to a testing library. (its also in aof_test.go)
type RecordWriter []byte

func (this *RecordWriter) Write(b []byte) (int, error) {
	*this = append(*this, b...)
	return len(b), nil
}

func TestProcessInput(t *testing.T) {
	var ftr filter.Filter
	var expected = "*2\r\n$6\r\nSELECT\r\n$1\r\n0\r\n"
	var input = bufio.NewReader(strings.NewReader(expected))
	ftr.Command = regexp.MustCompile("SELECT")
	rec := &RecordWriter{}
	matched, processed, err := processInput(input, rec, ftr, false)
	if err != nil {
		t.Errorf("processInput:'%s'", err.Error())
		return
	}
	if matched != 1 {
		t.Errorf("Invalid match count:'%d' expected:'1'", matched)
		return
	}
	if processed != 1 {
		t.Errorf("Invalid processed count:'%d' expected:'1'", processed)
		return
	}
	if string(*rec) != expected {
		t.Errorf("Invalid output:'%s' expected:'%s'", string(*rec), expected)
		return
	}
}

func TestProcessInputNoMatch(t *testing.T) {
	var ftr filter.Filter
	var expected = ""
	var input = bufio.NewReader(strings.NewReader("*2\r\n$6\r\nSELECT\r\n$1\r\n0\r\n"))
	ftr.Command = regexp.MustCompile("SADD")
	rec := &RecordWriter{}
	matched, processed, err := processInput(input, rec, ftr, false)
	if err != nil {
		t.Errorf("processInput:'%s'", err.Error())
		return
	}
	if matched != 0 {
		t.Errorf("Invalid match count:'%d' expected:'0'", matched)
		return
	}
	if processed != 1 {
		t.Errorf("Invalid processed count:'%d' expected:'1'", processed)
		return
	}
	if string(*rec) != expected {
		t.Errorf("Invalid output:'%s' expected:'%s'", string(*rec), expected)
		return
	}
}
