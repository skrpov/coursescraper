package main

import (
	"strings"
	"unicode"
)

type TokenType = int

var isLowerCaseCourseCode = map[string]bool{
	"cosc": true,
	"engl": true,
	"math": true,
	"phys": true,
}

const (
	TokenTypeCourseCode = 0
	TokenTypePrereq     = 1
	TokenTypeCorreq     = 2
	TokenTypeAny        = 3
	TokenTypeAll        = 4
)

type Token struct {
	TType TokenType
	Val   string
}

type Tokenizer struct {
	words []string
}

func MakeTokenizer(src string) *Tokenizer {
	return &Tokenizer{
		words: strings.Fields(strings.ToLower(src)),
	}
}

func (t *Tokenizer) parseCourseCode() (Token, bool) {
	if len(t.words) >= 2 {
		if isLowerCaseCourseCode[t.words[0]] {
			code := t.words[1]
			for len(code) > 0 && !unicode.IsDigit(rune(code[len(code)-1])) {
				code = code[:len(code)-1]
			}

			// there was no course code, just a mention of course id
			if len(code) == 0 {
				return Token{}, false
			}

			val := t.words[0] + "_" + code
			t.words = t.words[2:]
			return Token{TType: TokenTypeCourseCode, Val: val}, true
		}
	}

	return Token{}, false
}

func (t *Tokenizer) parsePrereq() (Token, bool) {
	if len(t.words) >= 1 {
		if t.words[0] == "prerequisite:" {
			t.words = t.words[1:]
			return Token{TType: TokenTypePrereq}, true
		}
	}

	return Token{}, false
}

func (t *Tokenizer) parseCorreq() (Token, bool) {
	if len(t.words) >= 1 {
		if t.words[0] == "correquisite:" {
			t.words = t.words[1:]
			return Token{TType: TokenTypePrereq}, true
		}
	}

	return Token{}, false
}

func (t *Tokenizer) parseAny() (Token, bool) {
	if len(t.words) >= 1 {
		if t.words[0] == "one" && t.words[1] == "of" {
			t.words = t.words[1:]
			return Token{TType: TokenTypeAny}, true
		}
	}

	return Token{}, false
}

func (t *Tokenizer) parseAll() (Token, bool) {
	if len(t.words) >= 1 {
		if t.words[0] == "all" && t.words[1] == "of" {
			t.words = t.words[1:]
			return Token{TType: TokenTypeAll}, true
		}
	}

	return Token{}, false
}

func (t *Tokenizer) NextToken() (Token, bool) {
	for len(t.words) > 0 {
		if val, res := t.parseCourseCode(); res {
			return val, true
		}
		if val, res := t.parsePrereq(); res {
			return val, true
		}
		if val, res := t.parseCorreq(); res {
			return val, true
		}
		if val, res := t.parseAny(); res {
			return val, true
		}
		if val, res := t.parseAll(); res {
			return val, true
		}
		t.words = t.words[1:]
	}
	return Token{}, false
}
