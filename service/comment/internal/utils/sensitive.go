package utils

import (
	"log"
	"strings"

	ahocorasick "github.com/anknown/ahocorasick"
)

type SensitiveFilter struct {
	machine *ahocorasick.Machine
}

func NewSensitiveFilter(words []string) *SensitiveFilter {
	if len(words) == 0 {
		log.Println("Warning!敏感词库为空!")
		return &SensitiveFilter{machine: new(ahocorasick.Machine)}
	}
	m := new(ahocorasick.Machine)
	var runeWords [][]rune
	for _, word := range words {
		lowerWord := strings.ToLower(word)
		runeWords = append(runeWords, []rune(lowerWord))
	}
	if err := m.Build(runeWords); err != nil {
		log.Fatalf("构建AC自动机失败%v", err)
	}
	return &SensitiveFilter{
		machine: m,
	}
}

func (sf *SensitiveFilter) Match(text string) (bool, string) {
	if sf.machine == nil {
		return false, ""
	}
	text = strings.ToLower(text)
	terms := sf.machine.MultiPatternSearch([]rune(text), false)
	if len(terms) > 0 {
		return true, string(terms[0].Word)
	}
	return false, ""
}
