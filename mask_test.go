package mask

import (
	"regexp"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNullMask(t *testing.T) {
	mask := NewMask(nil, nil, false)

	if mask.Masked("1a") != "1a" {
		t.Errorf("Expected '1a', got '%s'", mask.Masked("1a"))
	}
}

func TestEmptyStringMask(t *testing.T) {
	mask := NewMask("", nil, false)

	if mask.Masked("1a") != "1a" {
		t.Errorf("Expected '1a', got '%s'", mask.Masked("1a"))
	}
}

func TestUndefinedMask(t *testing.T) {
	mask := NewMask(nil, nil, false)

	if mask.Masked("1a") != "1a" {
		t.Errorf("Expected '1a', got '%s'", mask.Masked("1a"))
	}
}

func TestAtAtMask(t *testing.T) {
	mask := NewMask("@ @", nil, false)

	if mask.Masked("1") != "" {
		t.Errorf("Expected '', got '%s'", mask.Masked("1"))
	}
	if mask.Masked("a") != "a" {
		t.Errorf("Expected 'a', got '%s'", mask.Masked("a"))
	}
	if mask.Masked("ab") != "a b" {
		t.Errorf("Expected 'a b', got '%s'", mask.Masked("ab"))
	}
	if mask.Masked("abc") != "a b" {
		t.Errorf("Expected 'a b', got '%s'", mask.Masked("abc"))
	}
	if mask.Masked("1abc") != "a b" {
		t.Errorf("Expected 'a b', got '%s'", mask.Masked("1abc"))
	}

	if mask.Validate("a") != false {
		t.Errorf("Expected false, got %v", mask.Validate("a"))
	}
	if mask.Validate("ab") != true {
		t.Errorf("Expected true, got %v", mask.Validate("ab"))
	}
}

func TestCustomMask(t *testing.T) {
	tokens := Tokens{
		"D": {
			Pattern:  regexp.MustCompile(`[0-9]`),
			Optional: false,
			Multiple: false,
			Repeated: false,
		},
		"d": {
			Pattern:  regexp.MustCompile(`[0-9]`),
			Optional: true,
			Multiple: false,
			Repeated: false,
		},
		".": {
			Pattern:  regexp.MustCompile(`\.`),
			Optional: false,
			Multiple: false,
			Repeated: false,
		},
	}
	mask := NewMask("dddD.D", tokens, true)

	assert.Equal(t, mask.Masked("1234.5"), "1234.5")
	assert.Equal(t, mask.Masked(".1"), ".1")
	assert.Equal(t, mask.Masked("1.1"), "1.1")
	assert.Equal(t, mask.Masked("12.3"), "12.3")
	assert.NotEqual(t, mask.Masked("12"), "12")
	assert.Equal(t, mask.Masked("123.4"), "123.4")
	assert.Equal(t, mask.Masked("1234.5"), "1234.5")
	assert.Equal(t, mask.Masked("12345.6"), "2345.6")

}

func TestCustomMaskValidate(t *testing.T) {
	tokens := Tokens{
		"D": {
			Pattern:  regexp.MustCompile(`[0-9]`),
			Optional: false,
			Multiple: false,
			Repeated: false,
		},
		"d": {
			Pattern:  regexp.MustCompile(`[0-9]`),
			Optional: true,
			Multiple: false,
			Repeated: false,
		},
		".": {
			Pattern:  regexp.MustCompile(`\.`),
			Optional: false,
			Multiple: false,
			Repeated: false,
		},
	}
	mask := NewMask("dddD.D", tokens, true)

	assert.Equal(t, mask.Validate(".1"), false)
	assert.Equal(t, mask.Validate("1"), false)
	assert.Equal(t, mask.Validate("12"), false)
	assert.Equal(t, mask.Validate(".12"), false)
	assert.Equal(t, mask.Validate("1."), false)
	assert.Equal(t, mask.Validate("1.23"), false)
	assert.Equal(t, mask.Validate("1.1"), true)
	assert.Equal(t, mask.Validate("12.3"), true)
	assert.Equal(t, mask.Validate("123.4"), true)
	assert.Equal(t, mask.Validate("1234.5"), true)
	assert.Equal(t, mask.Validate("12345"), false)
	assert.Equal(t, mask.Validate("12345.6"), false)
}
