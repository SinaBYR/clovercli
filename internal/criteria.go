package internal

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	q "github.com/ostafen/clover/v2/query"
)

// parseFilter converts a parsed JSON5 filter object into a clover criteria.
//
// A filter is a mongo-like query document:
//
//	{ author: "Christopher Nolan", year: { $gt: 2010 } }
//	{ $and: [ { rating: { $gte: 8.5 } }, { year: { $lt: 2000 } } ] }
//
// Multiple keys of the same object are combined with a logical AND, a list of
// operators on the same field is combined with a logical AND as well.
//
// It returns a nil criteria when the filter is empty, meaning "match every
// document", along with an error for any unrecognized or malformed directive.
func parseFilter(filter map[string]interface{}) (q.Criteria, error) {
	if len(filter) == 0 {
		return nil, nil
	}

	var criteria []q.Criteria
	for _, key := range sortedKeys(filter) {
		if strings.HasPrefix(key, "$") {
			c, err := logicalCriteria(key, filter[key])
			if err != nil {
				return nil, err
			}
			criteria = append(criteria, c)
			continue
		}

		c, err := fieldCriteria(key, filter[key])
		if err != nil {
			return nil, err
		}
		criteria = append(criteria, c)
	}

	return andAll(criteria), nil
}

// fieldCriteria builds the criteria for a single field, which is either a
// plain value (implicit equality) or a list of operators, e.g. { $gt: 1, $lt: 5 }.
func fieldCriteria(field string, value interface{}) (q.Criteria, error) {
	conditions, ok := value.(map[string]interface{})
	if !ok {
		return q.Field(field).Eq(value), nil
	}

	if len(conditions) == 0 {
		return nil, fmt.Errorf("field %q: empty operator list", field)
	}

	var criteria []q.Criteria
	for _, op := range sortedKeys(conditions) {
		c, err := operatorCriteria(field, op, conditions[op])
		if err != nil {
			return nil, err
		}
		criteria = append(criteria, c)
	}

	return andAll(criteria), nil
}

// operatorCriteria translates a single query operator into a clover criteria.
func operatorCriteria(field string, op string, value interface{}) (q.Criteria, error) {
	switch op {
	case "$eq":
		return q.Field(field).Eq(value), nil

	case "$ne":
		return q.Field(field).Neq(value), nil

	case "$gt":
		return q.Field(field).Gt(value), nil

	case "$gte":
		return q.Field(field).GtEq(value), nil

	case "$lt":
		return q.Field(field).Lt(value), nil

	case "$lte":
		return q.Field(field).LtEq(value), nil

	case "$in", "$nin":
		values, ok := value.([]interface{})
		if !ok {
			return nil, fmt.Errorf("%s must be an array", op)
		}

		c := q.Field(field).In(values...)
		if op == "$nin" {
			c = c.Not()
		}
		return c, nil

	case "$exists":
		exists, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("$exists must be a boolean")
		}

		if exists {
			return q.Field(field).Exists(), nil
		}
		return q.Field(field).NotExists(), nil

	// $regex is the mongo name, $like the clover one
	case "$regex", "$like":
		pattern, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%s must be a string", op)
		}

		if _, err := regexp.Compile(pattern); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return q.Field(field).Like(pattern), nil

	case "$not":
		conditions, ok := value.(map[string]interface{})
		if !ok || len(conditions) == 0 {
			return nil, fmt.Errorf("$not must be a non-empty operator list")
		}

		c, err := fieldCriteria(field, value)
		if err != nil {
			return nil, err
		}
		return c.Not(), nil

	default:
		return nil, fmt.Errorf("unknown directive %q", op)
	}
}

// logicalCriteria translates the query operators which act on a whole
// sub-filter instead of a single field: $and, $or, $nor and $not.
func logicalCriteria(op string, value interface{}) (q.Criteria, error) {
	switch op {
	case "$and", "$or", "$nor":
		filters, ok := value.([]interface{})
		if !ok {
			return nil, fmt.Errorf("%s must be an array of filters", op)
		}
		if len(filters) == 0 {
			return nil, fmt.Errorf("%s must contain at least one filter", op)
		}

		var criteria []q.Criteria
		for i, filter := range filters {
			m, ok := filter.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("%s: filter at position %d must be an object", op, i)
			}

			c, err := parseFilter(m)
			if err != nil {
				return nil, fmt.Errorf("%s: filter at position %d: %w", op, i, err)
			}
			criteria = append(criteria, c)
		}

		join := andAll
		if op == "$or" || op == "$nor" {
			join = orAll
		}

		joined := join(criteria)
		if op == "$nor" {
			return joined.Not(), nil
		}
		return joined, nil

	case "$not":
		filter, ok := value.(map[string]interface{})
		if !ok || len(filter) == 0 {
			return nil, fmt.Errorf("$not must be a non-empty filter object")
		}

		c, err := parseFilter(filter)
		if err != nil {
			return nil, err
		}
		if c == nil {
			return nil, fmt.Errorf("$not must contain at least one condition")
		}
		return c.Not(), nil

	default:
		return nil, fmt.Errorf("unknown directive %q", op)
	}
}

// andAll chains the given criteria with a logical AND. It returns nil when no
// criteria is passed, so that the query selects every document.
func andAll(criteria []q.Criteria) q.Criteria {
	if len(criteria) == 0 {
		return nil
	}

	result := criteria[0]
	for _, c := range criteria[1:] {
		result = result.And(c)
	}
	return result
}

// orAll chains the given criteria with a logical OR.
func orAll(criteria []q.Criteria) q.Criteria {
	if len(criteria) == 0 {
		return nil
	}

	result := criteria[0]
	for _, c := range criteria[1:] {
		result = result.Or(c)
	}
	return result
}

// sortedKeys returns the keys of the given map in a stable order, so that
// errors and queries don't depend on the map iteration order.
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
