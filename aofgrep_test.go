package main

import (
	"fmt"
	"github.com/gato/aof"
	"github.com/gato/aofgrep/filter"
	"regexp"
	"strings"
	"testing"
)

// should move this to a testing library. (its also in aof_test.go)
type RecordWriter []byte

func (r *RecordWriter) Write(b []byte) (int, error) {
	*r = append(*r, b...)
	return len(b), nil
}

type ErrorNWriter struct {
	current int
	failing int
}

func (r *ErrorNWriter) Write(b []byte) (int, error) {
	r.current++
	if r.current == r.failing {
		return len(b), fmt.Errorf("Some error")
	}
	return len(b), nil
}

func newErrorNWriter(failing int) ErrorNWriter {
	return ErrorNWriter{current: 0, failing: failing}
}

func TestProcessInput(t *testing.T) {
	var ftr filter.Filter
	var expected = "*2\r\n$6\r\nSELECT\r\n$1\r\n0\r\n"
	var input = aof.NewBufioReader(strings.NewReader(expected))
	ftr.Command = regexp.MustCompile("SELECT")
	rec := RecordWriter{}
	matched, processed, err := processInput(input, &rec, ftr, false)
	if err != nil {
		t.Errorf("processInput returned error:'%s'", err.Error())
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
	if string(rec) != expected {
		t.Errorf("Invalid output:'%s' expected:'%s'", string(rec), expected)
		return
	}
}

func TestProcessInputNoMatch(t *testing.T) {
	var ftr filter.Filter
	var expected = ""
	var input = aof.NewBufioReader(strings.NewReader("*2\r\n$6\r\nSELECT\r\n$1\r\n0\r\n"))
	ftr.Command = regexp.MustCompile("SADD")
	rec := RecordWriter{}
	matched, processed, err := processInput(input, &rec, ftr, false)
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
	if string(rec) != expected {
		t.Errorf("Invalid output:'%s' expected:'%s'", string(rec), expected)
		return
	}
}

func TestProcessInputEofError(t *testing.T) {
	var ftr filter.Filter
	var input = aof.NewBufioReader(strings.NewReader("*2\r\n$6\r\nSELECT\r\n$1\r\n"))
	var expected = "Error processing command 0 Error:"
	ftr.Command = regexp.MustCompile("SELECT")
	rec := RecordWriter{}
	_, _, err := processInput(input, &rec, ftr, false)
	if err == nil {
		t.Errorf("Error was expected got nil")
		return
	}
	got := err.Error()[0:len(expected)]
	if got != expected {
		t.Errorf("Invalid error:'%s' expected:'%s'", got, expected)
		return
	}
}

func TestProcessInputErrorOnWrite(t *testing.T) {
	var ftr filter.Filter
	var input = aof.NewBufioReader(strings.NewReader("*2\r\n$6\r\nSELECT\r\n$1\r\n0\r\n"))
	var expected = "Error writing command 1 Error:Some error\n"
	ftr.Command = regexp.MustCompile("SELECT")
	rec := newErrorNWriter(1)
	_, _, err := processInput(input, &rec, ftr, false)
	if err == nil {
		t.Errorf("Error was expected got nil")
		return
	}
	if err.Error() != expected {
		t.Errorf("Invalid error:'%s' expected:'%s'", err.Error(), expected)
		return
	}
}

func TestProcessFiles(t *testing.T) {
	var opt options
	opt.Debug = false
	opt.Filter.Command = regexp.MustCompile("SET")
	opt.Files = []string{"test-data-bitop.aof", "test-data.aof"}
	rec := RecordWriter{}
	expected := "*3\r\n$3\r\nset\r\n$2\r\nk4\r\n$4\r\nyolo\r\n*3\r\n$3\r\nSET\r\n$6\r\nlast:1\r\n$32\r\nea4b640e462d11e2a119005056a3cdd9\r\n"
	matched, processed, err := processFiles(&rec, opt)
	if err != nil {
		t.Errorf("processFiles returned error:'%s'", err.Error())
		return
	}
	if matched != 2 {
		t.Errorf("Invalid match count:'%d' expected:'1'", matched)
		return
	}
	if processed != 6 {
		t.Errorf("Invalid processed count:'%d' expected:'6'", processed)
		return
	}
	if string(rec) != expected {
		t.Errorf("Invalid output:'%s' expected:'%s'", string(rec), expected)
		return
	}
}

func TestProcessFilesNotFound(t *testing.T) {
	var opt options
	opt.Debug = false
	opt.Filter.Command = regexp.MustCompile("SET")
	opt.Files = []string{"test-data-bitop.aof", "not-existing.aof"}
	rec := RecordWriter{}
	_, _, err := processFiles(&rec, opt)
	if err == nil {
		t.Errorf("processFiles should return and error")
		return
	}
}
func TestProcessFilesErrorOnWrite(t *testing.T) {
	var opt options
	opt.Debug = false
	opt.Filter.Command = regexp.MustCompile("SET")
	opt.Files = []string{"test-data-bitop.aof", "test-data.aof"}
	rec := newErrorNWriter(1)
	_, _, err := processFiles(&rec, opt)
	if err == nil {
		t.Errorf("processFiles should return and error")
		return
	}
}
