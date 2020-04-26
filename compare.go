package rehapt

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"time"
)

func (r *Rehapt) timeDeltaCompare(ctx compareCtx) error {
	timeDelta, ok := ctx.Expected.(TimeDelta)
	if ok == false {
		// This should never happened because we normally
		// arrive here only when ExpectedType is TimeDelta
		panic("Expected is not TimeDelta")
	}

	// TimeDelta can only compare with actual string values
	if ctx.ActualKind != reflect.String {
		return fmt.Errorf("different kinds. Expected string, got %v", ctx.ActualKind)
	}

	// Use specific time format or default one if not specified
	format := r.defaultTimeDeltaFormat
	if timeDelta.Format != "" {
		format = timeDelta.Format
	}

	// Parse the actual value given the format
	actualTime, err := time.Parse(format, ctx.ActualValue.String())
	if err != nil {
		return fmt.Errorf("invalid time. %v", err)
	}

	dt := timeDelta.Time.Sub(actualTime)
	if dt < -timeDelta.Delta || dt > timeDelta.Delta {
		return fmt.Errorf("max difference between %v and %v allowed is %v, but difference was %v", timeDelta.Time, actualTime, timeDelta.Delta, dt)
	}
	return nil
}

func (r *Rehapt) numberDeltaCompare(ctx compareCtx) error {
	numDelta, ok := ctx.Expected.(NumberDelta)
	if ok == false {
		// This should never happened because we normally
		// arrive here only when ExpectedType is NumberDelta
		panic("Expected is not NumberDelta")
	}

	actualFloatValue := 0.0
	switch ctx.ActualKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		actualFloatValue = float64(ctx.ActualValue.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		actualFloatValue = float64(ctx.ActualValue.Uint())
	case reflect.Float32, reflect.Float64:
		actualFloatValue = ctx.ActualValue.Float()
	default:
		return fmt.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got %v", ctx.ActualKind)
	}

	delta := math.Abs(numDelta.Value - actualFloatValue)
	if delta > numDelta.Delta {
		return fmt.Errorf("max difference between %v and %v allowed is %v, but difference was %v", numDelta.Value, ctx.Actual, numDelta.Delta, delta)
	}
	return nil
}

func (r *Rehapt) regexpVarsCompare(ctx compareCtx) error {
	reVars, ok := ctx.Expected.(RegexpVars)
	if ok == false {
		// This should never happened because we normally
		// arrive here only when ExpectedType is RegexpVars
		panic("Expected is not RegexpVars")
	}

	// RegexpVars can only compare with actual string values
	if ctx.ActualKind != reflect.String {
		return fmt.Errorf("different kinds. Expected string, got %v", ctx.ActualKind)
	}

	actualStr := ctx.ActualValue.String()

	re, err := regexp.Compile(reVars.Regexp)
	if err != nil {
		return err
	}
	elements := re.FindStringSubmatch(actualStr)
	if len(elements) == 0 {
		return fmt.Errorf("regexp '%v' does not match '%v'", reVars.Regexp, actualStr)
	}

	for groupid, varname := range reVars.Vars {
		if groupid >= len(elements) {
			return fmt.Errorf("expected variable index %d overflow regexp group count of %d", groupid, len(elements))
		}
		if err := r.SetVariable(varname, elements[groupid]); err != nil {
			return err
		}
	}
	return nil
}

func (r *Rehapt) unsortedSliceCompare(ctx compareCtx) error {
	if ctx.ActualKind != reflect.Slice {
		return fmt.Errorf("different kinds. Expected slice, got %v", ctx.ActualKind)
	}

	expectedLen := ctx.ExpectedValue.Len()
	actualLen := ctx.ActualValue.Len()
	if expectedLen != actualLen {
		return fmt.Errorf("different slice sizes. Expected %v, got %v. Expected %v got %v", expectedLen, actualLen, ctx.Expected, ctx.Actual)
	}

	// Unordered comparison
	// We build a list of all the indexes (0,1,2,..,N-1)
	// So each time we find a matching element, we can remove its index from this list
	// and ignore it on next search
	actualIndexes := make([]int, actualLen)
	for i := range actualIndexes {
		actualIndexes[i] = i
	}

nextExpected:
	for i := 0; i < expectedLen; i++ {
		expectedElement := ctx.ExpectedValue.Index(i)

		// Now find a matching element in actual object.
		// Once found, ignore the index.
		for j := 0; j < len(actualIndexes); j++ {
			idx := actualIndexes[j]
			actualElement := ctx.ActualValue.Index(idx)

			if err := r.compare(expectedElement.Interface(), actualElement.Interface()); err == nil {
				// Thats a match, ignore this index now, and continue to next expected.
				actualIndexes = append(actualIndexes[:j], actualIndexes[j+1:]...)
				continue nextExpected
			}
		}

		// If we arrive here, we have an expected not matching any actual
		return fmt.Errorf("expected element %v at index %v not found", expectedElement, i)
	}

	// If here we still have actual index, it means unmatched element thats bad
	if len(actualIndexes) > 0 {
		return fmt.Errorf("actual elements at indexes %v not found", actualIndexes)
	}

	return nil
}

func (r *Rehapt) sliceCompare(ctx compareCtx) error {
	if ctx.ActualKind != reflect.Slice {
		return fmt.Errorf("different kinds. Expected slice, got %v", ctx.ActualKind)
	}

	expectedLen := ctx.ExpectedValue.Len()
	actualLen := ctx.ActualValue.Len()
	if expectedLen != actualLen {
		return fmt.Errorf("different slice sizes. Expected %v, got %v. Expected %v got %v", expectedLen, actualLen, ctx.Expected, ctx.Actual)
	}

	// ordered comparison
	for i := 0; i < expectedLen; i++ {
		expectedElement := ctx.ExpectedValue.Index(i)
		actualElement := ctx.ActualValue.Index(i)
		if err := r.compare(expectedElement.Interface(), actualElement.Interface()); err != nil {
			return fmt.Errorf("slice element %v does not match. %v", i, err)
		}
	}

	return nil
}

func (r *Rehapt) partialMapCompare(ctx compareCtx) error {
	if ctx.ActualKind != reflect.Map {
		return fmt.Errorf("different kinds. Expected map, got %v", ctx.ActualKind)
	}

	// Key types have to be the same
	if ctx.ExpectedType.Key() != ctx.ActualType.Key() {
		return fmt.Errorf("different map key types. Expected %v, got %v", ctx.ExpectedType.Key(), ctx.ActualType.Key())
	}

	// Partial match. Ignore the keys not listed in expected map
	// to do this we just have to skip the map size comparison
	keys := ctx.ExpectedValue.MapKeys()
	for _, key := range keys {
		expectedElement := ctx.ExpectedValue.MapIndex(key)
		actualElement := ctx.ActualValue.MapIndex(key)

		if actualElement.IsValid() == false {
			return fmt.Errorf("expected key %v not found", key)
		}

		if err := r.compare(expectedElement.Interface(), actualElement.Interface()); err != nil {
			return fmt.Errorf("map element [%v] does not match. %v", key, err)
		}
	}

	return nil
}

func (r *Rehapt) mapCompare(ctx compareCtx) error {
	if ctx.ActualKind != reflect.Map {
		return fmt.Errorf("different kinds. Expected map, got %v", ctx.ActualKind)
	}

	// Key types have to be the same
	if ctx.ExpectedType.Key() != ctx.ActualType.Key() {
		return fmt.Errorf("different map key types. Expected %v, got %v", ctx.ExpectedType.Key(), ctx.ActualType.Key())
	}

	if ctx.ExpectedValue.Len() != ctx.ActualValue.Len() {
		return fmt.Errorf("different map sizes. Expected %v, got %v", ctx.ExpectedValue.Len(), ctx.ActualValue.Len())
	}

	keys := ctx.ExpectedValue.MapKeys()
	for _, key := range keys {
		expectedElement := ctx.ExpectedValue.MapIndex(key)
		actualElement := ctx.ActualValue.MapIndex(key)

		if actualElement.IsValid() == false {
			return fmt.Errorf("expected key %v not found", key)
		}

		if err := r.compare(expectedElement.Interface(), actualElement.Interface()); err != nil {
			return fmt.Errorf("map element [%v] does not match. %v", key, err)
		}
	}

	return nil
}

func (r *Rehapt) anyCompare(ctx compareCtx) error {
	// We have to ignore this completely as requested by the user
	return nil
}

func (r *Rehapt) storeVarCompare(ctx compareCtx) error {
	expectedStr := ctx.ExpectedValue.String()

	// Don't compare but store the actual value using the expectedStr as variable name
	if err := r.SetVariable(expectedStr, ctx.Actual); err != nil {
		return err
	}
	return nil
}

func (r *Rehapt) loadVarCompare(ctx compareCtx) error {
	expectedStr := ctx.ExpectedValue.String()

	// Compare actual value with the loaded value (which might be a string or not)
	value := r.GetVariable(expectedStr)
	return r.compare(value, ctx.Actual)
}

func (r *Rehapt) regexpCompare(ctx compareCtx) error {
	if ctx.ActualKind != reflect.String {
		return fmt.Errorf("different kinds. Expected string, got %v", ctx.ActualKind)
	}

	expectedStr := ctx.ExpectedValue.String()
	actualStr := ctx.ActualValue.String()

	// Make variable replacement
	var err error
	expectedStr, err = r.replaceVars(expectedStr)
	if err != nil {
		return err
	}

	re, err := regexp.Compile(expectedStr)
	if err != nil {
		return err
	}
	if re.MatchString(actualStr) == false {
		return fmt.Errorf("regexp '%v' does not match '%v'", expectedStr, actualStr)
	}
	return nil
}

func (r *Rehapt) stringCompare(ctx compareCtx) error {
	expectedStr := ctx.ExpectedValue.String()

	// This might be a StoreVar shortcut
	// even if actual value is not a string
	if r.storeIfVariable(expectedStr, ctx.Actual) == true {
		// This was a variable store operation. no comparison to do
		return nil
	}

	if ctx.ActualKind != reflect.String {
		return fmt.Errorf("different kinds. Expected string, got %v", ctx.ActualKind)
	}

	actualStr := ctx.ActualValue.String()

	// Make variable replacement
	var err error
	expectedStr, err = r.replaceVars(expectedStr)
	if err != nil {
		return err
	}

	// classic comparison
	if expectedStr != actualStr {
		return fmt.Errorf("strings does not match. Expected '%v', got '%v'", expectedStr, actualStr)
	}
	return nil
}

func (r *Rehapt) boolCompare(ctx compareCtx) error {
	if ctx.ActualKind != reflect.Bool {
		return fmt.Errorf("different kinds. Expected bool, got %v", ctx.ActualKind)
	}

	expectedBool := ctx.ExpectedValue.Bool()
	actualBool := ctx.ActualValue.Bool()

	// classic comparison
	if expectedBool != actualBool {
		return fmt.Errorf("bools does not match. Expected %v, got %v", expectedBool, actualBool)
	}
	return nil
}

func (r *Rehapt) intCompare(ctx compareCtx) error {
	expectedInt := ctx.ExpectedValue.Int()

	switch ctx.ActualKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		actualInt := ctx.ActualValue.Int()
		// classic comparison
		if expectedInt != actualInt {
			return fmt.Errorf("integers does not match. Expected %v, got %v", expectedInt, actualInt)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		actualInt := ctx.ActualValue.Uint()
		// classic comparison
		if uint64(expectedInt) != actualInt {
			return fmt.Errorf("uintegers does not match. Expected %v, got %v", expectedInt, actualInt)
		}
	case reflect.Float32, reflect.Float64:
		actualFloat := ctx.ActualValue.Float()
		// classic comparison
		if float64(expectedInt) != actualFloat {
			return fmt.Errorf("floats does not match. Expected %v, got %v", expectedInt, actualFloat)
		}
	default:
		return fmt.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got %v", ctx.ActualKind)
	}

	return nil
}

func (r *Rehapt) uintCompare(ctx compareCtx) error {
	expectedInt := ctx.ExpectedValue.Uint()

	switch ctx.ActualKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		actualInt := ctx.ActualValue.Int()
		// classic comparison
		if int64(expectedInt) != actualInt {
			return fmt.Errorf("integers does not match. Expected %v, got %v", expectedInt, actualInt)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		actualInt := ctx.ActualValue.Uint()
		// classic comparison
		if expectedInt != actualInt {
			return fmt.Errorf("uintegers does not match. Expected %v, got %v", expectedInt, actualInt)
		}
	case reflect.Float32, reflect.Float64:
		actualFloat := ctx.ActualValue.Float()
		// classic comparison
		if float64(expectedInt) != actualFloat {
			return fmt.Errorf("floats does not match. Expected %v, got %v", expectedInt, actualFloat)
		}
	default:
		return fmt.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got %v", ctx.ActualKind)
	}

	return nil
}

func (r *Rehapt) floatCompare(ctx compareCtx) error {
	expectedFloat := ctx.ExpectedValue.Float()

	switch ctx.ActualKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		actualInt := ctx.ActualValue.Int()
		// classic comparison
		if int64(expectedFloat) != actualInt {
			return fmt.Errorf("integers does not match. Expected %v, got %v", expectedFloat, actualInt)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		actualInt := ctx.ActualValue.Uint()
		// classic comparison
		if uint64(expectedFloat) != actualInt {
			return fmt.Errorf("uintegers does not match. Expected %v, got %v", expectedFloat, actualInt)
		}
	case reflect.Float32, reflect.Float64:
		actualFloat := ctx.ActualValue.Float()
		// classic comparison
		if expectedFloat != actualFloat {
			return fmt.Errorf("floats does not match. Expected %v, got %v", expectedFloat, actualFloat)
		}
	default:
		return fmt.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got %v", ctx.ActualKind)
	}

	return nil
}
