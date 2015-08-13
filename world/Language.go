package world

type Person int

const (
	Third Person = iota // hence it is also the default
	Second
	First
	Plural
)

func MakeIndefinite(word string) string {
	if word == "" {
		return ""
	}
	switch word[0] {
	case 'a', 'e', 'i', 'o', 'u':
		return "an " + word
	}
	return "a " + word
}

func ConjugatePast(word string, person Person) string {
	switch word {
	case "be":
		switch person {
		case First, Third:
			return "was"
		case Second, Plural:
			return "were"
		}
	case "have":
		return "had"
	}
	return word + "ed"
}

func ConjugatePresent(word string, person Person) string {
	if word == "" {
		panic("empty word for conjugation")
	}
	switch word {
	case "be":
		switch person {
		case First:
			return "am"
		case Second, Plural:
			return "are"
		case Third:
			return "is"
		}
	case "have":
		switch person {
		case First, Second, Plural:
			return "have"
		case Third:
			return "has"
		}

	}
	// word: "learn"
	// I learn
	// You learn
	// He learns
	// We/They learn
	if person == Third {
		if word[len(word)-1] == 's' {
			return word + "es"
		}
		return word + "s"
	} else {
		return word
	}
}
