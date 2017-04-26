package ec2go

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"strings"
)

const (
	QueryPrompt      string = "[Query]> "
	InitialCursorIdx        = 8
)

var page int
var itemPerPage = 10
var currentPage = 0

type Screen struct {
	prompt         string
	cursorIdx      int
	selectedLine   int
	input          []rune
	candidates     []string
	originContents []string
}

func (s *Screen) Run() {
	output := ""
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer func() {
		termbox.Close()
		if output != "" {
			fmt.Println(output)
		}
	}()

	s.draw()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			//s.draw()
		case termbox.EventKey:
			switch ev.Key {
			case 0:
				s.inputChar(ev.Ch)
			case termbox.KeyEsc, termbox.KeyCtrlC:
				break loop
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				s.deleteChar()
			case termbox.KeyEnter:
				output = s.selectCandidate()
				break loop
			case termbox.KeyArrowDown:
				s.selectNextLine()
			case termbox.KeyArrowUp:
				s.selectPreviousLine()
			case termbox.KeyArrowLeft:
				if currentPage != 0 {
					currentPage--
				}
				s.draw()
			case termbox.KeyArrowRight:
				if page > currentPage+1 {
					currentPage++
				}
				s.draw()
			}
		default:
			//s.draw()
		}
	}
}

func (s *Screen) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	_, height := termbox.Size()

	itemPerPage = height - 2

	page = len(s.originContents) / itemPerPage
	if len(s.originContents)%itemPerPage != 0 {
		page++
	}

	start := currentPage * itemPerPage
	end := start + itemPerPage
	if len(s.originContents[start:]) < itemPerPage {
		end = start + len(s.originContents[start:])
	}

	s.candidates = s.originContents[start:end]

	s.drawQueryLine(string(s.input))
	s.drawCandidates()
	s.setHighlight(1)
	termbox.SetCursor(s.cursorIdx+1, 0)
	termbox.Flush()
}

func (s *Screen) drawQueryLine(qs string) {
	fs := s.prompt + qs
	cells := []termbox.Cell{}

	for _, s := range fs {
		cells = append(cells, termbox.Cell{
			Ch: s,
			Fg: termbox.ColorDefault,
			Bg: termbox.ColorDefault,
		})
	}
	s.drawCells(0, 0, cells)
}

func (s *Screen) drawCells(x, y int, cells []termbox.Cell) {
	for i := 0; i < len(cells); i++ {
		termbox.SetCell(x+i, y, cells[i].Ch, cells[i].Fg, cells[i].Bg)
	}
}

func (s *Screen) inputChar(ch rune) {
	defer s.setHighlight(1)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	s.input = append(s.input, ch)
	s.cursorIdx++
	s.drawQueryLine(string(s.input))
	s.filterCandidates()
	s.drawCandidates()
	termbox.SetCursor(s.cursorIdx+1, 0)
	termbox.Flush()
}

func (s *Screen) filterCandidates() {
	s.candidates = []string{}
	tempContents := []string{}
	for _, c := range s.originContents {
		if strings.ContainsAny(c, string(s.input)) {
			tempContents = append(tempContents, c)
		}
	}
	// pager
	_, height := termbox.Size()

	itemPerPage = height - 2

	page = len(s.originContents) / itemPerPage
	if len(s.originContents)%itemPerPage != 0 {
		page++
	}

	start := currentPage * itemPerPage
	end := start + itemPerPage
	if len(s.originContents[start:]) < itemPerPage {
		end = start + len(s.originContents[start:])
	}

	s.candidates = tempContents[start:end]
	//s.candidates = tempContents
}

func (s *Screen) drawCandidates() {
	for rowIdx, c := range s.candidates {
		var cells []termbox.Cell
		for _, s := range c {
			cells = append(cells, termbox.Cell{
				Ch: s,
				Fg: termbox.ColorDefault,
				Bg: termbox.ColorDefault,
			})
		}
		s.drawCells(0, rowIdx+1, cells)
	}
}

func (s *Screen) deleteChar() {
	lastIdx := len(s.input)
	defer s.setHighlight(1)
	if lastIdx > 1 {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		s.input = s.input[:lastIdx-1]
		s.cursorIdx--
		s.drawQueryLine(string(s.input))
		s.filterCandidates()
		s.drawCandidates()
		termbox.SetCursor(s.cursorIdx+1, 0)
		termbox.Flush()
	} else if lastIdx == 1 {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		s.candidates = s.originContents
		s.input = []rune("")
		s.cursorIdx--
		s.draw()
	}
}

func (s *Screen) selectNextLine() {
	if s.selectedLine == len(s.candidates) {
		s.selectedLine = 0
	} else {
		s.selectedLine++
	}
	s.setHighlight(s.selectedLine)
}

func (s *Screen) selectPreviousLine() {
	if s.selectedLine == 0 {
		s.selectedLine = len(s.candidates)
	} else {
		s.selectedLine--
	}
	s.setHighlight(s.selectedLine)
}

func (s *Screen) selectCandidate() string {
	return s.candidates[s.selectedLine-1]
}

func (s *Screen) setHighlight(targetLine int) {
	width, _ := termbox.Size()
	s.selectedLine = targetLine
	for row := 1; row <= len(s.candidates); row++ {
		bgColor := termbox.ColorDefault
		if row == targetLine {
			bgColor = termbox.ColorGreen
		}
		for col := 0; col < width; col++ {
			char := termbox.CellBuffer()[(width*row)+col].Ch
			termbox.SetCell(col, row, char, termbox.ColorDefault, bgColor)
		}
	}
	termbox.Flush()
}
