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
	gen.WriteString("print:\n")
	gen.WriteString("  mov     r9, -3689348814741910323\n")
	gen.WriteString("  sub     rsp, 40\n")
	gen.WriteString("  mov     BYTE [rsp+31], 10\n")
	gen.WriteString("  lea     rcx, [rsp+30]\n")
	gen.WriteString(".L2:\n")
	gen.WriteString("  mov     rax, rdi\n")
	gen.WriteString("  lea     r8, [rsp+32]\n")
	gen.WriteString("  mul     r9\n")
	gen.WriteString("  mov     rax, rdi\n")
	gen.WriteString("  sub     r8, rcx\n")
	gen.WriteString("  shr     rdx, 3\n")
	gen.WriteString("  lea     rsi, [rdx+rdx*4]\n")
	gen.WriteString("  add     rsi, rsi\n")
	gen.WriteString("  sub     rax, rsi\n")
	gen.WriteString("  add     eax, 48\n")
	gen.WriteString("  mov     BYTE [rcx], al\n")
	gen.WriteString("  mov     rax, rdi\n")
	gen.WriteString("  mov     rdi, rdx\n")
	gen.WriteString("  mov     rdx, rcx\n")
	gen.WriteString("  sub     rcx, 1\n")
	gen.WriteString("  cmp     rax, 9\n")
	gen.WriteString("  ja      .L2\n")
	gen.WriteString("  lea     rax, [rsp+32]\n")
	gen.WriteString("  mov     edi, 1\n")
	gen.WriteString("  sub     rdx, rax\n")
	gen.WriteString("  xor     eax, eax\n")
	gen.WriteString("  lea     rsi, [rsp+32+rdx]\n")
	gen.WriteString("  mov     rdx, r8\n")
	gen.WriteString("  mov     rax, 1\n")
	gen.WriteString("  syscall\n")
	gen.WriteString("  add     rsp, 40\n")
	gen.WriteString("  ret\n")
	gen.WriteString("\n")

	gen.WriteString("_start:\n")
	for idx, inst := range program {
		gen.WriteString(fmt.Sprintf("addr_%d:\n", idx))
		generateX8664Instruction(&gen, inst)
	}

	gen.WriteString("\n")
	gen.WriteString(fmt.Sprintf("addr_%d:\n", len(program)))
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
	case InstKeyword:
		generateX8664Keyword(gen, inst.ValueKeyword)
	case InstKindIntrinsic:
		generateX8664Intrinsic(gen, inst.ValueIntrinsic)
	default:
		panic(fmt.Sprintf("unknown instruction kind '%s'", inst.Kind))
	}
}

func generateX8664Keyword(gen *strings.Builder, keyword Keyword) {
	switch keyword.Kind {
	case KeywordKindIf:
		gen.WriteString("  ;; if\n")
		gen.WriteString("  pop rax\n")
		gen.WriteString("  test rax rax\n")
		gen.WriteString(fmt.Sprintf("  jz addr_%d\n", keyword.JmpAddress))
	case KeywordKindElse:
		gen.WriteString("  ;; else\n")
	case KeywordKindEnd:
		gen.WriteString("  ;; end\n")
		if keyword.HasJmp {
			gen.WriteString(fmt.Sprintf("  jmp addr_%d\n", keyword.JmpAddress))
		}
	case KeywordKindWhile:
		gen.WriteString("  ;; while\n")
	case KeywordKindDo:
		gen.WriteString("  ;; do\n")
		gen.WriteString("  pop rax\n")
		gen.WriteString("  test rax rax\n")
		gen.WriteString(fmt.Sprintf("  jz addr_%d\n", keyword.JmpAddress))
	default:
		panic(fmt.Sprintf("unknown keyword '%s'", keyword.Kind))
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
	case IntrinsicGreather:
		gen.WriteString("  ;; greather\n")
		gen.WriteString("  mov rcx, 0\n")
		gen.WriteString("  mov rdx, 1\n")
		gen.WriteString("  pop rbx\n")
		gen.WriteString("  pop rax\n")
		gen.WriteString("  cmp rax, rbx\n")
		gen.WriteString("  cmovg rcx, rdx\n")
		gen.WriteString("  push rcx\n")
	case IntrinsicLess:
		gen.WriteString("  ;; greather\n")
		gen.WriteString("  mov rcx, 0\n")
		gen.WriteString("  mov rdx, 1\n")
		gen.WriteString("  pop rbx\n")
		gen.WriteString("  pop rax\n")
		gen.WriteString("  cmp rax, rbx\n")
		gen.WriteString("  cmovl rcx, rdx\n")
		gen.WriteString("  push rcx\n")
	case IntrinsicNotEqual:
		gen.WriteString("  ;; not equal\n")
		gen.WriteString("  mov rcx, 0\n")
		gen.WriteString("  mov rdx, 1\n")
		gen.WriteString("  pop rbx\n")
		gen.WriteString("  pop rax\n")
		gen.WriteString("  cmp rax, rbx\n")
		gen.WriteString("  cmovne rcx, rdx\n")
		gen.WriteString("  push rcx\n")
	case IntrinsicDup:
		gen.WriteString("  ;; dup\n")
		gen.WriteString("  pop rax\n")
		gen.WriteString("  push rax\n")
		gen.WriteString("  push rax\n")
	case IntrinsicPrint:
		gen.WriteString("  ;; print\n")
		gen.WriteString("  pop rdi\n")
		gen.WriteString("  call print\n")
	default:
		panic(fmt.Sprintf("unknown intrinsic '%s'", intrinsic))
	}
}
