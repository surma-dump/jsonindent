package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	outputFlag = flag.String("o", "-", "Output file")
	helpFlag   = flag.Bool("h", false, "Show this help")
	output     io.Writer
)

func main() {
	flag.Parse()
	if *helpFlag {
		fmt.Printf("Usage: jsonindent [options] [input]\n")
		flag.PrintDefaults()
		return
	}
	log.SetPrefix("")

	if *outputFlag == "-" {
		output = os.Stdout
	} else {
		f, e := os.Create(*outputFlag)
		if e != nil {
			log.Fatalf("Could not open %s: %s\n", *outputFlag, e)
		}
		defer f.Close()
		output = f
	}

	if flag.NArg() == 0 {
		jsonindent(os.Stdin)
		return
	}
	for _, file := range flag.Args() {
		func() {
			f, e := os.Open(file)
			if e != nil {
				log.Printf("Could not open %s: %s\n", file, e)
				return
			}
			defer f.Close()
			log.SetPrefix(file+": ")
			jsonindent(f)
			log.SetPrefix("")
		}()
	}
}

func jsonindent(r io.Reader) {
	d := json.NewDecoder(r)
	var v interface{}
	e := d.Decode(&v)
	if e != nil {
		log.Printf("Could not unmarshall: %s\n", e)
		return
	}

	b, e := json.MarshalIndent(v, "", "\t")
	if e != nil {
		log.Printf("Could not marshall: %s\n", e)
		return
	}
	output.Write(b)
}
