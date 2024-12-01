package regexp

import "regexp"

type Regexp struct {
	*regexp.Regexp
}

func MustNew(str string) *Regexp {
	return &Regexp{
		Regexp: regexp.MustCompile(str),
	}
}

func (re *Regexp) MustFindSubmatch(b []byte) [][]byte {
	sm := re.Regexp.FindSubmatch(b)
	if sm == nil {
		panic("regexp did not match")
	}
	return sm
}
