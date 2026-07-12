package rehapt

import (
	"fmt"
	"regexp"
	"strconv"
)

// GetVariable allows retrieving a variable value from its name.
// Nil is returned if the variable is not found
func (r *Rehapt) GetVariable(name string) interface{} {
	return r.variables[name]
}

// LookupVariable allows retrieving a variable value from its name.
func (r *Rehapt) LookupVariable(name string) (interface{}, bool) {
	value, ok := r.variables[name]
	return value, ok
}

// GetVariableString allows retrieving a variable value as a string from its name.
// An empty string is returned if the variable is not found
func (r *Rehapt) GetVariableString(name string) string {
	if value, ok := r.variables[name].(string); ok == true {
		return value
	}
	return ""
}

// LookupVariableString allows retrieving a variable value as a string from its name
func (r *Rehapt) LookupVariableString(name string) (string, bool) {
	if value, ok := r.variables[name].(string); ok == true {
		return value, true
	}
	return "", false
}

// SetVariable allows manually defining a variable.
// Variable names are strings, however values can be any type
func (r *Rehapt) SetVariable(name string, value interface{}) error {
	if r.validVarname(name) == false {
		return fmt.Errorf("invalid variable name %v", name)
	}
	r.variables[name] = value
	return nil
}

// SetStoreShortcutBounds modifies the strings used as prefix and suffix to identify
// a shortcut version of the store variable operation. The default prefix and suffix is "$" which makes
// the default shortcut form like "$myvar$".
func (r *Rehapt) SetStoreShortcutBounds(prefix string, suffix string) error {
	if prefix == "" {
		return fmt.Errorf("invalid prefix, cannot be empty")
	}
	if suffix == "" {
		return fmt.Errorf("invalid suffix, cannot be empty")
	}
	prefixEscaped := regexp.QuoteMeta(prefix)
	suffixEscaped := regexp.QuoteMeta(suffix)
	re, err := regexp.Compile(`^` + prefixEscaped + `([a-zA-Z0-9]+)` + suffixEscaped + `$`)
	if err != nil {
		return err
	}
	r.variableStoreRegexp = re
	return nil
}

// SetLoadShortcutBounds modifies the strings used as prefix and suffix to identify
// a shortcut version of the load variable operation. The default prefix and suffix is "_" which makes
// the default shortcut form like "_myvar_".
func (r *Rehapt) SetLoadShortcutBounds(prefix string, suffix string) error {
	if prefix == "" {
		return fmt.Errorf("invalid prefix, cannot be empty")
	}
	if suffix == "" {
		return fmt.Errorf("invalid suffix, cannot be empty")
	}
	prefixEscaped := regexp.QuoteMeta(prefix)
	suffixEscaped := regexp.QuoteMeta(suffix)
	re, err := regexp.Compile(prefixEscaped + `([a-zA-Z0-9]+)` + suffixEscaped)
	if err != nil {
		return err
	}
	r.variableLoadRegexp = re
	return nil
}

// SetLoadShortcutFloatPrecision changes the precision of float formatting when
// used with a load shortcut. For example "value is _myvar_" can be replaced by
// "value is 10.50" or "value is 10.500000".
func (r *Rehapt) SetLoadShortcutFloatPrecision(precision int) {
	r.floatPrecision = precision
}

func (r *Rehapt) validVarname(name string) bool {
	return r.variableNameRegexp.MatchString(name)
}

func (r *Rehapt) ReplaceVars(str string) string {
	s, _ := r.replaceVars(str)
	return s
}

func (r *Rehapt) replaceVars(str string) (string, error) {
	matches := r.variableLoadRegexp.FindAllStringSubmatchIndex(str, -1)
	if len(matches) == 0 {
		return str, nil
	}

	replaced := make([]byte, 0, len(str)*2)
	offset := 0
	for _, match := range matches {
		// Match should be 4 elements
		// For example, "the _var_ move" should return [4, 9, 5, 8] :
		//  0   45  89
		// "the _var_ move"
		// with the 4 indexes : [prefix start, suffix end, varname start, varname end]
		if len(match) < 4 {
			continue
		}
		prefix := match[0]
		suffix := match[1]
		varnameStart := match[2]
		varnameEnd := match[3]

		// remove the prefix and suffix
		varname := str[varnameStart:varnameEnd]
		value := ""

		// Make sure variable exists, or report error
		ivalue, ok := r.variables[varname]
		if ok == false {
			return "", fmt.Errorf("variable %v is not defined", varname)
		}

		// Try to convert value to string
		switch ival := ivalue.(type) {
		case string:
			value = ival
		case int:
			value = strconv.FormatInt(int64(ival), 10)
		case int8:
			value = strconv.FormatInt(int64(ival), 10)
		case int16:
			value = strconv.FormatInt(int64(ival), 10)
		case int32:
			value = strconv.FormatInt(int64(ival), 10)
		case int64:
			value = strconv.FormatInt(ival, 10)
		case uint:
			value = strconv.FormatUint(uint64(ival), 10)
		case uint8:
			value = strconv.FormatUint(uint64(ival), 10)
		case uint16:
			value = strconv.FormatUint(uint64(ival), 10)
		case uint32:
			value = strconv.FormatUint(uint64(ival), 10)
		case uint64:
			value = strconv.FormatUint(ival, 10)
		case float32:
			value = strconv.FormatFloat(float64(ival), 'f', r.floatPrecision, 32)
		case float64:
			value = strconv.FormatFloat(ival, 'f', r.floatPrecision, 64)
		case bool:
			value = strconv.FormatBool(ival)
		default:
			return "", fmt.Errorf("variable %v of type %T cannot be used inside string", varname, ivalue)
		}

		replaced = append(replaced, str[offset:prefix]...)
		replaced = append(replaced, value...)
		offset = suffix
	}

	// Finish with end of str, if any
	if offset < len(str) {
		replaced = append(replaced, str[offset:]...)
	}

	return string(replaced), nil
}

func (r *Rehapt) storeIfVariable(expected string, actual interface{}) bool {
	elements := r.variableStoreRegexp.FindStringSubmatch(expected)
	if len(elements) > 1 {
		// index 0 is the full match.
		// index 1 is the first group, our variable name without the '_' prefix and suffix
		varname := elements[1]
		// We override any stored value
		r.variables[varname] = actual
		return true
	}
	return false
}
