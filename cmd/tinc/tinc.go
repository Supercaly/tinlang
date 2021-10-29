package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Supercaly/tinlang/pkg/tin"
)

func usage(stream io.Writer, program string) {
	fmt.Fprintf(stream, "Usage %s [OPTIONS] <input.tin>\n", program)
	fmt.Fprintf(stream, "OPTIONS:\n")
	fmt.Fprintf(stream, "  -h	Print this help message\n")
}

func main() {
	program := os.Args[0]
	os.Args = os.Args[1:]

	if len(os.Args) <= 0 {
		usage(os.Stderr, program)
		log.Fatal("ERROR: missing input file name")
	}

	if os.Args[0] == "-h" {
		usage(os.Stdout, program)
		os.Exit(0)
	}

	inputFilePath := os.Args[0]
	outputFilePathNoExt := strings.TrimSuffix(inputFilePath, filepath.Ext(inputFilePath))

	option := tin.CompilerOption{
		InputPath:  inputFilePath,
		OutputPath: outputFilePathNoExt + ".asm",
	}

	if err := tin.CompileFile(option); err != nil {
		log.Fatal(err)
	}

	nasmCmd := exec.Command("nasm", "-felf64", outputFilePathNoExt+".asm")
	if err := nasmCmd.Run(); err != nil {
		log.Fatal(err)
	}
	ldCmd := exec.Command("ld", "-o", outputFilePathNoExt, outputFilePathNoExt+".o")
	if err := ldCmd.Run(); err != nil {
		log.Fatal(err)
	}
}
