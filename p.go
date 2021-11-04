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

func print(f *os.File, pagesize int, t *bufio.Reader) error {
	r := bufio.NewReaderSize(f, 1024)
	w := bufio.NewWriter(os.Stdout)
	nlines := 0
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				w.Flush()
				return nil
			}
			return err
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

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
	os.Exit(1)
}

func main() {
	progname := os.Args[0]

	n := flag.Int("n", 22, "number of lines to print")
	flag.Parse()

	s := "/dev/tty"
	f, err := os.Open(s)
	if err != nil {
		exit(fmt.Errorf("%s: %s\n", progname, err))
	}
	tty := bufio.NewReaderSize(f, 1024)

	if flag.NArg() == 0 {
		if err := print(os.Stdin, *n, tty); err != nil {
			exit(fmt.Errorf("%s: stdin:, %s\n", progname, err))
		}
		os.Exit(0)
	}

	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if err != nil {
			exit(fmt.Errorf("%s: %s\n", progname, err))
		}
		if err := print(f, *n, tty); err != nil {
			exit(fmt.Errorf("%s: %s: %s\n", progname, flag.Arg(i), err))
		}
		f.Close()
	}
}
