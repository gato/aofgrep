package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gato/aof"
	"io"
	"os"
	"regexp"
)

func processFile(input *bufio.Reader, filter Filter, invert bool) (matched, processed int) {
	processed = 0
	matched = 0
	for {
		processed++
		op, err := aof.ReadOperation(input)
		if err != nil {
			if err == io.EOF {
				return
			}
			info := fmt.Sprintf("Error processing command %d Error:%s\n", processed, err.Error())
			os.Stderr.WriteString(info)
			os.Exit(2)
		}
		if match(op, filter, invert) {
			err = op.ToAof(os.Stdout)
			matched++
			if err != nil {
				info := fmt.Sprintf("Error writing command %d Error:%s\n", processed, err.Error())
				os.Stderr.WriteString(info)
				os.Exit(3)
			}
		}
	}
}

type Filter struct {
	Command   *regexp.Regexp
	SubOp     *regexp.Regexp
	Key       *regexp.Regexp
	Parameter *regexp.Regexp
}

func match(op aof.Operation, filter Filter, inverse bool) bool {
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
	return true
}

func main() {

	var matched, processed int
	var filter Filter = Filter{}

	filterCommand := flag.String("command", "", "a regexp for filtering by command")
	filterSubop := flag.String("subop", "", "a regexp for filtering by sub operation keys")
	filterKey := flag.String("key", "", "a regexp for filtering by key")
	filterParameter := flag.String("param", "", "a regexp for filtering by parameter")
	debug := flag.Bool("d", false, "output debug information (to STDERR)")
	invert := flag.Bool("v", false, "output command if does not match")

	flag.Parse()

	var err error
	if *filterCommand != "" {
		filter.Command, err = regexp.Compile(*filterCommand)
		if err != nil {
			info := fmt.Sprintf("Can't compile command regexp:%s Error:%s\n", *filterCommand, err.Error())
			os.Stderr.WriteString(info)
			os.Exit(1)
		}
	}

	if *filterSubop != "" {
		filter.SubOp, err = regexp.Compile(*filterSubop)
		if err != nil {
			info := fmt.Sprintf("Can't compile subop regexp:%s Error:%s\n", *filterSubop, err.Error())
			os.Stderr.WriteString(info)
			os.Exit(1)
		}
	}
	if *filterKey != "" {
		filter.Key, err = regexp.Compile(*filterKey)
		if err != nil {
			info := fmt.Sprintf("Can't compile key regexp:%s Error:%s\n", *filterKey, err.Error())
			os.Stderr.WriteString(info)
			os.Exit(1)
		}
	}
	if *filterParameter != "" {
		filter.Parameter, err = regexp.Compile(*filterParameter)
		if err != nil {
			info := fmt.Sprintf("Can't compile parameter regexp:%s Error:%s\n", *filterParameter, err.Error())
			os.Stderr.WriteString(info)
			os.Exit(1)
		}
	}

	if len(flag.Args()) > 0 {
		for _, file := range flag.Args() {
			if *debug {
				info := fmt.Sprintf("Parsing file %s\n", file)
				os.Stderr.WriteString(info)
			}
			f, err := os.Open(file)
			if err != nil {
				info := fmt.Sprintf("Can't open file:%s Error:%s\n", file, err.Error())
				os.Stderr.WriteString(info)
				os.Exit(1)
			}
			defer f.Close()
			m, p := processFile(bufio.NewReader(f), filter, *invert)
			matched += m
			processed += p
		}
	} else {
		// process stdin
		matched, processed = processFile(bufio.NewReader(os.Stdin), filter, *invert)
	}
	if *debug {
		info := fmt.Sprintf("%d matches found %d commands processed\n", matched, processed)
		os.Stderr.WriteString(info)
	}
}
