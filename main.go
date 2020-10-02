// Copyright (c) 2020, Microsoft Corporation, Sean Hinchee
// Licensed under the MIT License.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/seh-msft/burpxml"
)

var (
	useStdin    = flag.Bool("s", false, "read from stdin (rather than first argument)")
	toJson      = flag.Bool("j", false, "emit XML as JSON only")
	toCsv       = flag.Bool("c", false, "emit XML as CSV only")
	toGo        = flag.Bool("g", false, "emit XML as valid Go syntax only")
	outFileName = flag.String("o", "", "output file name (rather than stdout)")
	inFileName  = flag.String("i", "", "input file name (rather than first argument)")
	noReq       = flag.Bool("r", false, "omit requests in CSV (as they may corrupt output in excel)")
	noResp      = flag.Bool("R", false, "omit responses in CSV (as they may corrupt output in excel)")
	decode      = flag.Bool("d", false, "decode base64 bodies (may corrupt output)")
)

// Parse burp proxy history XML output
// You might have to strip the burp XML preamble and version string
func main() {
	flag.Parse()
	args := flag.Args()

	/* Handle argument combinatorics */

	if len(args) < 1 && !(*useStdin || *inFileName != "") {
		fatal("err: specify an xml file to parse")
	}

	if !xor(*useStdin, *inFileName != "", len(args) > 0) {
		fatal("err: specify input as stdin ⊻ infile ⊻ argument")
	}

	if (*toCsv || *toJson || *toGo) && !xor(*toCsv, *toJson, *toGo) {
		fatal("err: specify output as CSV ⊻ JSON ⊻ Go")
	}

	var inFile string // Input file name
	var f io.Reader   // Input file struct

	var of io.Writer = os.Stdout
	if *outFileName != "" && *outFileName != "-" {
		file, err := os.Create(*outFileName)
		if err != nil {
			fatal("err: could not open output file ⇒ ", err)
		}
		of = file
		defer file.Close()
	}

	if len(args) > 0 && args[0] == "-" {
		f = os.Stdin
	}

	if *useStdin {
		f = os.Stdin
	} else if *inFileName != "" {
		inFile = *inFileName
	} else {
		inFile = args[0]
	}

	// Open XML file
	if f == nil {
		file, err := os.Open(inFile)
		if err != nil {
			fatal("err: could not open input file ⇒ ", err)
		}
		f = file
		defer file.Close()
	}

	// Use bufio
	bf := bufio.NewReader(f)
	bof := bufio.NewWriter(of)
	defer bof.Flush()

	/* Decoding */

	items, err := burpxml.Parse(bf, *decode)
	if err != nil {
		fatal("err: parsing failed ⇒ ", err)
	}

	if *toJson {
		err := items.Json(bof)
		if err != nil {
			fatal("err: json conversion failed ⇒ ", err)
		}
		return
	}

	if *toCsv {
		_, err := items.Csv(bof, *noReq, *noResp)
		if err != nil {
			fatal("err: csv conversion failed ⇒ ", err)
		}
		return
	}

	if *toGo {
		bof.WriteString(items.Go())
		return
	}

	// Emit items structure
	bof.WriteString(fmt.Sprintln(items))
}

/* Utility routines */

// Fatal - end program with an error message and newline
func fatal(s ...interface{}) {
	fmt.Fprintln(os.Stderr, s...)
	os.Exit(1)
}

// One and only one is true
func xor(in ...bool) bool {
	f := false

	for _, b := range in {
		switch {
		case b && f:
			// Two are true
			return false
		case b && !f:
			f = true
		}
	}

	return f
}
