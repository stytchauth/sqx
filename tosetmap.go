package sqx

import (
	"fmt"
	"reflect"
	"strings"

	scan "github.com/blockloop/scan/v2"
)

const (
	dbTag                   = "db"
	sqxTag                  = "sqx"
	excludeOnInsertTagValue = "excludeOnInsert"
)

// ToSetMap converts a struct into a map[string]any based on the presence of "db" struct tags
// Nil values are skipped over automatically
// Add fields to the "excluded" arg to exclude them from the row
func ToSetMap(v any, excluded ...string) (map[string]any, error) {
	if isNil(v) {
		return map[string]any{}, nil
	}

	model, err := reflectValue(v)
	if err != nil {
		return nil, fmt.Errorf("ToSetMap: %w", err)
	}

	newExclusions := excludedTags(model)
	excluded = append(excluded, newExclusions...)
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

func excludedTags(model reflect.Value) []string {
	numfield := model.NumField()
	res := []string{}

	for i := 0; i < numfield; i++ {
		valField := model.Field(i)
		if !valField.IsValid() || !valField.CanSet() {
			continue
		}

		typeField := model.Type().Field(i)

		dbTag, hasDBTag := typeField.Tag.Lookup(dbTag)
		if !hasDBTag || dbTag == "-" {
			continue
		}

		sqxTag, hasSQXTag := typeField.Tag.Lookup(sqxTag)
		if !hasSQXTag || sqxTag == "-" {
			continue
		}

		tagValues := strings.Split(sqxTag, ",")
		for _, tagValue := range tagValues {
			if tagValue == excludeOnInsertTagValue {
				res = append(res, dbTag)
				break
			}
		}
	}

	return res
}

// Same logic as blockloop.scan for extracting the reflect.Value from a pointer.
func reflectValue(v interface{}) (reflect.Value, error) {
	vType := reflect.TypeOf(v)
	vKind := vType.Kind()
	if vKind != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("%q must be a pointer: %w", vKind.String(), scan.ErrNotAPointer)
	}

	vVal := reflect.Indirect(reflect.ValueOf(v))
	if vVal.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("%q must be a pointer to a struct: %w", vKind.String(), scan.ErrNotAStructPointer)
	}
	return vVal, nil
}
