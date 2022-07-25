package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if err := Copy(filepath.Join(dir, from), filepath.Join(dir, to), offset, limit); err != nil {
		log.Fatal(err)
	}
}
