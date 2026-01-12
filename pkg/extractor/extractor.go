// Copyright 2025 SGNL.ai, Inc.
package extractor

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ValueFromList attempts to extract a value from the provided list based on the provided
// includedPrefix and excludedSuffix.
//
// For example, if you want to extract a URL from a header with the following format:
// `<https://test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2>; rel="next"`
//
// You could specify a `includedPrefix` of `https://` and an `excludedSuffix` of `>; rel="next"`.
// This would return the first value found matching this pattern in the following format:
// `https://test-instance.oktapreview.com/api/v1/users?after=100u65xtp32NovHoPx1d7&limit=2`.
//
// All whitespace is ignored in values, prefixes, and suffixes. If there are multiple matching values
// provided, the value from the first matched will be returned.
func ValueFromList(values []string, includedPrefix, excludedSuffix string) string {
	strippedExcludedSuffix := removeWhiteSpace(excludedSuffix)
	strippedIncludedPrefix := removeWhiteSpace(includedPrefix)

	for _, value := range values {
		value = removeWhiteSpace(value)

		// Check if the suffix is present in the current value. If not, continue.
		excludedSuffixIndex := strings.Index(value, strippedExcludedSuffix)
		if excludedSuffixIndex == -1 {
			continue
		}

		value = value[:excludedSuffixIndex]

		// Check if the prefix is present in the current value, if not continue.
		includedPrefixIndex := strings.LastIndex(value, strippedIncludedPrefix)
		if includedPrefixIndex == -1 {
			continue
		}

		return value[includedPrefixIndex:]
	}

	return ""
}

func removeWhiteSpace(str string) string {
	var b strings.Builder

	b.Grow(len(str))

	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}

	return b.String()
}

// AttributesFromJSONPath extracts a list of all valid attribute names from a JSON path expression based on the spec
// outlined in the article https://goessner.net/articles/JsonPath/, which github.com/PaesslerAG/jsonpath implements.
//
// This function ignores all Union operators, Slice operators, Subscript operators, Wildcards, Filter expressions, and
// Script expressions. For example, for a provided expression of `$.book[?(@.price<10)].authors[1]`, this function would
// return `[]string{"book", "authors"}`. Recursive descents will also be ignored and instead treated as a direct
// descent.
//
// All provided expressions must start with a `$.` prefix.
func AttributesFromJSONPath(expression string) (result []string, err error) {
	if !strings.HasPrefix(expression, "$.") {
		return nil, errors.New("expression missing required '$.' prefix")
	}

	expression = strings.TrimPrefix(expression, "$.")

	var (
		expressionComponents []string
		sb                   strings.Builder
		wrappedStack         []rune
	)

	// Split the provided expression into components using '.' / '[' as a separator (unless
	// '.' / '[' is wrapped by brackets). This is done to prevent splitting filter and script expressions.
	for i, r := range expression {
		if len(wrappedStack) == 0 {
			if r == '.' {
				if sb.String() != "" {
					expressionComponents = append(expressionComponents, sb.String())
					sb.Reset()
				}

				continue
			}

			if r == '[' && i > 0 && expression[i-1] == ']' {
				if sb.String() != "" {
					expressionComponents = append(expressionComponents, sb.String())
					sb.Reset()
				}
			}
		}

		// Note: This stack doesn't check if a set of quotes are both escaped or not,
		// so the following would be valid: `$.[\\"test"]`
		switch r {
		case '\'':
			if len(wrappedStack) > 0 && wrappedStack[len(wrappedStack)-1] == '\'' {
				wrappedStack = wrappedStack[:len(wrappedStack)-1]
			} else {
				wrappedStack = append(wrappedStack, '\'')
			}
		case '"':
			if len(wrappedStack) > 0 && wrappedStack[len(wrappedStack)-1] == '"' {
				wrappedStack = wrappedStack[:len(wrappedStack)-1]
			} else {
				wrappedStack = append(wrappedStack, '"')
			}
		case ']':
			if len(wrappedStack) > 0 && wrappedStack[len(wrappedStack)-1] == '[' {
				wrappedStack = wrappedStack[:len(wrappedStack)-1]
			} else {
				wrappedStack = append(wrappedStack, ']')
			}
		case '[':
			if len(wrappedStack) > 0 && wrappedStack[len(wrappedStack)-1] == ']' {
				wrappedStack = wrappedStack[:len(wrappedStack)-1]
			} else {
				wrappedStack = append(wrappedStack, '[')
			}
		}

		sb.WriteRune(r)
	}

	if len(expressionComponents) == 0 && sb.String() == "" {
		return nil, nil
	}

	expressionComponents = append(expressionComponents, sb.String())

	for _, component := range expressionComponents {
		attribute, err := processJSONPathExpressionComponent(component)
		if err != nil {
			return nil, err
		}

		if attribute != nil {
			result = append(result, *attribute)
		}
	}

	return
}

func processJSONPathExpressionComponent(component string) (*string, error) {
	if component == "" {
		return nil, errors.New("empty expression component provided")
	}

	var attribute string

	switch component[0] {
	case '[':
		// Bracket notation.
		var quotedStack []rune

		// Find the first non-quoted closing bracket.
		for i, r := range component {
			if r == ']' && len(quotedStack) == 0 {
				// Remove outer brackets from the component.
				attribute = component[1:i]

				break
			}

			switch r {
			case '\'':
				if len(quotedStack) > 0 && quotedStack[len(quotedStack)-1] == '\'' {
					quotedStack = quotedStack[:len(quotedStack)-1]
				} else {
					quotedStack = append(quotedStack, '\'')
				}
			case '"':
				if len(quotedStack) > 0 && quotedStack[len(quotedStack)-1] == '"' {
					quotedStack = quotedStack[:len(quotedStack)-1]
				} else {
					quotedStack = append(quotedStack, '"')
				}
			}
		}

		if attribute == "" {
			return nil, errors.New("invalid expression provided: missing closing bracket")
		}

		// First, check if the attribute is quoted with a matching prefix and suffix of `'`, `\'`, `"` or `\"`.
		// If so, remove the quotes. Otherwise, check for unsupported operators.
		switch {
		case strings.HasPrefix(attribute, `'`) && strings.HasSuffix(attribute, `'`):
			attribute = strings.Trim(attribute, `'`)
		case strings.HasPrefix(attribute, `\'`) && strings.HasSuffix(attribute, `\'`):
			attribute = strings.Trim(attribute, `\'`)
		case strings.HasPrefix(attribute, `"`) && strings.HasSuffix(attribute, `"`):
			attribute = strings.Trim(attribute, `"`)
		case strings.HasPrefix(attribute, `\"`) && strings.HasSuffix(attribute, `\"`):
			attribute = strings.Trim(attribute, `\"`)
		case attribute == "*":
			// Unsupported wildcard operator.
			return nil, nil
		case strings.Contains(attribute, ":"):
			// Unsupported array slice operator.
			return nil, nil
		case strings.Contains(attribute, ","):
			// Unsupported union operator.
			return nil, nil
		case strings.HasPrefix(attribute, "?("):
			// Unsupported filter expression.
			return nil, nil
		case strings.HasPrefix(attribute, "("):
			// Unsupported script expression.
			return nil, nil
		default:
			// If the current attribute is an int, this is the subscript operator and
			// should be ignored.
			if _, err := strconv.Atoi(attribute); err == nil {
				// Unsupported subscript operator.
				return nil, nil
			}
		}

		return &attribute, nil
	case '*':
		// Unsupported wildcard.
		return nil, nil
	default:
		// Dot notation.
		attribute = component

		// Ignore any filters or scripts defined in this component.
		attribute = strings.Split(attribute, "[")[0]

		return &attribute, nil
	}
}
