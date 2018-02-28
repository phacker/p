package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func tty(t *bufio.Reader) {
	line, _, e := t.ReadLine()
	if e != nil {
		log.Fatal(e)
	}
	if len(line) > 0 && line[0] == 'q' {
		os.Exit(0)
	}
}

func print(f *os.File, filename string, pagesize int, t *bufio.Reader, progname string) {
	r := bufio.NewReaderSize(f, 1024)
	w := bufio.NewWriter(os.Stdout)
	nlines := 0
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				w.Flush()
				return
			}
			log.Printf("%s: error reading %s: %s\n", progname, filename, err)
			break
		} else {
			line = strings.TrimRight(line, "\r\n")
			w.Write([]byte(line))
			w.Flush()
			nlines++

			if nlines >= pagesize {
				tty(t)
				nlines = 0
			} else {
				w.WriteRune('\n')
			}
		}
	}
}

func main() {
	flag.Parse()
	progname := flag.Args()[0]

	n := flag.Int("n", 22, "number of lines to print")
	flag.Parse()

	s := "/dev/tty"
	f, e := os.Open(s)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%s: error opening '%s': %s\n", progname, s, e)
		os.Exit(1)
	}
	tty := bufio.NewReaderSize(f, 1024)

	if flag.NArg() == 0 {
		print(os.Stdin, "stdin", *n, tty, progname)
	}

	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: couldn't open %s:  %s\n", progname, flag.Arg(i), err)
			os.Exit(1)
		}
		print(f, flag.Arg(i), *n, tty, progname)
		f.Close()
	}
}
