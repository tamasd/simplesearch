package trie_test

import (
	"crypto/rand"
	mrand "math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tamasd/simplesearch/trie"
)

func TestRadix_Find(t *testing.T) {
	search := []byte("asdf")
	r := trie.NewRadix()
	r.Add(search)
	require.True(t, r.Find(search))

	search2 := []byte("qwer")
	require.False(t, r.Find(search2))
	r.Add(search2)
	require.True(t, r.Find(search2))
}

func TestRadix_Add(t *testing.T) {
	r := trie.NewRadix()

	require.True(t, r.Find(nil))
	require.False(t, r.Find([]byte("zxcv")))

	r.Add(nil)
	require.True(t, r.Find(nil))
}

func benchmark(b *testing.B, min, max, n int) {
	b.ReportAllocs()

	t := trie.NewRadix()
	for i := 0; i < n; i++ {
		t.Add(randomBytes(min, max))
	}

	searches := make([][]byte, b.N)

	for i := 0; i < b.N; i++ {
		searches[i] = randomBytes(min, max)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t.Find(searches[i])
	}
}

func BenchmarkRadix_Find_Shallow(b *testing.B) {
	benchmark(b, 1, 16, 1024)
}

func BenchmarkRadix_Find_Deep(b *testing.B) {
	benchmark(b, 256, 1024, 128)
}

func randomBytes(min, max int) []byte {
	l := mrand.Intn(max-min) + min
	buf := make([]byte, l)

	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return buf
}
