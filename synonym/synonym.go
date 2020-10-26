package synonym

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/tamasd/simplesearch/disjoint_set"
)

type Transformer interface {
	Transform(s string) string
}

type TransformFunc func(s string) string

func (f TransformFunc) Transform(s string) string {
	return f(s)
}

type Synonyms struct {
	words       map[string]*disjoint_set.StringItem
	transformer Transformer
}

func NewSynonyms(t Transformer) *Synonyms {
	return &Synonyms{
		words:       map[string]*disjoint_set.StringItem{},
		transformer: t,
	}
}

func (s *Synonyms) maybeTransform(word string) string {
	if s.transformer != nil {
		return s.transformer.Transform(word)
	}

	return word
}

func (s *Synonyms) Ingest(list ...string) {
	if len(list) == 0 {
		return
	}

	root := s.ensure(list[0])
	for _, word := range list[1:] {
		disjoint_set.StringUnion(root, s.ensure(word))
	}
}

func (s *Synonyms) ensure(word string) *disjoint_set.StringItem {
	word = s.maybeTransform(word)

	if item, ok := s.words[word]; ok {
		return item
	}

	item := disjoint_set.NewStringItem(word)
	s.words[word] = item

	return item
}

func (s *Synonyms) Synonym(word string) string {
	word = s.maybeTransform(word)
	if item, found := s.words[word]; found {
		return item.Find().String()
	}

	return word
}

func (s *Synonyms) Synonyms(word0, word1 string) bool {
	word0 = s.maybeTransform(word0)
	word1 = s.maybeTransform(word1)

	if word0 == word1 {
		return true
	}

	item0, found := s.words[word0]
	if !found {
		return false
	}

	item1, found := s.words[word1]
	if !found {
		return false
	}

	return item0.Find() == item1.Find()
}

func (s *Synonyms) LoadFromDataFile(fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	var r io.Reader = f

	if strings.HasSuffix(fn, ".gz") {
		if r, err = gzip.NewReader(r); err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		words := make([]string, 0, len(line))
		for _, phrase := range strings.Split(line, ",") {
			phrase = strings.TrimSpace(phrase)
			if strings.IndexFunc(phrase, unicode.IsSpace) < 0 {
				words = append(words, phrase)
			}
		}

		s.Ingest(words...)
	}
	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}
