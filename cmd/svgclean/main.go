// Copyright (c) 2018, The Decred developers
// See LICENSE for details.

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chappjc/svgclean"
)

func main() {
	var inSVG string
	outFile := os.Stdout

	args := os.Args
	switch len(args) {
	case 1:
		fmt.Println("Usage:")
		fmt.Println("\tsvgclean in.svg [out.svg]")
		os.Exit(0)
	case 2:
		// input file, output to stdout
		b, err := ioutil.ReadFile(args[1])
		if err != nil {
			fmt.Printf("Failed to open input file %s: %v\n", args[1], err)
			os.Exit(1)
		}
		inSVG = string(b)
	case 3:
		// input and output files
		b, err := ioutil.ReadFile(args[1])
		if err != nil {
			fmt.Printf("Failed to open input file %s: %v\n", args[1], err)
			os.Exit(1)
		}
		inSVG = string(b)
		outFile, err = os.OpenFile(args[2], os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Printf("Failed to open output file %s: %v\n", args[2], err)
			os.Exit(1)
		}
		defer outFile.Close()
	default:
		os.Exit(1)
	}

	outSVG := svgclean.CleanSVGString(inSVG)

	_, err := outFile.WriteString(outSVG)
	if err != nil {
		fmt.Printf("Failed to write output file %s: %v", outFile.Name(), err)
	}
	if outFile == os.Stdout {
		fmt.Fprintf(outFile, "\n")
	}
}
