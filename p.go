package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	BUFSIZE  = 1024
	progname = "p"
)

var tty *bufio.Reader

func init() {
	s := "/dev/tty"
	f, e := os.Open(s)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%s: error opening '%s': %s\n", progname, s, e)
		os.Exit(1)
	}
	tty = bufio.NewReaderSize(f, BUFSIZE)
}

func ttyin() {
	line, _, e := tty.ReadLine()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%s: error reading ttyin: %s\n", e)
		os.Exit(1)
	}
	if line == nil || (len(line) > 0 && line[0] == 'q') {
		os.Exit(0)
	}
}

func print(f *os.File, filename string, pagesize int) {
	r := bufio.NewReaderSize(f, BUFSIZE)
	w := bufio.NewWriter(os.Stdout)
	nlines := 0
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				w.Flush()
				return
			}
			fmt.Fprintf(os.Stderr, "%s: error reading %s: %s\n", progname, filename, err)
			break
		} else {
			line = strings.TrimRight(line, "\r\n")
			w.Write([]byte(line))
			w.Flush()
			nlines += 1
			if nlines >= pagesize {
				ttyin()
				nlines = 0
			} else {
				w.WriteRune('\n')
			}
		}
	}
}

func main() {
	var n *int = flag.Int("n", 22, "number of lines to print")
	flag.Parse()
	if flag.NArg() == 0 {
		print(os.Stdin, "stdin", *n)
	}
	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: couldn't open %s:  %s\n", progname, flag.Arg(i), err)
			os.Exit(1)
		}
		print(f, flag.Arg(i), *n)
		f.Close()
	}
}
