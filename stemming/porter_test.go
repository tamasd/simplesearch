package stemming_test

import (
	"compress/gzip"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/tamasd/simplesearch/stemming"
)

func TestSimple(t *testing.T) {
	plurals := []string{
		"CARESSES",
		"FLIES",
		"DIES",
		"MULES",
		"DENIED",
		"DIED",
		"AGREED",
		"OWNED",
		"HUMBLED",
		"SIZED",
		"MEETING",
		"STATING",
		"SIEZING",
		"ITEMIZATION",
		"SENSATIONAL",
		"TRADITIONAL",
		"REFERENCE",
		"COLONIZER",
		"PLOTTED",
	}

	singles := make([]string, len(plurals))
	for i := range plurals {
		singles[i] = stemming.Stem(plurals[i])
	}

	stemmed := strings.Join(singles, " ")
	if stemmed != "caress fli die mule deni die agre own humbl size meet state siez item sensat tradit refer colon plot" {
		t.Fatalf("stemmer failed: %s\n", stemmed)
	}
}

func TestRegressions(t *testing.T) {
	regressions := map[string]string{
		"ion": "ion",
	}

	for word, expected := range regressions {
		if stemmed := stemming.Stem(word); stemmed != expected {
			t.Fatalf("stemmer failed: %s should be %s, not %s", word, expected, stemmed)
		}
	}
}

func TestFull(t *testing.T) {
	vocab, output := readData(t)

	if len(vocab) != len(output) {
		t.Fatalf("vocabulary length (%d) does not match output length (%d)\n", len(vocab), len(output))
	}

	checkData(vocab, output, t)
}

type fatalf interface {
	Fatalf(format string, args ...interface{})
}

type errorf interface {
	Errorf(format string, args ...interface{})
}

func readData(t fatalf) ([]string, []string) {
	vocab, err := readDataFile("porter_test_vocabulary.txt.gz")
	if err != nil {
		t.Fatalf("reading vocabulary failed: %v\n", err)
	}

	output, err := readDataFile("porter_test_output.txt.gz")
	if err != nil {
		t.Fatalf("reading output failed: %v\n", err)
	}

	return vocab, output
}

func checkData(vocab, output []string, t errorf) {
	for i := range vocab {
		if stemmed := stemming.Stem(vocab[i]); stemmed != output[i] {
			t.Errorf("stemming failed! %s became %s instead of %s\n", vocab[i], stemmed, output[i])
		}
	}
}

func readDataFile(fn string) ([]string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()

	all, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(all)), "\n"), nil
}

func BenchmarkStem(b *testing.B) {
	b.ReportAllocs()
	vocab, output := readData(b)
	b.ResetTimer()

	for n := min(b.N, len(vocab)); n > 0; n -= len(vocab) {
		checkData(vocab[:n], output[:n], b)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
