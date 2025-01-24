package main

import (
    "database/sql"
    "fmt"

    "os"

    tea "github.com/charmbracelet/bubbletea"
    _ "github.com/lib/pq"
)

type model struct {
    choices  []string
    cursor   int
    selected map[int]struct{}
    adding   bool
    newMovie string
}

func initialModel() model {
    return model{
        choices:  []string{"Attack on Titan", "Dark Knight", "Inglourious Basterds"},
        selected: make(map[int]struct{}),
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        case "a": // Start adding a new movie
            m.adding = true
            m.newMovie = ""
        case "s": // Save to database
            return m, saveToDatabase(m.choices, m.selected)
        }
    case tea.WindowSizeMsg:
        // Handle resizing if needed

    return m, nil
}

func (m model) View() string {
    if m.adding {
        return fmt.Sprintf("Enter new movie title: %s\n(Press Enter to add, Esc to cancel)", m.newMovie)
    }

    s := "What should we watch next?\n\n"

    for i, choice := range m.choices {
        cursor := ""
        if m.cursor == i {
            cursor = ">"
        }

        checked := " " // not selected
        if _, ok := m.selected[i]; ok {
            checked = "x" // selected!
        }

        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }
    s += "\nPress q to quit, a to add, s to save.\n"
    return s
}

// Saves the movie list to the database
func saveToDatabase(movies []string, selected map[int]struct{}) tea.Cmd {
    return func() tea.Msg {
        connStr := "host=localhost port=5432 user=postgres dbname=movie_watch_cli password= sslmode=disable"

        db, err := sql.Open("postgres", connStr)
        if err != nil {
            return fmt.Sprintf("Could not connect to database: %v", err)
        }
        defer db.Close()

        _, err = db.Exec("DELETE FROM watchlist")
        if err != nil {
            return fmt.Sprintf("Could not delete existing entries: %v", err)
        }

        for i, movie := range movies {
            selectedValue := 0
            if _, ok := selected[i]; ok {
                selectedValue = 1
            }
            _, err = db.Exec("INSERT INTO watchlist (title, selected) VALUES ($1, $2)", movie, selectedValue)
            if err != nil {
                return fmt.Sprintf("Could not insert movie: %v", err)
            }
        }

        return "Watchlist saved to database."
    }
}

func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v\n", err)
        os.Exit(1)
    }
}
