// Package celprocessor provides a simple interface for processing JSON data using CEL expressions
package processCel

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// AttributeConfig represents the minimal interface needed for attribute configuration
type AttributeConfig interface {
	GetExternalId() string
}

// ProcessCELAttributes processes attributes that start with @cel by applying CEL expressions
// to each object in the response and adding the results back to the objects.
// This is the main entry point for adapters to use CEL processing.
func ProcessCELAttributes(attributes []AttributeConfig, objects []map[string]any) error {
	// Find CEL attributes (those starting with @cel)
	var celAttributes []string
	for _, attr := range attributes {
		if strings.HasPrefix(attr.GetExternalId(), "@cel") {
			celAttributes = append(celAttributes, attr.GetExternalId())
		}
	}

	// If no CEL attributes, nothing to do
	if len(celAttributes) == 0 {
		return nil
	}

	// Process each object
	for i, objMap := range objects {
		// Convert object to JSON string for CEL processing
		objJSON, err := json.Marshal(objMap)
		if err != nil {
			return fmt.Errorf("failed to marshal object to JSON: %w", err)
		}

		// Process each CEL attribute
		for _, celAttr := range celAttributes {
			// Strip the @ prefix to get the CEL expression
			celExpression := strings.TrimPrefix(celAttr, "@")

			// Process with CEL
			result, err := ProcessJSON(string(objJSON), celExpression)
			if err != nil {
				return fmt.Errorf("failed to process CEL expression '%s': %w", celExpression, err)
			}

			// Add the result back to the object with the full @cel attribute name
			objMap[celAttr] = result
		}

		// Update the object in the slice (though this isn't strictly necessary since we're modifying the map in place)
		objects[i] = objMap
	}

	return nil
}

// HasCELAttributes checks if any attributes start with @cel (quick check for optimization)
func HasCELAttributes(attributes []AttributeConfig) bool {
	for _, attr := range attributes {
		if strings.HasPrefix(attr.GetExternalId(), "@cel") {
			return true
		}
	}
	return false
}

// ProcessJSON takes a JSON string and CEL expression, returns the result
// This is the main entry point for the library
func ProcessJSON(jsonStr string, celExpression string) (interface{}, error) {
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Create CEL environment with enhanced functions that have access to data
	env, err := cel.NewEnv(
		cel.Lib(&celExtensionsLib{data: data}),
		cel.Variable("data", cel.DynType), // Make data available as a variable
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create CEL environment: %w", err)
	}

	// Compile and evaluate the expression
	ast, issues := env.Compile(celExpression)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("CEL compilation error: %w", issues.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		return nil, fmt.Errorf("CEL program creation error: %w", err)
	}

	out, _, err := prg.Eval(map[string]interface{}{
		"data": data, // Provide data in eval context
	})
	if err != nil {
		// Handle common evaluation errors gracefully
		errMsg := err.Error()
		if strings.Contains(errMsg, "index out of bounds") {
			return "", nil // Return empty string for index out of bounds
		}
		return nil, fmt.Errorf("CEL evaluation error: %w", err)
	}

	// Convert CEL types to standard Go types
	result := convertCELValue(out)
	return result, nil
}

// Enhanced CEL library with cel.* functions
type celExtensionsLib struct {
	data interface{}
}

func (lib *celExtensionsLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		// cel.split(string, delimiter) -> []string
		cel.Function("cel.split",
			cel.Overload("cel_split_string_string_list",
				[]*cel.Type{cel.StringType, cel.StringType},
				cel.ListType(cel.StringType),
				cel.FunctionBinding(celSplitFunc),
			),
		),

		// cel.filterBy(key, value) -> filtered array (works around CEL string literal limitation)
		cel.Function("cel.filterBy",
			cel.Overload("cel_filterBy_string_any_list",
				[]*cel.Type{cel.StringType, cel.DynType},
				cel.ListType(cel.DynType),
				cel.FunctionBinding(lib.celFilterFunc),
			),
		),
		// cel.find(key, value) -> first matching object
		cel.Function("cel.find",
			cel.Overload("cel_find_string_any_any",
				[]*cel.Type{cel.StringType, cel.DynType},
				cel.DynType,
				cel.FunctionBinding(lib.celFindFunc),
			),
		),
		// cel.get(key) -> value (recursive key search)
		cel.Function("cel.get",
			cel.Overload("cel_get_string_any",
				[]*cel.Type{cel.StringType},
				cel.DynType,
				cel.FunctionBinding(lib.celGetFunc),
			),
		),
		// cel.extract(key, delimiter, index) -> string part after splitting
		cel.Function("cel.extract",
			cel.Overload("cel_extract_string_string_int_string",
				[]*cel.Type{cel.StringType, cel.StringType, cel.IntType},
				cel.StringType,
				cel.FunctionBinding(lib.celExtractFunc),
			),
		),
		// cel.path(dotPath) -> value at path
		cel.Function("cel.path",
			cel.Overload("cel_path_string_any",
				[]*cel.Type{cel.StringType},
				cel.DynType,
				cel.FunctionBinding(lib.celPathFunc),
			),
		),
	}
}

func (lib *celExtensionsLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

// cel.split(string, delimiter) - splits a string by delimiter
func celSplitFunc(args ...ref.Val) ref.Val {
	if len(args) != 2 {
		return types.NewErr("cel.split requires exactly 2 arguments")
	}

	// Handle null values by converting to empty string
	var str string
	if args[0].Value() == nil {
		str = ""
	} else {
		str = args[0].Value().(string)
	}

	delimiter := args[1].Value().(string)

	parts := strings.Split(str, delimiter)

	celList := make([]ref.Val, len(parts))
	for i, part := range parts {
		celList[i] = types.String(part)
	}

	return types.NewDynamicList(types.DefaultTypeAdapter, celList)
}

// cel.filter(key, value) - filters array by key/value match
func (lib *celExtensionsLib) celFilterFunc(args ...ref.Val) ref.Val {
	if len(args) != 2 {
		return types.NewErr("cel.filter requires exactly 2 arguments")
	}

	// Use data from the library struct
	dataList, ok := lib.data.([]interface{})
	if !ok {
		return types.NewErr("data is not an array")
	}

	key := args[0].Value().(string)
	searchValue := args[1].Value()

	var filtered []interface{}
	for _, item := range dataList {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if val, exists := itemMap[key]; exists && val == searchValue {
				filtered = append(filtered, item)
			}
		}
	}

	return types.DefaultTypeAdapter.NativeToValue(filtered)
}

// cel.find(key, value) - finds first object matching key/value
func (lib *celExtensionsLib) celFindFunc(args ...ref.Val) ref.Val {
	if len(args) != 2 {
		return types.NewErr("cel.find requires exactly 2 arguments")
	}

	// Use data from the library struct
	dataList, ok := lib.data.([]interface{})
	if !ok {
		return types.NewErr("data is not an array")
	}

	key := args[0].Value().(string)
	searchValue := args[1].Value()

	for _, item := range dataList {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if val, exists := itemMap[key]; exists && val == searchValue {
				return types.DefaultTypeAdapter.NativeToValue(item)
			}
		}
	}

	return types.NewErr("no object found with %s = %v", key, searchValue)
}

// cel.get(key) - recursively searches for a key
func (lib *celExtensionsLib) celGetFunc(args ...ref.Val) ref.Val {
	if len(args) != 1 {
		return types.NewErr("cel.get requires exactly 1 argument")
	}

	key := args[0].Value().(string)

	result := findKeyRecursive(lib.data, key)
	if result != nil {
		return types.DefaultTypeAdapter.NativeToValue(result)
	}

	return types.String("") // Return empty string instead of error
}

// cel.path(dotPath) - access nested values using dot notation
func (lib *celExtensionsLib) celPathFunc(args ...ref.Val) ref.Val {
	if len(args) != 1 {
		return types.NewErr("cel.path requires exactly 1 argument")
	}

	path := args[0].Value().(string)
	parts := strings.Split(path, ".")

	current := lib.data
	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			if val, exists := v[part]; exists {
				current = val
			} else {
				return types.String("") // Return empty string instead of error
			}
		case []interface{}:
			// Handle array index access
			if idx, err := strconv.Atoi(part); err == nil && idx >= 0 && idx < len(v) {
				current = v[idx]
			} else {
				return types.NewErr("invalid array index '%s' in path '%s'", part, path)
			}
		default:
			return types.NewErr("cannot access '%s' in path '%s' - not an object or array", part, path)
		}
	}

	// If current is nil, return empty string instead of null
	if current == nil {
		return types.String("")
	}

	return types.DefaultTypeAdapter.NativeToValue(current)
}

// cel.extract(key, delimiter, index) - finds key, splits value, returns index
func (lib *celExtensionsLib) celExtractFunc(args ...ref.Val) ref.Val {
	if len(args) != 3 {
		return types.NewErr("cel.extract requires exactly 3 arguments")
	}

	key := args[0].Value().(string)
	delimiter := args[1].Value().(string)
	index := int(args[2].Value().(int64))

	// Find the key recursively
	result := findKeyRecursive(lib.data, key)
	if result == nil {
		return types.NewErr("key '%s' not found", key)
	}

	// Convert to string and split
	str, ok := result.(string)
	if !ok {
		return types.NewErr("value for key '%s' is not a string", key)
	}

	parts := strings.Split(str, delimiter)
	if index < 0 || index >= len(parts) {
		return types.NewErr("index %d out of range for split result (length: %d)", index, len(parts))
	}

	return types.String(parts[index])
}

// Helper function for recursive key search
func findKeyRecursive(data interface{}, key string) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		if val, exists := v[key]; exists {
			return val
		}
		for _, value := range v {
			if result := findKeyRecursive(value, key); result != nil {
				return result
			}
		}
	case []interface{}:
		for _, item := range v {
			if result := findKeyRecursive(item, key); result != nil {
				return result
			}
		}
	}
	return nil
}

// convertCELValue converts CEL-specific types to standard Go types
func convertCELValue(val ref.Val) interface{} {
	// Check if it's a CEL list/array
	if lister, ok := val.(traits.Lister); ok {
		// Convert CEL list to Go slice
		size := int(lister.Size().(types.Int))
		result := make([]interface{}, size)

		for i := 0; i < size; i++ {
			item := lister.Get(types.Int(i))
			result[i] = convertCELValue(item)
		}

		return result
	}

	// For non-list values, just return the native Go value
	return val.Value()
}
