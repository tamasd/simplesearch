// Ported from: https://github.com/nltk/nltk/blob/develop/nltk/stem/porter.py
package stemming

import (
	"bytes"
	"strings"
	"sync"
)

var (
	vowels = map[byte]struct{}{
		'a': {},
		'e': {},
		'i': {},
		'o': {},
		'u': {},
	}

	pool = map[string]string{}

	step1aRules = []rule{
		{"sses", "ss", nil},
		{"ies", "i", nil},
		{"ss", "ss", nil},
		{"s", "", nil},
	}

	measureBuf = sync.Pool{New: func() interface{} {
		buf := bytes.NewBuffer(nil)
		buf.Grow(128)
		return buf
	}}
)

func init() {
	irregularForms := map[string][]string{
		"sky":     {"sky", "skies"},
		"die":     {"dying"},
		"lie":     {"lying"},
		"tie":     {"tying"},
		"news":    {"news"},
		"inning":  {"innings", "inning"},
		"outing":  {"outings", "outing"},
		"canning": {"cannings", "canning"},
		"howe":    {"howe"},
		"proceed": {"proceed"},
		"exceed":  {"exceed"},
		"succeed": {"succeed"},
	}

	for key, line := range irregularForms {
		for _, element := range line {
			pool[element] = key
		}
	}
}

func isConsonant(word string, i int) bool {
	if _, found := vowels[word[i]]; found {
		return false
	}

	if word[i] == 'y' {
		if i == 0 {
			return true
		}

		return !isConsonant(word, i-1)
	}

	return true
}

func measure(stem string) int {
	// Using a buffer here to reduce the allocations
	buf := measureBuf.Get().(*bytes.Buffer)

	for i := range stem {
		if isConsonant(stem, i) {
			buf.WriteByte('c')
		} else {
			buf.WriteByte('v')
		}
	}

	res := strings.Count(buf.String(), `vc`)

	buf.Reset()
	measureBuf.Put(buf)

	return res
}

func hasPositiveMeasure(stem string) bool {
	return measure(stem) > 0
}

func hasGt1Measure(stem string) bool {
	return measure(stem) > 1
}

func containsVowel(stem string) bool {
	for i := range stem {
		if !isConsonant(stem, i) {
			return true
		}
	}

	return false
}

func endsDoubleConsonant(word string) bool {
	l := len(word)
	return l >= 2 && word[l-1] == word[l-2] && isConsonant(word, l-1)
}

func endsCVC(word string) bool {
	l := len(word)
	return (l >= 3 && isConsonant(word, l-3) && !isConsonant(word, l-2) && isConsonant(word, l-1) && word[l-1] != 'w' && word[l-1] != 'x' && word[l-1] != 'y') || (l == 2 && !isConsonant(word, 0) && isConsonant(word, 1))
}

func replaceSuffix(word, suffix, replacement string) string {
	if suffix == "" {
		return word + replacement
	}

	return word[:len(word)-len(suffix)] + replacement
}

type rule struct {
	suffix      string
	replacement string
	condition   func(stem string) bool
}

func applyRuleList(word string, rules []rule) string {
	stem := ""
	for _, rule := range rules {
		if rule.suffix == "*d" && endsDoubleConsonant(word) {
			stem = word[:len(word)-2]
			if rule.condition == nil || rule.condition(stem) {
				return stem + rule.replacement
			} else {
				return word
			}
		}
		if strings.HasSuffix(word, rule.suffix) {
			stem = replaceSuffix(word, rule.suffix, "")
			if rule.condition == nil || rule.condition(stem) {
				return stem + rule.replacement
			} else {
				return word
			}
		}
	}

	return word
}

func step1a(word string) string {
	if strings.HasSuffix(word, "ies") && len(word) == 4 {
		return replaceSuffix(word, "ies", "ie")
	}

	return applyRuleList(word, step1aRules)
}

func step1b(word string) string {
	if strings.HasSuffix(word, "ied") {
		if len(word) == 4 {
			return replaceSuffix(word, "ied", "ie")
		} else {
			return replaceSuffix(word, "ied", "i")
		}
	}

	if strings.HasSuffix(word, "eed") {
		stem := replaceSuffix(word, "eed", "")
		if measure(stem) > 0 {
			return stem + "ee"
		} else {
			return word
		}
	}

	ruleSucceeded := false

	intermediateStem := ""
	for _, suffix := range []string{"ed", "ing"} {
		if strings.HasSuffix(word, suffix) {
			intermediateStem = replaceSuffix(word, suffix, "")
			if containsVowel(intermediateStem) {
				ruleSucceeded = true
				break
			}
		}
	}

	if !ruleSucceeded {
		return word
	}

	lis := len(intermediateStem)
	intermediateStemLast := intermediateStem[lis-1]

	return applyRuleList(intermediateStem, []rule{
		{"at", "ate", nil},
		{"bl", "ble", nil},
		{"iz", "ize", nil},
		{"*d", string(intermediateStemLast), func(_ string) bool {
			return !(intermediateStemLast == 'l' || intermediateStemLast == 's' || intermediateStemLast == 'z')
		}},
		{"", "e", func(stem string) bool {
			return measure(stem) == 1 && endsCVC(stem)
		}},
	})
}

func step1c(word string) string {
	return applyRuleList(word, []rule{
		{"y", "i", step1cRuleCondition},
	})
}

func step1cRuleCondition(stem string) bool {
	return len(stem) > 1 && isConsonant(stem, len(stem)-1)
}

func step2(word string) string {
	if strings.HasSuffix(word, "alli") && hasPositiveMeasure(replaceSuffix(word, "alli", "")) {
		return step2(replaceSuffix(word, "alli", "al"))
	}

	return applyRuleList(word, []rule{
		{"ational", "ate", hasPositiveMeasure},
		{"tional", "tion", hasPositiveMeasure},
		{"enci", "ence", hasPositiveMeasure},
		{"anci", "ance", hasPositiveMeasure},
		{"izer", "ize", hasPositiveMeasure},
		{"bli", "ble", hasPositiveMeasure},
		{"alli", "al", hasPositiveMeasure},
		{"entli", "ent", hasPositiveMeasure},
		{"eli", "e", hasPositiveMeasure},
		{"ousli", "ous", hasPositiveMeasure},
		{"ization", "ize", hasPositiveMeasure},
		{"ation", "ate", hasPositiveMeasure},
		{"ator", "ate", hasPositiveMeasure},
		{"alism", "al", hasPositiveMeasure},
		{"iveness", "ive", hasPositiveMeasure},
		{"fulness", "ful", hasPositiveMeasure},
		{"ousness", "ous", hasPositiveMeasure},
		{"aliti", "al", hasPositiveMeasure},
		{"iviti", "ive", hasPositiveMeasure},
		{"biliti", "ble", hasPositiveMeasure},
		{"fulli", "ful", hasPositiveMeasure},
		{"logi", "log", func(stem string) bool {
			return hasPositiveMeasure(word[:len(word)-3])
		}},
	})
}

func step3(word string) string {
	return applyRuleList(word, []rule{
		{"icate", "ic", hasPositiveMeasure},
		{"ative", "", hasPositiveMeasure},
		{"alize", "al", hasPositiveMeasure},
		{"iciti", "ic", hasPositiveMeasure},
		{"ical", "ic", hasPositiveMeasure},
		{"ful", "", hasPositiveMeasure},
		{"ness", "", hasPositiveMeasure},
	})
}

func step4(word string) string {
	return applyRuleList(word, []rule{
		{"al", "", hasGt1Measure},
		{"ance", "", hasGt1Measure},
		{"ence", "", hasGt1Measure},
		{"er", "", hasGt1Measure},
		{"ic", "", hasGt1Measure},
		{"able", "", hasGt1Measure},
		{"ible", "", hasGt1Measure},
		{"ant", "", hasGt1Measure},
		{"ement", "", hasGt1Measure},
		{"ment", "", hasGt1Measure},
		{"ent", "", hasGt1Measure},
		{"ion", "", step4RuleCondition},
		{"ou", "", hasGt1Measure},
		{"ism", "", hasGt1Measure},
		{"ate", "", hasGt1Measure},
		{"iti", "", hasGt1Measure},
		{"ous", "", hasGt1Measure},
		{"ive", "", hasGt1Measure},
		{"ize", "", hasGt1Measure},
	})
}

func step4RuleCondition(stem string) bool {
	if len(stem) == 0 {
		return false
	}

	last := stem[len(stem)-1]
	return hasGt1Measure(stem) && (last == 's' || last == 't')
}

func step5a(word string) string {
	if strings.HasSuffix(word, "e") {
		stem := replaceSuffix(word, "e", "")
		if hasGt1Measure(stem) {
			return stem
		}
		if measure(stem) == 1 && !endsCVC(stem) {
			return stem
		}
	}

	return word
}

func step5b(word string) string {
	return applyRuleList(word, []rule{
		{"ll", "l", func(stem string) bool {
			return hasGt1Measure(word[:len(word)-1])
		}},
	})
}

func Stem(word string) string {
	stem := strings.ToLower(word)

	if result, found := pool[stem]; found {
		return result
	}

	if len(word) <= 2 {
		return word
	}

	stem = step1a(stem)
	stem = step1b(stem)
	stem = step1c(stem)
	stem = step2(stem)
	stem = step3(stem)
	stem = step4(stem)
	stem = step5a(stem)
	stem = step5b(stem)

	return stem
}
