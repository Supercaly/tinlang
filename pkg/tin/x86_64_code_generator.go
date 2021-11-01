package tin

import (
	"fmt"
	"strings"
)

type x86_64Generator struct {
	text    strings.Builder
	strings []string
}

func generateNasmX8664(program Program) string {
	gen := x86_64Generator{}

	// Text section
	gen.text.WriteString("section .text\n")
	gen.text.WriteString("global _start\n")

	// Main function
	gen.text.WriteString("\n")
	gen.text.WriteString("print:\n")
	gen.text.WriteString("  mov     r9, -3689348814741910323\n")
	gen.text.WriteString("  sub     rsp, 40\n")
	gen.text.WriteString("  mov     BYTE [rsp+31], 10\n")
	gen.text.WriteString("  lea     rcx, [rsp+30]\n")
	gen.text.WriteString(".L2:\n")
	gen.text.WriteString("  mov     rax, rdi\n")
	gen.text.WriteString("  lea     r8, [rsp+32]\n")
	gen.text.WriteString("  mul     r9\n")
	gen.text.WriteString("  mov     rax, rdi\n")
	gen.text.WriteString("  sub     r8, rcx\n")
	gen.text.WriteString("  shr     rdx, 3\n")
	gen.text.WriteString("  lea     rsi, [rdx+rdx*4]\n")
	gen.text.WriteString("  add     rsi, rsi\n")
	gen.text.WriteString("  sub     rax, rsi\n")
	gen.text.WriteString("  add     eax, 48\n")
	gen.text.WriteString("  mov     BYTE [rcx], al\n")
	gen.text.WriteString("  mov     rax, rdi\n")
	gen.text.WriteString("  mov     rdi, rdx\n")
	gen.text.WriteString("  mov     rdx, rcx\n")
	gen.text.WriteString("  sub     rcx, 1\n")
	gen.text.WriteString("  cmp     rax, 9\n")
	gen.text.WriteString("  ja      .L2\n")
	gen.text.WriteString("  lea     rax, [rsp+32]\n")
	gen.text.WriteString("  mov     edi, 1\n")
	gen.text.WriteString("  sub     rdx, rax\n")
	gen.text.WriteString("  xor     eax, eax\n")
	gen.text.WriteString("  lea     rsi, [rsp+32+rdx]\n")
	gen.text.WriteString("  mov     rdx, r8\n")
	gen.text.WriteString("  mov     rax, 1\n")
	gen.text.WriteString("  syscall\n")
	gen.text.WriteString("  add     rsp, 40\n")
	gen.text.WriteString("  ret\n")
	gen.text.WriteString("\n")

	gen.text.WriteString("_start:\n")
	for idx, inst := range program {
		gen.text.WriteString(fmt.Sprintf("addr_%d:\n", idx))
		generateX8664Instruction(&gen, inst)
	}

	gen.text.WriteString("\n")
	gen.text.WriteString(fmt.Sprintf("addr_%d:\n", len(program)))
	gen.text.WriteString("  ;; exit syscall\n")
	gen.text.WriteString("  mov rax, 0x3c\n")
	gen.text.WriteString("  mov rdi, 0\n")
	gen.text.WriteString("  syscall\n")

	// Data section
	gen.text.WriteString("\n")
	gen.text.WriteString("section .data\n")
	for idx, str := range gen.strings {
		gen.text.WriteString(fmt.Sprintf("str_%d: db \"%s\"\n", idx, str))
	}

	return gen.text.String()
}

func generateX8664Instruction(gen *x86_64Generator, inst Instruction) {
	switch inst.Kind {
	case InstKindPushInt:
		gen.text.WriteString("  ;; push int\n")
		gen.text.WriteString(fmt.Sprintf("  mov rax, %d\n", inst.ValueInt))
		gen.text.WriteString("  push rax\n")
	case InstKindPushString:
		gen.text.WriteString("  ;; push string\n")
		gen.text.WriteString(fmt.Sprintf("  mov rax, %d\n", len(inst.ValueString)))
		gen.text.WriteString("  push rax\n")
		gen.text.WriteString(fmt.Sprintf("  push str_%d", len(gen.strings)))
		gen.strings = append(gen.strings, inst.ValueString)
	case InstKeyword:
		generateX8664Keyword(gen, inst.ValueKeyword)
	case InstKindIntrinsic:
		generateX8664Intrinsic(gen, inst.ValueIntrinsic)
	default:
		panic(fmt.Sprintf("unknown instruction kind '%s'", inst.Kind))
	}
}

func generateX8664Keyword(gen *x86_64Generator, keyword Keyword) {
	switch keyword.Kind {
	case KeywordKindIf:
		gen.text.WriteString("  ;; if\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  test rax rax\n")
		gen.text.WriteString(fmt.Sprintf("  jz addr_%d\n", keyword.JmpAddress))
	case KeywordKindElse:
		gen.text.WriteString("  ;; else\n")
	case KeywordKindEnd:
		gen.text.WriteString("  ;; end\n")
		if keyword.HasJmp {
			gen.text.WriteString(fmt.Sprintf("  jmp addr_%d\n", keyword.JmpAddress))
		}
	case KeywordKindWhile:
		gen.text.WriteString("  ;; while\n")
	case KeywordKindDo:
		gen.text.WriteString("  ;; do\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  test rax rax\n")
		gen.text.WriteString(fmt.Sprintf("  jz addr_%d\n", keyword.JmpAddress))
	default:
		panic(fmt.Sprintf("unknown keyword '%s'", keyword.Kind))
	}
}

func generateX8664Intrinsic(gen *x86_64Generator, intrinsic Intrinsic) {
	switch intrinsic {
	case IntrinsicPlus:
		gen.text.WriteString("  ;; add\n")
		gen.text.WriteString("  pop rbx\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  add rax, rbx\n")
		gen.text.WriteString("  push rax\n")
	case IntrinsicMinus:
		gen.text.WriteString("  ;; sub\n")
		gen.text.WriteString("  pop rbx\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  sub rax, rbx\n")
		gen.text.WriteString("  push rax\n")
	case IntrinsicTimes:
		gen.text.WriteString("  ;; mul\n")
		gen.text.WriteString("  pop rbx\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  imul rbx\n")
		gen.text.WriteString("  push rax\n")
	case IntrinsicDivMod:
		gen.text.WriteString("  ;; divmod\n")
		gen.text.WriteString("  pop rbx\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  idiv rbx\n")
		gen.text.WriteString("  push rax\n")
		gen.text.WriteString("  push rdx\n")
	case IntrinsicGreather:
		gen.text.WriteString("  ;; greather\n")
		gen.text.WriteString("  mov rcx, 0\n")
		gen.text.WriteString("  mov rdx, 1\n")
		gen.text.WriteString("  pop rbx\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  cmp rax, rbx\n")
		gen.text.WriteString("  cmovg rcx, rdx\n")
		gen.text.WriteString("  push rcx\n")
	case IntrinsicLess:
		gen.text.WriteString("  ;; greather\n")
		gen.text.WriteString("  mov rcx, 0\n")
		gen.text.WriteString("  mov rdx, 1\n")
		gen.text.WriteString("  pop rbx\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  cmp rax, rbx\n")
		gen.text.WriteString("  cmovl rcx, rdx\n")
		gen.text.WriteString("  push rcx\n")
	case IntrinsicNotEqual:
		gen.text.WriteString("  ;; not equal\n")
		gen.text.WriteString("  mov rcx, 0\n")
		gen.text.WriteString("  mov rdx, 1\n")
		gen.text.WriteString("  pop rbx\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  cmp rax, rbx\n")
		gen.text.WriteString("  cmovne rcx, rdx\n")
		gen.text.WriteString("  push rcx\n")
	case IntrinsicDup:
		gen.text.WriteString("  ;; dup\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  push rax\n")
		gen.text.WriteString("  push rax\n")
	case IntrinsicPrint:
		gen.text.WriteString("  ;; print\n")
		gen.text.WriteString("  pop rdi\n")
		gen.text.WriteString("  call print\n")
	default:
		panic(fmt.Sprintf("unknown intrinsic '%s'", intrinsic))
	}
}
