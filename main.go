package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var usage = `
Usage: %s OPTION... [FILE]...
Print selected parts of lines from each FILE to standard output.

With no FILE, or when FILE is -, read standard input.

    -f      --fields FIELDS                 select only these fields; also print any line
                                                that contains no delimiter character.
                                                Fields can be used more than once
    -cw     --cut-on-whitespace             cut on any whitespace but also trim following
                                                whitespace aswell until first non-whitepsace
    -cmw    --cut-on-multi-whitespace       cut on consecutive whitespace but also trim following
                                                whitespace aswell until first non-whitepsace
    -cs     --cut-on-seperator STR          cut whenever STR is encountered in a line
    -cf     --cut-on-format DELIMS          cut line applying one delimiter in STR after
                                                another
    -fsep   --format-seperator STR          use STR as seperator in the DELIMS specification
                                                Default: ','
    -osep   --ouput-seperator STR           use the STR as the output field seperator
                                                Default: ' '

Only one of the following can be used at a time:
    -cw
    -cmw
    -cs
    -cf

FIELDS is made up of one range, or many ranges separated by commas.
Selected input is written in the same order that it is read.
Each range is one of:

    N     N'th field, counted from 1
   -N     (num items on line) - N
    N:    from N'th field, to end of line
   -N:    (num items on line) - N, to the end of line
    N:M   from N'th to M'th (included) field
    :M    from first to M'th (included) field
    :-M   from first to ((num items on line) - M)'th (included) field
    :     from beginning to end of line

DELIMS is made up on one seperator, or many seperators seperated by commas.
The seperator of DELIMS can be changed using --format-seperator.
Each delimiter can be one of:
    s                   cut on next whitespace
    t                   cut on next tab
    w                   cut on next whitespace but consume all consecutive aswell
    m                   cut on next multi whitespace and consume all consecutive aswell
    <str> cut on next encounter of str, where str can by any string

Note: whitespace always means tab and space.

Examples:
    $ echo "A B C" | gut -cw -f 2:
    B C

    $ echo "A  B  C" | gut -f -2:
    B C

    $ echo "A;;B     C" | gut -cf "<;;>,m"
    A B C

    $ echo "A  B,C" | gut -cf "s|<,>" -fsep "|"
    A  B C
`

// selection options
var fieldsArg = arg[string]{aliases: []string{"f", "fields"}}

// cutting options
var cutOnWhitespaceArg = arg[bool]{aliases: []string{"cw", "cut-on-whitespace"}}
var cutOnMultiWhitespaceArg = arg[bool]{aliases: []string{"cmw", "cut-on-multi-whitespace"}}
var cutOnSeperatorArg = arg[string]{aliases: []string{"cs", "cut-on-seperator"}}
var cutOnFormatArg = arg[string]{aliases: []string{"cf", "cut-on-format"}}

// seperators
var formatSeperatorArg = arg[string]{aliases: []string{"fsep", "format-seperator"}, defaultValue: ","}
var outputSeperatorArg = arg[string]{aliases: []string{"osep", "output-seperator"}, defaultValue: " "}

func setupFlags() {
	sArgs := []*arg[string]{&fieldsArg, &cutOnSeperatorArg, &cutOnFormatArg, &formatSeperatorArg, &outputSeperatorArg}
	bArgs := []*arg[bool]{&cutOnWhitespaceArg, &cutOnMultiWhitespaceArg}

	for _, sArg := range sArgs {
		for _, alias := range sArg.aliases {
			flag.StringVar(&sArg.value, alias, sArg.defaultValue, "")
		}
	}

	for _, bArg := range bArgs {
		for _, alias := range bArg.aliases {
			flag.BoolVar(&bArg.value, alias, bArg.defaultValue, "")
		}
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}

	flag.Parse()
}

func getGutter() StringSplitter {
	// ensure that only one action is performed
	actionCount := countValue(
		true,
		cutOnWhitespaceArg.value,
		cutOnMultiWhitespaceArg.value,
		len(cutOnSeperatorArg.value) != 0,
		len(cutOnFormatArg.value) != 0)

	if actionCount > 1 {
		die("You can only use one of the cutting actions but you specified more than one")
	}

	switch {
	case cutOnWhitespaceArg.value:
		return cutterToSplitter(singleWsCutter)
	case cutOnMultiWhitespaceArg.value:
		return cutterToSplitter(multiWsCutter)
	case len(cutOnSeperatorArg.value) != 0:
		return cutterToSplitter(cutterFromSeperator(cutOnSeperatorArg.value))
	case len(cutOnFormatArg.value) != 0:
		cutters, err := flagToCutters(cutOnFormatArg.value, formatSeperatorArg.value)
		if err != nil {
			die("Error: %v", err)
		}
		return func(s string) []string {
			return cutWithCutters(s, cutters)
		}
	}

	return cutterToSplitter(multiWsCutter)
}

func getSpans() []span {
	// get the spans either by default or user provided value
	if len(fieldsArg.value) == 0 {
		return []span{{}}
	}

	spans, err := flagToSpans(fieldsArg.value, ",")
	if err != nil {
		die("Error: %v\n", err)
	}

	return spans
}

func getReaders() []io.Reader {
	files := flag.Args()

	if len(files) == 0 || (len(files) == 1 && files[0] == "-") {
		return []io.Reader{os.Stdin}
	}

	var readers []io.Reader
	for _, fileName := range files {
		if file, err := os.Open(filepath.Clean(fileName)); err != nil {
			die("The file '%s' could not be opened for reading!", fileName)
		} else {
			readers = append(readers, AutoCloseReader{file})
		}
	}
	return readers
}

func do(writer io.Writer, readers []io.Reader, oSep string, spans []span, chunker StringSplitter) {
	for _, reader := range readers {
		lineScanner := bufio.NewScanner(reader)
		for lineScanner.Scan() {
			isFirstPart := true // one can not know in advance what the last one will be
			parts := chunker(lineScanner.Text())
			for _, span := range spans {
				for _, selected := range access(span, parts) {
					if !isFirstPart {
						io.WriteString(writer, oSep)
					}
					io.WriteString(writer, selected)
					isFirstPart = false
				}
			}

			io.WriteString(writer, "\n")
		}

		if lineScanner.Err() != nil {
			die("An error occured during reading: %v", lineScanner.Err())
		}
	}
}

func main() {
	setupFlags()

	gutter := getGutter()
	spans := getSpans()
	readers := getReaders()

	do(os.Stdout, readers, outputSeperatorArg.value, spans, gutter)
}
