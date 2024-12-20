package mask

type MaskFunc func(value string) string

// Mask represents the main masking structure
type Mask struct {
	mask     any
	tokens   Tokens
	reversed bool
}

// NewMask creates a new Mask instance with the given options
func NewMask(mask any, tokens Tokens, reversed bool) *Mask {
	v_tokens := defaultTokens

	for key, token := range tokens {
		v_tokens[key] = token
	}

	return &Mask{
		mask:     mask,
		tokens:   v_tokens,
		reversed: reversed,
	}
}

// Masked returns the masked value
func (m *Mask) Masked(value string) string {
	return m.process(value, m.findMask(value))
}

// Unmasked returns the unmasked value
func (m *Mask) Unmasked(value string) string {
	return m.process(value, m.findMask(value))
}

// findMask determines which mask to use
func (m *Mask) findMask(value string) string {
	switch v := m.mask.(type) {
	case nil:
		return ""
	case string:
		if v == "" {
			return ""
		}
		return v
	case []string:
		if len(v) == 0 {
			return ""
		}
		if len(v) == 1 {
			return v[0]
		}

		// 获取最后���个mask处理后的长度作为参考
		lastMask := v[len(v)-1]
		referenceLen := len(m.process(value, lastMask))

		// 找到第一个处理后长度大于等于参考长度的mask
		for _, mask := range v {
			processedLen := len(m.process(value, mask))
			if processedLen >= referenceLen {
				return mask
			}
		}
		return ""

	case MaskFunc:
		return v(value)
	}
	return ""
}

type escapedMask struct {
	mask    string
	escaped []int
}

func (m *Mask) escapeMask(maskRaw string) escapedMask {
	var chars []rune
	var escaped []int
	runes := []rune(maskRaw)

	for i, ch := range runes {
		if ch == '!' && (i == 0 || runes[i-1] != '!') {
			escaped = append(escaped, i-len(escaped))
		} else {
			chars = append(chars, ch)
		}
	}

	return escapedMask{
		mask:    string(chars),
		escaped: escaped,
	}
}

// process handles the main masking/unmasking logic
func (m *Mask) process(value string, maskRaw string) string {
	if maskRaw == "" {
		return value
	}

	escaped := m.escapeMask(maskRaw)
	var result []rune
	tokens := m.tokens
	offset := 1

	maskRunes := []rune(escaped.mask)
	valueRunes := []rune(value)

	lastMaskChar := 0

	var lastRawMaskChar rune
	repeatedPos := -1
	maskPos := 0
	valuePos := 0

	check := func() bool {
		return maskPos < len(maskRunes) && valuePos < len(valueRunes)
	}

	push := func(r rune) {
		result = append(result, r)
	}

	multipleMatched := false

	for check() {
		maskChar := maskRunes[maskPos]
		token, hasToken := tokens[string(maskChar)]

		var valueChar rune
		if valuePos >= 0 && valuePos < len(valueRunes) {
			if token.Transform != nil {
				valueChar = []rune(token.Transform(string(valueRunes[valuePos])))[0]
			} else {
				valueChar = valueRunes[valuePos]
			}
		}

		// Check if current position is not escaped and has a token
		isEscaped := false
		for _, pos := range escaped.escaped {
			if pos == maskPos {
				isEscaped = true
				break
			}
		}

		if !isEscaped && hasToken {
			// Value symbol matched token
			if token.Pattern != nil && token.Pattern.MatchString(string(valueChar)) {
				push(valueChar)

				if token.Repeated {
					if repeatedPos == -1 {
						repeatedPos = maskPos
					} else if maskPos == lastMaskChar && maskPos != repeatedPos {
						maskPos = repeatedPos - offset
					}

					if lastMaskChar == repeatedPos {
						maskPos -= offset
					}
				} else if token.Multiple {
					multipleMatched = true
					maskPos -= offset
				}

				maskPos += offset
			} else if token.Multiple {
				if multipleMatched {
					maskPos += offset
					valuePos -= offset
					multipleMatched = false
				}
				// Invalid input - skip
			} else if valueChar == lastRawMaskChar {
				// Matched the last raw mask character
				lastRawMaskChar = 0
			} else if token.Optional {
				maskPos += offset
				valuePos -= offset
			}
			// Invalid input - skip

			valuePos += offset
		} else {
			// Mask symbol is literal character
			push(maskChar)

			if valueChar == maskChar {
				valuePos += offset
			} else {
				lastRawMaskChar = maskChar
			}

			maskPos += offset
		}
	}

	return string(result)
}

// Completed checks if the value matches the full mask pattern
func (m *Mask) Completed(value string) bool {
	mask := m.findMask(value)

	if m.mask == nil || mask == "" {
		return false
	}

	processedLen := len(m.process(value, mask))

	switch v := m.mask.(type) {
	case string:
		return processedLen >= len(v)
	default:
		return processedLen >= len(mask)
	}
}
