package main

func main() {
	deprecated()             // ISSUE
	deprecatedAndCommented() // ISSUE
	notQuiteDeprecated()
	badlyDeprecated()
	loudlyDeprecated()
	notDeprecated()
}

// Deprecated: some reason
func deprecated() {}

// This is an empty function
//
// Deprecated: some reason
func deprecatedAndCommented() {}

/// Deprecated: some reason
func notQuiteDeprecated() {}

// This type is deprecated for some reason
func badlyDeprecated() {}

// DEPRECATED: some reason
func loudlyDeprecated() {}

// This function is not deprecated
func notDeprecated() {}
