package tin

import (
	"fmt"
	"io/ioutil"
	"os"
)

type CompilerOption struct {
	InputPath  string
	OutputPath string
}

func CompileFile(option CompilerOption) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	source, err := ioutil.ReadFile(option.InputPath)
	if err != nil {
		panic(err)
	}

	tokens := tokenizeSource(string(source))
	program := parseProgramFromTokens(tokens)
	asm := generateNasmX8664(program)

	if e := ioutil.WriteFile(option.OutputPath, []byte(asm), os.ModePerm); e != nil {
		panic(e)
	}

	return err
}
