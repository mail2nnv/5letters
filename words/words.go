package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"unicode"
	"unicode/utf8"
)

type rule struct {
	should map[int]rune
	may    map[rune]int
	stop   map[rune]bool
}

func makeRule(arg string) (rule, error) {
	r := rule{
		should: map[int]rune{},
		may:    map[rune]int{},
		stop:   map[rune]bool{},
	}
	err := r.fill(arg)
	return r, err
}

// arg `ПИЛКА-02001`
func (r *rule) fill(arg string) error {
	if utf8.RuneCountInString(arg) != 11 {
		return fmt.Errorf("invalid input «%s», len should be 11", arg)
	}
	m := make([]struct {
		rune
		int
	}, 5)
	i := 0
	for _, rune := range arg {
		switch i {
		case 0, 1, 2, 3, 4:
			if !unicode.IsLetter(rune) {
				return fmt.Errorf("invalid input «%s», not a letter `%c` at index %d", arg, rune, i)
			}
			m[i].rune = unicode.ToLower(rune)
		case 6, 7, 8, 9, 10:
			if (rune != '0') && (rune != '1') && (rune != '2') {
				return fmt.Errorf("invalid input «%s», not a indicator `%c` at index %d", arg, rune, i)
			}
			m[i-6].int = int(rune) - 48
		}
		i++
	}

	for i, m := range m {
		switch m.int {
		case 0:
			r.stop[m.rune] = true
		case 1:
			r.may[m.rune] = i
		case 2:
			r.should[i] = m.rune
		}
	}

	return nil
}

func (r rule) match(w string) bool {
	i := 0
	for _, rune := range w {
		rune = unicode.ToLower(rune)
		if should, ok := r.should[i]; ok {
			if rune != should {
				return false
			}
			i++
			continue
		}
		if r.stop[rune] {
			return false
		}
		i++
	}
	for may, pos := range r.may {
		ok := false
		i := 0
		for _, rune := range w {
			rune = unicode.ToLower(rune)
			if rune == may {
				ok = i != pos
			}
			i++
		}
		if !ok {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		panic(errors.New("missed command line arguments"))
	}

	rules := make([]rule, 0, 6)

	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		r, err := makeRule(arg)
		if err != nil {
			panic(err)
		}
		rules = append(rules, r)
	}

	func() {
		src, err := os.Open("all.txt")
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		defer src.Close()

		scanner := bufio.NewScanner(src)
		for scanner.Scan() {
			w := scanner.Text()
			if utf8.RuneCountInString(w) == 5 {
				ok := true
				for _, r := range rules {
					if !r.match(w) {
						ok = false
						break
					}
				}
				if ok {
					fmt.Println(w)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
			panic(err)
		}
	}()
}
