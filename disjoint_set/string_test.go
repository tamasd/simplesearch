package disjoint_set_test

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tamasd/simplesearch/disjoint_set"
)

func TestSet(t *testing.T) {
	m := make(map[string]*disjoint_set.StringItem)

	items := []string{"a", "b", "c", "d"}
	for _, v := range items {
		m[v] = disjoint_set.NewStringItem(v)
		require.Equal(t, m[v].String(), v)
	}

	require.NotEqual(t, m["a"].Find(), m["b"].Find())
	require.Equal(t, m["a"].Find(), m["a"])

	disjoint_set.StringUnion(m["a"], m["a"])
	disjoint_set.StringUnion(m["a"], m["b"])

	require.Equal(t, m["a"].Find(), m["b"].Find())
	require.Equal(t, m["a"], m["b"].Find())

	disjoint_set.StringUnion(m["c"], m["a"])
	require.Equal(t, m["a"], m["c"].Find())

	disjoint_set.StringUnion(m["a"], m["d"])
	require.Equal(t, m["a"], m["d"].Find())
}

func BenchmarkStringLinearUnion(b *testing.B) {
	b.ReportAllocs()

	items := make([]*disjoint_set.StringItem, b.N)
	for i := 0; i < b.N; i++ {
		items[i] = disjoint_set.NewStringItem(randomString())
	}

	b.ResetTimer()

	for i := 1; i < b.N; i++ {
		disjoint_set.StringUnion(items[0], items[i])
	}
}

func BenchmarkStringLoopUnion(b *testing.B) {
	b.ReportAllocs()

	items := make([]*disjoint_set.StringItem, b.N)
	for i := 0; i < b.N; i++ {
		items[i] = disjoint_set.NewStringItem(randomString())
	}

	b.ResetTimer()

	for i := 1; i < b.N; i++ {
		disjoint_set.StringUnion(items[i-1], items[i])
	}
}

func BenchmarkStringItem_Find(b *testing.B) {
	b.ReportAllocs()

	items := make([]*disjoint_set.StringItem, b.N)
	for i := 0; i < b.N; i++ {
		items[i] = disjoint_set.NewStringItem(randomString())
		disjoint_set.StringUnion(items[0], items[i])
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		items[i].Find()
	}
}

func randomString() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)

	return hex.EncodeToString(b)
}
