package main

import (
	"bufio"
	"fmt"
	"iter"
	"log"
	"maps"
	"os"
	"slices"
	"unicode/utf8"
)

type word = string

func doubles(w word) bool {
	uniq := map[rune]bool{}
	for _, r := range w {
		if uniq[r] {
			return true
		}
		uniq[r] = true
	}
	return false
}

type alphabet struct {
	runes map[rune]int
	all   int
}

func newAlphabet() *alphabet {
	return &alphabet{runes: map[rune]int{}}
}

func (a *alphabet) add(w word) {
	for _, r := range w {
		i := a.runes[r]
		i++
		a.runes[r] = i
		a.all++
	}
}

func (a alphabet) allRunes() iter.Seq2[rune, float64] {
	runes := slices.Collect(maps.Keys(a.runes))
	slices.Sort(runes)
	return func(visit func(rune, float64) bool) {
		for _, r := range runes {
			if !visit(r, 100.0*float64(a.runes[r])/float64(a.all)) {
				break
			}
		}
	}
}

func (a alphabet) clone() *alphabet {
	clone := &alphabet{
		runes: maps.Clone(a.runes),
		all:   a.all,
	}
	return clone
}

func (a alphabet) match(w word) bool {
	for _, r := range w {
		if _, ok := a.runes[r]; !ok {
			return false
		}
	}
	return true
}

func (a *alphabet) remove(w word) {
	for _, r := range w {
		if i, ok := a.runes[r]; ok {
			delete(a.runes, r)
			a.all -= i
		}
	}
}

func (a alphabet) weight(w word) float64 {
	result := 0
	for _, r := range w {
		result += a.runes[r]
	}
	return 100.0 * float64(result) / float64(a.all)
}

type wordRec struct {
	word
	parent   *wordRec
	children *wordRecs
}

func newWordRec(word word) *wordRec {
	wr := &wordRec{
		word:     word,
		children: newWordRecs(),
	}
	wr.children.parent = wr
	return wr
}

type report = func(wr *wordRec)

func (wr wordRec) depth() int {
	result := 0
	for wr := &wr; wr != nil; wr = wr.parent {
		result++
	}
	return result
}

func (wr *wordRec) fillChildren(words []word, alphabet *alphabet, report report) {
	a1 := alphabet.clone()
	a1.remove(wr.word)

	wr.children.fill(words, a1, report)

	report(wr)
}

func (wr *wordRec) report(alphabet *alphabet) string {
	report, ww := "", 0.0
	for wr := wr; wr != nil; wr = wr.parent {
		w := alphabet.weight(wr.word)
		line := fmt.Sprintf("%s [%6.3f]", wr.word, w)
		if report == "" {
			report = line
		} else {
			report = line + " + " + report
		}
		ww += w
	}
	return fmt.Sprintf("%6.3f : %s", ww, report)
}

type wordRecs struct {
	parent *wordRec
	all    []*wordRec
}

func newWordRecs() *wordRecs {
	return &wordRecs{all: []*wordRec{}}
}

func (ww *wordRecs) add(word word) *wordRec {
	wr := newWordRec(word)
	wr.parent = ww.parent
	ww.all = append(ww.all, wr)
	return wr
}

func (ww *wordRecs) fill(words []word, alphabet *alphabet, report report) {
	for i, word := range words {
		if alphabet.match(word) {
			child := ww.add(word)
			child.fillChildren(words[i:], alphabet, report)
		}
	}
}

type reports struct {
	root  string
	files map[string]*os.File
}

func makeReports(root string) reports {
	return reports{
		root:  root,
		files: map[string]*os.File{},
	}
}

func (r *reports) close() {
	for _, f := range r.files {
		f.Close()
	}
	clear(r.files)
}

func (r *reports) file(path string) *os.File {
	if f, ok := r.files[path]; ok {
		return f
	}
	f, err := os.Create(r.root + path)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	r.files[path] = f
	return f
}

func (r *reports) writeLn(path string, v ...any) {
	f := r.file(path)
	fmt.Fprintln(f, v...)
}

func main() {

	alphabet := newAlphabet()
	words5 := []word{}

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
				alphabet.add(w)
				if !doubles(w) {
					words5 = append(words5, w)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
			panic(err)
		}
	}()

	reports := makeReports("./reports/")
	defer reports.close()

	for r, w := range alphabet.allRunes() {
		reports.writeLn("alphabet.txt", fmt.Sprintf("%c %d %6.3f", r, r, w))
	}

	for _, w := range words5 {
		reports.writeLn("all5.txt", fmt.Sprintf("%s %6.3f", w, alphabet.weight(w)))
	}

	for i, w := range words5 {
		fmt.Println(w)

		wr := newWordRec(w)
		wr.fillChildren(words5[i:], alphabet, func(wr *wordRec) {
			s := wr.report(alphabet)
			reports.writeLn("allRecs.txt", s)
			switch wr.depth() {
			case 4:
				reports.writeLn("allRecs4.txt", s)
			case 5:
				reports.writeLn("allRecs5.txt", s)
			case 6:
				reports.writeLn("allRecs6.txt", s)
			}
		})
	}
}
