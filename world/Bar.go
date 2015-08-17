package world

import "fmt"

type Bar struct {
	Value   int
	Maximum int
}

func (bar Bar) Render(width int) string {
	total := width - 2
	amount := bar.Value * total / bar.Maximum
	text := []rune(fmt.Sprintf("%d/%d", bar.Value, bar.Maximum))
	padding := (total - len(text)) / 2
	if amount < 0 {
		amount = 1
	}
	result := "{f:white|"
	result += "{b:red|▌"
	for i := 0; i < total; i++ {
		if i == amount {
			result += "}"
		}
		if i >= padding && i < padding+len(text) {
			result += string([]rune{text[i-padding]})
		} else {
			result += " "
		}
	}
	if amount == total {
		result += "▐}"
		result += "}"
	} else {
		result += "▐}"
	}

	return result
}
