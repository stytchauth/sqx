package sqx

import (
	"reflect"

	scan "github.com/blockloop/scan"
)

// ToSetMap converts a struct into a map[string]any based on the presence of "db" struct tags
// Nil values are skipped over automatically
// Add fields to the "excluded" arg to exclude them from the row
func ToSetMap(v any, excluded ...string) (map[string]any, error) {
	if isNil(v) {
		return map[string]any{}, nil
	}
	cols, err := scan.ColumnsStrict(v, excluded...)
	if err != nil {
		return nil, err
	}
	vals, err := scan.Values(cols, v)
	if err != nil {
		return nil, err
	}
	setMap := make(map[string]any, len(cols))
	for i := range cols {
		if !isNil(vals[i]) {
			setMap[cols[i]] = vals[i]
		}
	}
	return setMap, nil
}

// isNil determines if an interface - which may be a ptr - is nil, ignoring the ptr type
// see https://stackoverflow.com/questions/13476349/check-for-nil-and-nil-interface-in-go
func isNil(v any) bool {
	if v == nil {
		return true
	}
	return reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil()
}

// ToSetMapAlias is like ToSetMap, but takes in a table alias
func ToSetMapAlias(tableName string, v any, excluded ...string) (map[string]any, error) {
	setMap, err := ToSetMap(v, excluded...)
	if err != nil {
		return nil, err
	}
	aliasedSetMap := make(map[string]any, len(setMap))
	for key, value := range setMap {
		aliasedSetMap[tableName+"."+key] = value
	}
	return aliasedSetMap, nil
}
