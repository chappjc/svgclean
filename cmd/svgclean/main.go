// Copyright (c) 2018, The Decred developers
// See LICENSE for details.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chappjc/svgclean"
)

var black = flag.Bool("blacklist", false, "Use element blacklist instead of whitelist.")

func main() {
	flag.Parse()

	var inSVG string
	outFile := os.Stdout

	args := flag.Args()
	fmt.Println(args)
	switch len(args) {
	case 0:
		fmt.Println("Usage:")
		fmt.Println("\tsvgclean in.svg [out.svg]")
		os.Exit(0)
	case 1:
		// input file, output to stdout
		b, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Printf("Failed to open input file %s: %v\n", args[0], err)
			os.Exit(1)
		}
		inSVG = string(b)
	case 2:
		// input and output files
		b, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Printf("Failed to open input file %s: %v\n", args[0], err)
			os.Exit(1)
		}
		inSVG = string(b)
		outFile, err = os.OpenFile(args[1], os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Printf("Failed to open output file %s: %v\n", args[1], err)
			os.Exit(1)
		}
		defer outFile.Close()
	default:
		os.Exit(1)
	}

	var outSVG string
	if *black {
		outSVG = svgclean.CleanSVGStringWhite(inSVG)
	} else {
		outSVG = svgclean.CleanSVGStringWhite(inSVG)
	}

	_, err := outFile.WriteString(outSVG)
	if err != nil {
		fmt.Printf("Failed to write output file %s: %v", outFile.Name(), err)
	}
	if outFile == os.Stdout {
		fmt.Fprintf(outFile, "\n")
	}
}
