package theme

import lg "github.com/charmbracelet/lipgloss"

var DefaultText = lg.NewStyle().
	Foreground(lg.Color("#FFFFFF"))

var CursorText = lg.NewStyle().
	Foreground(lg.Color("#000000")).
	Background(lg.Color("#FFFFFF"))

var NumberColors = [9]lg.Color{
	lg.Color("236"), // 0
	lg.Color("21"),  // 1
	lg.Color("2"),   // 2
	lg.Color("9"),   // 3
	lg.Color("5"),   // 4
	lg.Color("1"),   // 5
	lg.Color("6"),   // 6
	lg.Color("7"),   // 7
	lg.Color("8"),   // 8
}
