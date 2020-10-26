package synonym_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tamasd/simplesearch/stemming"
	"github.com/tamasd/simplesearch/synonym"
)

func TestSynonyms_IngestEmpty(t *testing.T) {
	s := synonym.NewSynonyms(nil)
	s.Ingest()
	require.Equal(t, "", s.Synonym(""))
}

func TestSynonyms_IngestOne(t *testing.T) {
	s := synonym.NewSynonyms(nil)
	s.Ingest("a")
	require.Equal(t, "a", s.Synonym("a"))
}

func TestSynonyms_Synonym(t *testing.T) {
	s := synonym.NewSynonyms(nil)
	s.Ingest("a", "b")
	s.Ingest("b", "c")

	require.Equal(t, s.Synonym("b"), s.Synonym("c"))
	require.True(t, s.Synonyms("b", "c"))
	require.False(t, s.Synonyms("d", "a"))
	require.False(t, s.Synonyms("a", "d"))

	require.True(t, s.Synonyms("f", "f"))
}

func TestSynonyms_Synonyms(t *testing.T) {
	s := synonym.NewSynonyms(synonym.TransformFunc(stemming.Stem))
	err := s.LoadFromDataFile("../data/words.txt.gz")
	require.Nil(t, err)

	table := []struct {
		a string
		b string
		e bool
	}{
		{"adore", "admire", true},
		{"adore", "accomodating", false},
	}

	for _, row := range table {
		require.Equal(t, row.e, s.Synonyms(row.a, row.b))
	}
}
