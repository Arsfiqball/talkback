package talkback

const (
	OpIsNull    = "isnull"    // IS NULL
	OpEq        = "eq"        // EQUALS
	OpNe        = "ne"        // NOT EQUALS
	OpGt        = "gt"        // GREATER THAN
	OpLt        = "lt"        // LESS THAN
	OpGte       = "gte"       // GREATER THAN OR EQUALS
	OpLte       = "lte"       // LESS THAN OR EQUALS
	OpContain   = "contain"   // CONTAINS
	OpNcontain  = "ncontain"  // NOT CONTAINS
	OpContains  = "contains"  // CONTAINS CASE SENSITIVE
	OpNcontains = "ncontains" // NOT CONTAINS CASE SENSITIVE
	OpIn        = "in"        // IN
	OpNin       = "nin"       // NOT IN
)

// validOps is a list of valid operations.
var validOps = []string{
	OpIsNull,
	OpEq,
	OpNe,
	OpGt,
	OpLt,
	OpGte,
	OpLte,
	OpContain,
	OpNcontain,
	OpContains,
	OpNcontains,
	OpIn,
	OpNin,
}

// Op is a string representing a valid operation.
type Condition struct {
	Field  string   // Field is the name of the field to filter on.
	Op     string   // Op is the operation to perform.
	Values []string // Values is a list of values to filter on.
}

// Valid returns true if the condition is valid.
func (c Condition) Valid() bool {
	return sliceContainsString(validOps, c.Op) &&
		len(c.Field) > 0 &&
		len(c.Values) > 0
}

// Query is a query to filter on.
type Sort struct {
	Field   string
	Reverse bool
}

// Query is a query to filter on.
type Query struct {
	Conditions  []Condition
	With        []string
	Group       []string
	Accumulator []string
	Sort        []Sort
	Limit       int
	Skip        int
}
