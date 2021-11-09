package tin

import (
	"fmt"
	"strings"
)

type x86_64Generator struct {
	text    strings.Builder
	strings []string
}

const (
	addressPrefix string = "addr"
	stringPrefix  string = "str"
)

func generateNasmX8664(program Program) string {
	gen := x86_64Generator{}

	// Text section
	gen.text.WriteString("section .text\n")
	gen.text.WriteString("global _start\n")

	// Main function
	{
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
	}

	gen.text.WriteString("_start:\n")
	gen.text.WriteString("  mov rax, ret_stack\n")
	gen.text.WriteString("  mov [ret_base], rax\n")

	for idx, inst := range program {
		gen.text.WriteString(fmt.Sprintf("%s:\n", getAddrName(idx)))
		generateX8664Instruction(&gen, inst)
	}

	gen.text.WriteString("\n")
	gen.text.WriteString(fmt.Sprintf("%s:\n", getAddrName(len(program))))
	gen.text.WriteString("  ;; exit syscall\n")
	gen.text.WriteString("  mov rax, 0x3c\n")
	gen.text.WriteString("  mov rdi, 0\n")
	gen.text.WriteString("  syscall\n")

	// Data section
	gen.text.WriteString("\n")
	gen.text.WriteString("section .data\n")
	for idx, str := range gen.strings {
		gen.text.WriteString(fmt.Sprintf("%s: db `%s`\n", getStringName(idx), str))
	}

	// Bss section
	gen.text.WriteString("\n")
	gen.text.WriteString("section .bss\n")
	gen.text.WriteString("	ret_base: resq 1\n")
	gen.text.WriteString("	ret_stack: resb 1024\n")

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
		gen.text.WriteString(fmt.Sprintf("  push str_%d\n", len(gen.strings)))
		gen.strings = append(gen.strings, inst.ValueString)
	case InstKindTestCondition:
		gen.text.WriteString("  ;; test condition\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  test rax, rax\n")
		gen.text.WriteString(fmt.Sprintf("  jz %s\n", getAddrName(inst.JmpAddress)))
	case InstKindElse:
		gen.text.WriteString("  ;; else\n")
		gen.text.WriteString(fmt.Sprintf("  jmp %s\n", getAddrName(inst.JmpAddress)))
	case InstKindWhile:
		gen.text.WriteString("  ;; while\n")
	case InstKindEnd:
		gen.text.WriteString("  ;; end\n")
		gen.text.WriteString(fmt.Sprintf("  jmp %s\n", getAddrName(inst.JmpAddress)))
	case InstKindFunSkip:
		gen.text.WriteString("  ;; fun skip\n")
		gen.text.WriteString(fmt.Sprintf("  jmp %s\n", getAddrName(inst.JmpAddress)))
	case InstKindFunDef:
		gen.text.WriteString("  ;; fun def\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  mov rbx, [ret_base]\n")
		gen.text.WriteString("  mov [rbx], rax\n")
		gen.text.WriteString("  add rbx, 8\n")
		gen.text.WriteString("  mov [ret_base], rbx\n")
	case InstKindFunRet:
		gen.text.WriteString("  ;; fun ret\n")
		gen.text.WriteString("  mov rbx, [ret_base]\n")
		gen.text.WriteString("  sub rbx, 8\n")
		gen.text.WriteString("  mov [ret_base], rbx\n")
		gen.text.WriteString("  mov rax, [rbx]\n")
		gen.text.WriteString("  push rax\n")
		gen.text.WriteString("  ret\n")
	case InstKindFunCall:
		gen.text.WriteString("  ;; fun call\n")
		gen.text.WriteString(fmt.Sprintf("  call %s\n", getAddrName(inst.JmpAddress)))
	case InstKindIntrinsic:
		generateX8664Intrinsic(gen, inst.ValueIntrinsic)
	default:
		panic(fmt.Sprintf("unknown instruction kind '%s'", inst.Kind))
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
		gen.text.WriteString("  xor rdx, rdx\n")
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
		gen.text.WriteString("  ;; less\n")
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
	case IntrinsicSyscall0:
		gen.text.WriteString("  ;; syscall0\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  syscall\n")
	case IntrinsicSyscall1:
		gen.text.WriteString("  ;; syscall1\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  pop rdi\n")
		gen.text.WriteString("  syscall\n")
	case IntrinsicSyscall2:
		gen.text.WriteString("  ;; syscall2\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  pop rdi\n")
		gen.text.WriteString("  pop rsi\n")
		gen.text.WriteString("  syscall\n")
	case IntrinsicSyscall3:
		gen.text.WriteString("  ;; syscall3\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  pop rdi\n")
		gen.text.WriteString("  pop rsi\n")
		gen.text.WriteString("  pop rdx\n")
		gen.text.WriteString("  syscall\n")
	case IntrinsicSyscall4:
		gen.text.WriteString("  ;; syscall4\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  pop rdi\n")
		gen.text.WriteString("  pop rsi\n")
		gen.text.WriteString("  pop rdx\n")
		gen.text.WriteString("  pop r10\n")
		gen.text.WriteString("  syscall\n")
	case IntrinsicSyscall5:
		gen.text.WriteString("  ;; syscall5\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  pop rdi\n")
		gen.text.WriteString("  pop rsi\n")
		gen.text.WriteString("  pop rdx\n")
		gen.text.WriteString("  pop r10\n")
		gen.text.WriteString("  pop r8\n")
		gen.text.WriteString("  syscall\n")
	case IntrinsicSyscall6:
		gen.text.WriteString("  ;; syscall6\n")
		gen.text.WriteString("  pop rax\n")
		gen.text.WriteString("  pop rdi\n")
		gen.text.WriteString("  pop rsi\n")
		gen.text.WriteString("  pop rdx\n")
		gen.text.WriteString("  pop r10\n")
		gen.text.WriteString("  pop r8\n")
		gen.text.WriteString("  pop r9\n")
		gen.text.WriteString("  syscall\n")
	default:
		panic(fmt.Sprintf("unknown intrinsic '%s'", intrinsic))
	}
}

func getAddrName(addr int) string {
	return fmt.Sprintf("%s_%d", addressPrefix, addr)
}

func getStringName(strNum int) string {
	return fmt.Sprintf("%s_%d", stringPrefix, strNum)
}
