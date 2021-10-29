package tin

import (
	"fmt"
	"strings"
)

func generateNasmX8664(program Program) string {
	gen := strings.Builder{}

	// Text section
	gen.WriteString("section .text\n")
	gen.WriteString("global _start\n")

	// Main function
	gen.WriteString("\n")
	gen.WriteString("_start:\n")
	for _, inst := range program {
		generateX8664Instruction(&gen, inst)
	}
	gen.WriteString("  ;; exit syscall\n")
	gen.WriteString("  mov rax, 0x3c\n")
	gen.WriteString("  mov rdi, 0\n")
	gen.WriteString("  syscall\n")

	return gen.String()
}

func generateX8664Instruction(gen *strings.Builder, inst Instruction) {
	switch inst.Kind {
	case InstKindPushInt:
		gen.WriteString("  ;; push int\n")
		gen.WriteString(fmt.Sprintf("  mov rax, %d\n", inst.ValueInt))
		gen.WriteString("  push rax\n")
	case InstKindIntrinsic:
		generateX8664Intrinsic(gen, inst.ValueIntrinsic)
	default:
		panic("generateX8664Instruction: unreachable")
	}
}

func generateX8664Intrinsic(gen *strings.Builder, intrinsic Intrinsic) {
	switch intrinsic {
	case IntrinsicPlus:
		gen.WriteString("  ;; add\n")
		gen.WriteString("  pop rbx\n")
		gen.WriteString("  pop rax\n")
		gen.WriteString("  add rax, rbx\n")
		gen.WriteString("  push rax\n")
	case IntrinsicMinus:
		gen.WriteString("  ;; sub\n")
		gen.WriteString("  pop rbx\n")
		gen.WriteString("  pop rax\n")
		gen.WriteString("  sub rax, rbx\n")
		gen.WriteString("  push rax\n")
	case IntrinsicTimes:
		gen.WriteString("  ;; mul\n")
		gen.WriteString("  pop rbx\n")
		gen.WriteString("  pop rax\n")
		gen.WriteString("  imul rbx\n")
		gen.WriteString("  push rax\n")
	case IntrinsicDivMod:
		gen.WriteString("  ;; divmod\n")
		gen.WriteString("  pop rbx\n")
		gen.WriteString("  pop rax\n")
		gen.WriteString("  idiv rbx\n")
		gen.WriteString("  push rax\n")
		gen.WriteString("  push rdx\n")
	default:
		panic("generateX8664Intrinsic: unreachable")
	}
}
