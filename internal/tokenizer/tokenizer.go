package tokenizer

import (
	"bytes"
)

type Key int

type Tokenizer struct {
	tokens  map[Key][][]byte
	overlap bool
}

func New() *Tokenizer {
	return &Tokenizer{
		tokens: make(map[Key][][]byte),
	}
}

func (t *Tokenizer) WithOverlap(overlap bool) *Tokenizer {
	t.overlap = overlap
	return t
}

func (t *Tokenizer) DefineTokens(key Key, tokens ...[]byte) {
	t.tokens[key] = tokens
}

func (t *Tokenizer) DefineTokensString(key Key, tokens ...string) {
	strs := make([][]byte, 0, len(tokens))
	for _, t := range tokens {
		strs = append(strs, []byte(t))
	}
	t.DefineTokens(key, strs...)
}

func (t *Tokenizer) Parse(data []byte) *Stream {
	return &Stream{
		buffer:    data,
		tokenizer: t,
		current:   -1,
	}
}

type Stream struct {
	tokenizer *Tokenizer
	buffer    []byte
	current   Key
	unknown   bool
}

func (s *Stream) Scan() bool {
	for {
		if len(s.buffer) == 0 {
			// return ongoing unknown token
			if s.unknown {
				s.unknown = false
				return true
			}
			return false
		}
		for key, values := range s.tokenizer.tokens {
			for _, value := range values {
				if bytes.HasPrefix(s.buffer, value) {
					// return ongoing unknown token
					if s.unknown {
						s.unknown = false
						return true
					}
					s.current = key
					if s.tokenizer.overlap {
						s.buffer = s.buffer[1:]
					} else {
						s.buffer = s.buffer[len(value):]
					}
					return true
				}
			}
		}
		s.current = -1
		s.unknown = true
		s.buffer = s.buffer[1:]
	}
}

func (s *Stream) Token() Key {
	return s.current
}
