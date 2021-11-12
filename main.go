package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"eagain.net/go/read-file-to-env/internal"
	"golang.org/x/sys/unix"
)

var prog = filepath.Base(os.Args[0])

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", prog)
	fmt.Fprintf(flag.CommandLine.Output(), "  %s [-one-line=ENV=FILE].. [--] CMD..\n", prog)
	fmt.Fprintf(flag.CommandLine.Output(), "\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	var oneLines internal.OneLineFlag
	flag.Var(&oneLines, "one-line", "for `ENV=FILE`, read ENV from FILE")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		log.Printf("missing command to run")
		os.Exit(2)
	}

	env := os.Environ()
	// If we ever add other types of file reading than one line, make sure to collect everything to a single slice to retain relative order.
	// Something like `[]EnvVar` that contains a `Source` interface.
	for _, e := range oneLines {
		line, err := internal.ReadOneLine(e.Path)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if strings.IndexByte(line, '\x00') >= 0 {
			log.Fatalf("read: %q: environment variable cannot contain NUL", e.Path)
		}
		env = append(env, fmt.Sprintf("%s=%s", e.Name, line))
	}

	args := flag.Args()
	cmd := args[0]
	if !filepath.IsAbs(cmd) {
		c, err := exec.LookPath(cmd)
		if err != nil {
			// LookPath adds an "exec:" prefix etc to error message
			log.Fatalf("%v", err)
		}
		cmd = c
	}
	if err := unix.Exec(cmd, args, env); err != nil {
		log.Fatalf("exec: %q: %v", cmd, err)
	}
}
