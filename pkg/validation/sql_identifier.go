// Copyright 2026 SGNL.ai, Inc.

// ABOUTME: Provides SQL identifier validation patterns and functions.
// ABOUTME: Used to prevent SQL injection by validating table, schema, and column names.

package validation

import "regexp"

// ValidSQLIdentifier matches valid SQL identifiers for security purposes.
// SQL identifiers must:
// - Contain only alphanumeric characters, dollar signs, and underscores
// - Be between 1 and 128 characters in length
// This helps prevent SQL injection and ensures compatibility across databases.
var ValidSQLIdentifier = regexp.MustCompile(`^[a-zA-Z0-9$_]{1,128}$`)

// ValidColumnIdentifier is a more permissive regex for column names that allows
// characters which can be safely quoted in DB2 (/, -, space).
// These characters require the identifier to be wrapped in double quotes.
var ValidColumnIdentifier = regexp.MustCompile(`^[a-zA-Z0-9$_/\- ]{1,128}$`)

// IsValidSQLIdentifier checks if a string is a valid SQL identifier.
func IsValidSQLIdentifier(s string) bool {
	return ValidSQLIdentifier.MatchString(s)
}

// IsValidColumnIdentifier checks if a string is a valid column identifier.
// This is more permissive than IsValidSQLIdentifier, allowing /, -, and space.
func IsValidColumnIdentifier(s string) bool {
	return ValidColumnIdentifier.MatchString(s)
}
