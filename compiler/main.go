package main

import (
	"io"
	"os"
)

func main() {
	m := Machine{
		buf:    make([]byte, 1),
		input:  os.Stdin,
		output: os.Stdout,
	}

	helloWorld := `
	++++++++[>++++[>++>+++>+++>+<<
	<<-]>+>+>->>+[<]<-]>>.>---.+++
	++++..+++.>>.<-.<.+++.------.-
	-------.>>+.>++.
	`

	m.code = Compile(helloWorld)
	// for i, v := range m.code {
	// 	fmt.Printf("%d(%c %d) ", i, v.Type, v.Arg)
	// }

	Execute(&m)
}

type Machine struct {
	code []Instruction
	ip   int

	memory [20]byte
	dp     int

	input  io.Reader
	output io.Writer
	buf    []byte
}

func Execute(m *Machine) {
	for m.ip < len(m.code) {
		switch ins := m.code[m.ip]; ins.Type {
		case '+':
			m.memory[m.dp] += byte(ins.Arg)
		case '-':
			m.memory[m.dp] -= byte(ins.Arg)
		case '>':
			m.dp += ins.Arg
		case '<':
			m.dp -= ins.Arg
		case ',':
			for i := 0; i < ins.Arg; i++ {
				m.input.Read(m.buf)
				m.memory[m.dp] = m.buf[0]
			}
		case '.':
			for i := 0; i < ins.Arg; i++ {
				m.buf[0] = m.memory[m.dp]
				m.output.Write(m.buf)
			}
		case '[':
			if m.memory[m.dp] == 0 {
				m.ip = ins.Arg
				continue
			}
		case ']':
			if m.memory[m.dp] != 0 {
				m.ip = ins.Arg
				continue
			}
		}

		m.ip++
	}
}

type Instruction struct {
	Type byte
	Arg  int
}

func Compile(input string) []Instruction {
	pos := 0
	loopStack := []int{}
	program := []Instruction{}

	for pos < len(input) {
		switch tok := input[pos]; tok {
		case '+', '-', '>', '<', ',', '.':
			start := pos
			for pos < len(input)-1 && input[pos+1] == tok {
				pos++
			}
			length := pos - start + 1
			program = append(program, Instruction{Type: tok, Arg: length})

		case '[':
			// dummy value into Arg for now - will update in the ']' case
			program = append(program, Instruction{Type: tok, Arg: 0})

			insPos := len(program) - 1
			loopStack = append(loopStack, insPos)

		case ']':
			if len(loopStack) == 0 {
				panic("missing opening bracket '['")
			}
			// pop position of last opening bracket off the stack
			openInsPos := loopStack[len(loopStack)-1]
			loopStack = loopStack[:len(loopStack)-1]

			program = append(program, Instruction{Type: tok, Arg: openInsPos})
			insPos := len(program) - 1

			// backfill correct instruction position for matching opening bracket
			program[openInsPos].Arg = insPos
		}

		pos++
	}

	if len(loopStack) != 0 {
		panic("missing closing bracket ']'")
	}

	return program
}
