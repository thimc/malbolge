package main

import (
	"fmt"
	"io"
	"os"
)

const (
	MemorySize   = 59048
	encode       = "+b(29e*j1VMEKLyC})8&m#~W>qxdRp0wkrUo[D7,XTcA\"lI.v%{gJh4G\\-=O@5`_3i<?Z';FNQuY]szf$!BS/|t:Pn6^Ha"
	decode       = "5z]&gqtyfr$(we4{WP)H-Zn,[%\\3dL+Q;>U!pJS72FhOA1CB6v^=I_0/8|jsb9m<.TVac`uY*MK'X~xDl}REokN:#?G\"i@"
	validOpcodes = "i</*jpov"
	dataLower    = 33
	dataUpper    = 126
)

type Malbolge struct {
	memory []uint32
}

func NewMalbolge(memorySize int, program []byte) (*Malbolge, error) {
	m := &Malbolge{
		memory: make([]uint32, memorySize+1),
	}

	line := 1
	tok := 0
	n := 0
	for _, ch := range program {
		tok++
		if ch == '\r' || ch == '\n' || ch == ' ' || ch == '\t' {
			tok = 0
			line++
			continue
		}
		if ch > 32 && ch < 127 {
			if !m.valid(ch, n) {
				return nil, fmt.Errorf("unknown instruction at line %d, character %d: \"%c\"", line, tok, ch)
			}
		}
		if n >= len(m.memory) {
			return nil, fmt.Errorf("input file too long")
		}
		m.memory[n] = uint32(ch)
		n += 1
	}
	for n < len(m.memory) {
		m.memory[n] = m.crazy(m.memory[n-1], m.memory[n-2])
		n += 1
	}
	return m, nil
}

func (m *Malbolge) valid(ch byte, idx int) bool {
	op := encode[(int(ch)-33+idx)%94]
	for _, opcode := range validOpcodes {
		if opcode == rune(op) {
			return true
		}
	}
	return false
}

func (m *Malbolge) rotr(x uint32) uint32 {
	return x/3 + x%3*19683
}

func (m *Malbolge) crazy(x, y uint32) uint32 {
	var result uint32
	p9 := [5]uint32{1, 9, 81, 729, 6561}
	o := [9][9]uint32{
		{4, 3, 3, 1, 0, 0, 1, 0, 0},
		{4, 3, 5, 1, 0, 2, 1, 0, 2},
		{5, 5, 4, 2, 2, 1, 2, 2, 1},
		{4, 3, 3, 1, 0, 0, 7, 6, 6},
		{4, 3, 5, 1, 0, 2, 7, 6, 8},
		{5, 5, 4, 2, 2, 1, 8, 8, 7},
		{7, 6, 6, 7, 6, 6, 4, 3, 3},
		{7, 6, 8, 7, 6, 8, 4, 3, 5},
		{8, 8, 7, 8, 8, 7, 5, 5, 4},
	}
	for i := 0; i < 5; i++ {
		result += o[y/p9[i]%9][x/p9[i]%9] * p9[i]
	}
	return result
}

func (m *Malbolge) Run() {
	var c, d, a uint32
	for {
		if m.memory[c] < dataLower || m.memory[c] > dataUpper {
			fmt.Fprintf(os.Stderr, "ERROR: memory[%d] [C] = %d is not a valid instruction\n", c, m.memory[c])
			continue
		}
		switch encode[(m.memory[c]-33+c)%94] {
		case 'j':
			d = m.memory[d]
		case 'i':
			c = m.memory[d]
		case '*':
			m.memory[d] = m.rotr(m.memory[d])
			a = m.memory[d]
		case 'p':
			m.memory[d] = m.crazy(a, m.memory[d])
			a = m.memory[d]
		case '<':
			fmt.Printf("%c", a&0xFF)
		case '/':
			var ch byte
			_, err := fmt.Scanf("%c", &ch)
			if err != nil {
				a = MemorySize
			} else {
				a = uint32(ch)
			}
		case 'v':
			return
		}

		m.memory[c] = uint32(decode[m.memory[c]-33])

		if c >= MemorySize {
			c = 0
		} else {
			c += 1
		}

		if d >= MemorySize {
			d = 0
		} else {
			d += 1
		}
	}
}

func main() {
	file := os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		file = f
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	if len(content) < 1 {
		fmt.Fprintf(os.Stderr, "usage: %s <file.mal>\n", os.Args[0])
		os.Exit(1)
	}

	m, err := NewMalbolge(MemorySize, content)
	if err != nil {
		panic(err)
	}
	m.Run()
	fmt.Println()
}
