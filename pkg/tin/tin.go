package tin

import (
	"fmt"
	"io/ioutil"
)

func CompileFile(filePath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	tokens := tokenizeSource(string(source))
	fmt.Println(tokens)

	program := parseProgramFromTokens(tokens)
	fmt.Println(program)

	asm := generateNasmAppleSilicon(program)
	fmt.Println(asm)

	return err
}
