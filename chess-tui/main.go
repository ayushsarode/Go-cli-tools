package main

import (
	"fmt"
	"log"
	"strings"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/notnil/chess"
)

type model struct {
	game   *chess.Game
	input  string
	status string
}

func initialModel() model {
	return model{
		game:   chess.NewGame(),
		status: "Enter your move in UCI format (e.g., e2e4):",
	}
}

func (m model) Init() bubbletea.Cmd {
	return nil
}

func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, bubbletea.Quit

		case "enter":
			move, err := chess.UCINotation{}.Decode(m.game.Position(), strings.TrimSpace(m.input))
			if err == nil && m.game.Move(move) == nil {
				m.input = ""
			} else {
				m.status = "Invalid move! Try again:"
			}
		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m model) View() string {
	board := m.game.Position().Board().String()
	return fmt.Sprintf(
		"%s\n\n%s\n\nMove: %s\n[Press 'q' to quit]",
		board, m.status, m.input,
	)
}

func main() {
	p := bubbletea.NewProgram(initialModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
