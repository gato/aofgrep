package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gato/aof"
	"github.com/gato/aofgrep/filter"
	"io"
	"os"
	"regexp"
)

func processInput(input *bufio.Reader, ftr filter.Filter, invert bool) (matched, processed int) {
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
		if filter.Match(op, ftr, invert) {
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

func main() {

	var matched, processed int
	var ftr filter.Filter = filter.Filter{}

	filterCommand := flag.String("command", "", "a regexp for filtering by command")
	filterSubop := flag.String("subop", "", "a regexp for filtering by sub operation keys")
	filterKey := flag.String("key", "", "a regexp for filtering by key")
	filterParameter := flag.String("param", "", "a regexp for filtering by parameter")
	debug := flag.Bool("d", false, "output debug information (to STDERR)")
	invert := flag.Bool("v", false, "output command if does not match")

	flag.Parse()

	var err error
	if *filterCommand != "" {
		ftr.Command, err = regexp.Compile(*filterCommand)
		if err != nil {
			info := fmt.Sprintf("Can't compile command regexp:%s Error:%s\n", *filterCommand, err.Error())
			os.Stderr.WriteString(info)
			os.Exit(1)
		}
	}

	if *filterSubop != "" {
		ftr.SubOp, err = regexp.Compile(*filterSubop)
		if err != nil {
			info := fmt.Sprintf("Can't compile subop regexp:%s Error:%s\n", *filterSubop, err.Error())
			os.Stderr.WriteString(info)
			os.Exit(1)
		}
	}
	if *filterKey != "" {
		ftr.Key, err = regexp.Compile(*filterKey)
		if err != nil {
			info := fmt.Sprintf("Can't compile key regexp:%s Error:%s\n", *filterKey, err.Error())
			os.Stderr.WriteString(info)
			os.Exit(1)
		}
	}
	if *filterParameter != "" {
		ftr.Parameter, err = regexp.Compile(*filterParameter)
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
			m, p := processInput(bufio.NewReader(f), ftr, *invert)
			matched += m
			processed += p
		}
	} else {
		// process stdin
		matched, processed = processInput(bufio.NewReader(os.Stdin), ftr, *invert)
	}
	if *debug {
		info := fmt.Sprintf("%d matches found %d commands processed\n", matched, processed)
		os.Stderr.WriteString(info)
	}
}
