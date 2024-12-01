package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ProcessLine(parser *Tokenizer, data []byte) ([]Key, error) {
	stream := parser.Parse(data)
	var tokens []Key
	for stream.Scan() {
		tokens = append(tokens, stream.Token())
	}
	return tokens, nil
}

func TestTokenizer_nooverlap(t *testing.T) {
	p := New().WithOverlap(false)
	p.DefineTokensString(0, "0", "zero")
	p.DefineTokensString(1, "1", "one")
	p.DefineTokensString(2, "2", "two")
	p.DefineTokensString(3, "3", "three")
	p.DefineTokensString(4, "4", "four")
	p.DefineTokensString(5, "5", "five")
	p.DefineTokensString(6, "6", "six")
	p.DefineTokensString(7, "7", "seven")
	p.DefineTokensString(8, "8", "eight")
	p.DefineTokensString(9, "9", "nine")

	for _, tc := range []struct {
		input []byte
		value []Key
	}{
		{[]byte("four9tbnqhjlbmqnjq4gpzpvjtl2"), []Key{4, 9, -1, 4, -1, 2}},
		{[]byte("8three75sevenbbsbxjscvseven6mhpx"), []Key{8, 3, 7, 5, 7, -1, 7, 6, -1}},
		{[]byte("fivetmxkjczpjninefive5pss3onetwonetmq"), []Key{5, -1, 9, 5, 5, -1, 3, 1, 2, -1}},
		{[]byte("testfive5twonexx"), []Key{-1, 5, 5, 2, -1}},
	} {
		t.Run(string(tc.input), func(t *testing.T) {
			v, err := ProcessLine(p, tc.input)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.value, v)
			}
		})
	}
}

func TestTokenizer_overlap(t *testing.T) {
	p := New().WithOverlap(true)
	p.DefineTokensString(0, "0", "zero")
	p.DefineTokensString(1, "1", "one")
	p.DefineTokensString(2, "2", "two")
	p.DefineTokensString(3, "3", "three")
	p.DefineTokensString(4, "4", "four")
	p.DefineTokensString(5, "5", "five")
	p.DefineTokensString(6, "6", "six")
	p.DefineTokensString(7, "7", "seven")
	p.DefineTokensString(8, "8", "eight")
	p.DefineTokensString(9, "9", "nine")

	for _, tc := range []struct {
		input []byte
		value []Key
	}{
		{[]byte("four9tbnqhjlbmqnjq4gpzpvjtl2"), []Key{4, -1, 9, -1, 4, -1, 2}},
		{[]byte("8three75sevenbbsbxjscvseven6mhpx"), []Key{8, 3, -1, 7, 5, 7, -1, 7, -1, 6, -1}},
		{[]byte("fivetmxkjczpjninefive5pss3onetwonetmq"), []Key{5, -1, 9, -1, 5, -1, 5, -1, 3, 1, -1, 2, -1, 1, -1}},
		{[]byte("testfive5twonexx"), []Key{-1, 5, -1, 5, 2, -1, 1, -1}},
	} {
		t.Run(string(tc.input), func(t *testing.T) {
			v, err := ProcessLine(p, tc.input)
			if assert.NoError(t, err) {
				assert.Equal(t, tc.value, v)
			}
		})
	}
}
