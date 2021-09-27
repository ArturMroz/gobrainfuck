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

	m.code = helloWorld
	m.Execute()
}

type Machine struct {
	code string
	ip   int

	memory [100]byte
	dp     int

	input  io.Reader
	output io.Writer
	buf    []byte
}

func (m *Machine) Execute() {
	for m.ip < len(m.code) {
		switch ins := m.code[m.ip]; ins {
		case '+':
			m.memory[m.dp]++
		case '-':
			m.memory[m.dp]--
		case '>':
			m.dp++
		case '<':
			m.dp--
		case ',':
			m.input.Read(m.buf)
			m.memory[m.dp] = m.buf[0]
		case '.':
			m.buf[0] = m.memory[m.dp]
			m.output.Write(m.buf)
		case '[':
			if m.memory[m.dp] == 0 {
				depth := 1
				for depth != 0 {
					m.ip++
					switch m.code[m.ip] {
					case '[':
						depth++
					case ']':
						depth--
					}
				}
			}
		case ']':
			if m.memory[m.dp] != 0 {
				depth := 1
				for depth != 0 {
					m.ip--
					switch m.code[m.ip] {
					case '[':
						depth--
					case ']':
						depth++
					}
				}
			}
		}

		m.ip++
	}
}
