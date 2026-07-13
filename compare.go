package rehapt

import (
	"fmt"
	"reflect"
	"time"
)

func (r *Rehapt) requireKind(ctx compareCtx, expectedKind reflect.Kind) error {
	if ctx.Actual == nil {
		return ctx.Errorf("different kinds. Expected %s, got <nil>", expectedKind)
	}
	if ctx.ActualType.Kind() != expectedKind {
		return ctx.Errorf("different kinds. Expected %s, got %v", expectedKind, ctx.ActualType.Kind())
	}
	return nil
}

func (r *Rehapt) unsortedSliceCompare(ctx compareCtx) []error {
	err := r.requireKind(ctx, reflect.Slice)
	if err != nil {
		return []error{err}
	}

	expectedLen := ctx.ExpectedValue.Len()
	actualLen := ctx.ActualValue.Len()
	if expectedLen != actualLen {
		return []error{ctx.Errorf("different slice sizes. Expected length of %v, got %v. Expected %v got %v", expectedLen, actualLen, ctx.Expected, ctx.Actual)}
	}

	// Unordered comparison
	// We build a list of all the indexes (0,1,2,...,N-1)
	// So each time we find a matching element, we can remove its index from this list
	// and ignore it on next search
	actualIndexes := make([]int, actualLen)
	for i := range actualIndexes {
		actualIndexes[i] = i
	}

	var errs []error

nextExpected:
	for i := 0; i < expectedLen; i++ {
		expectedElement := ctx.ExpectedValue.Index(i)

		// Now find a matching element in actual object.
		// Once found, ignore the index.
		for j := 0; j < len(actualIndexes); j++ {
			idx := actualIndexes[j]
			actualElement := ctx.ActualValue.Index(idx)

			if cmpErrs := r.compare(expectedElement.Interface(), actualElement.Interface(), ""); len(cmpErrs) == 0 {
				// That's a match, ignore this index now, and continue to next expected.
				actualIndexes = append(actualIndexes[:j], actualIndexes[j+1:]...)
				continue nextExpected
			}
		}

		// If we arrive here, we have an expected element that doesn't match any actual element
		errs = append(errs, ctx.Errorf("expected element %v at index %v not found", expectedElement, i))
	}

	// If we still have actual indexes here, it means there are unmatched elements
	if len(actualIndexes) > 0 {
		errs = append(errs, ctx.Errorf("actual elements at indexes %v not found", actualIndexes))
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (r *Rehapt) sliceCompare(ctx compareCtx) []error {
	err := r.requireKind(ctx, reflect.Slice)
	if err != nil {
		return []error{err}
	}

	expectedLen := ctx.ExpectedValue.Len()
	actualLen := ctx.ActualValue.Len()
	if expectedLen != actualLen {
		return []error{ctx.Errorf("different slice sizes. Expected length of %d, got %d. Expected %v got %v", expectedLen, actualLen, ctx.Expected, ctx.Actual)}
	}

	var errs []error

	// ordered comparison
	for i := 0; i < expectedLen; i++ {
		expectedElement := ctx.ExpectedValue.Index(i)
		actualElement := ctx.ActualValue.Index(i)
		if cmpErrs := r.compare(expectedElement.Interface(), actualElement.Interface(), fmt.Sprintf("%v[%d]", ctx.Path, i)); len(cmpErrs) > 0 {
			for _, cmpErr := range cmpErrs {
				//errs = append(errs, fmt.Errorf("slice element %v does not match. %v", i, cmpErr))
				//errs = append(errs, fmt.Errorf("[%d]%v", i, cmpErr))
				errs = append(errs, cmpErr)
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (r *Rehapt) partialMapCompare(ctx compareCtx) []error {
	err := r.requireKind(ctx, reflect.Map)
	if err != nil {
		return []error{err}
	}

	// Key types have to be the same
	if ctx.ExpectedType.Key() != ctx.ActualType.Key() {
		return []error{ctx.Errorf("different map key types. Expected %v, got %v", ctx.ExpectedType.Key(), ctx.ActualType.Key())}
	}

	var errs []error

	// Partial match. Ignore the keys not listed in expected map
	// to do this we just have to skip the map size comparison
	keys := ctx.ExpectedValue.MapKeys()
	for _, key := range keys {
		expectedElement := ctx.ExpectedValue.MapIndex(key)
		actualElement := ctx.ActualValue.MapIndex(key)

		if actualElement.IsValid() == false {
			errs = append(errs, ctx.Errorf("expected key %v not found", key))
			continue
		}

		if cmpErrs := r.compare(expectedElement.Interface(), actualElement.Interface(), fmt.Sprintf("%v.%v", ctx.Path, key)); len(cmpErrs) > 0 {
			for _, cmpErr := range cmpErrs {
				//errs = append(errs, fmt.Errorf("map element [%v] does not match. %v", key, cmpErr))
				//errs = append(errs, fmt.Errorf("[%v]%v", key, cmpErr))
				errs = append(errs, cmpErr)
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (r *Rehapt) mapCompare(ctx compareCtx) []error {
	err := r.requireKind(ctx, reflect.Map)
	if err != nil {
		return []error{err}
	}

	var errs []error

	// Key types have to be the same
	if ctx.ExpectedType.Key() != ctx.ActualType.Key() {
		errs = append(errs, ctx.Errorf("different map key types. Expected %v, got %v", ctx.ExpectedType.Key(), ctx.ActualType.Key()))
	}

	if ctx.ExpectedValue.Len() != ctx.ActualValue.Len() {
		errs = append(errs, ctx.Errorf("different map sizes. Expected length of %d, got %d. Expected %v got %v", ctx.ExpectedValue.Len(), ctx.ActualValue.Len(), ctx.Expected, ctx.Actual))
	}

	// Cannot go any further
	if len(errs) > 0 {
		return errs
	}

	keys := ctx.ExpectedValue.MapKeys()
	for _, key := range keys {
		expectedElement := ctx.ExpectedValue.MapIndex(key)
		actualElement := ctx.ActualValue.MapIndex(key)

		if actualElement.IsValid() == false {
			errs = append(errs, ctx.Errorf("expected key %v not found in actual %v", key, ctx.Actual))
			continue
		}

		if cmpErrs := r.compare(expectedElement.Interface(), actualElement.Interface(), fmt.Sprintf("%v.%v", ctx.Path, key)); len(cmpErrs) > 0 {
			for _, cmpErr := range cmpErrs {
				//errs = append(errs, fmt.Errorf("map element [%v] does not match. %v", key, cmpErr))
				//errs = append(errs, fmt.Errorf("[%v]%v", key, cmpErr))
				errs = append(errs, cmpErr)
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (r *Rehapt) stringCompare(ctx compareCtx) []error {
	expectedStr := ctx.ExpectedValue.String()

	// This might be a StoreVar shortcut
	// even if the actual value is not a string
	if r.storeIfVariable(expectedStr, ctx.Actual) == true {
		// This was a variable store operation. no comparison to do
		return nil
	}

	err := r.requireKind(ctx, reflect.String)
	if err != nil {
		return []error{err}
	}

	actualStr := ctx.ActualValue.String()

	// Make variable replacement
	expectedStr, err = r.replaceVars(ctx, expectedStr)
	if err != nil {
		return []error{err}
	}

	// classic comparison
	if expectedStr != actualStr {
		return []error{ctx.Errorf("strings do not match. Expected '%v', got '%v'", expectedStr, actualStr)}
	}
	return nil
}

func (r *Rehapt) boolCompare(ctx compareCtx) []error {
	err := r.requireKind(ctx, reflect.Bool)
	if err != nil {
		return []error{err}
	}

	expectedBool := ctx.ExpectedValue.Bool()
	actualBool := ctx.ActualValue.Bool()

	// classic comparison
	if expectedBool != actualBool {
		return []error{ctx.Errorf("bools do not match. Expected %v, got %v", expectedBool, actualBool)}
	}
	return nil
}

func (r *Rehapt) intCompare(ctx compareCtx) []error {
	expectedInt := ctx.ExpectedValue.Int()

	if ctx.Actual == nil {
		return []error{ctx.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got <nil>")}
	}

	switch ctx.ActualType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		actualInt := ctx.ActualValue.Int()
		// classic comparison
		if expectedInt != actualInt {
			return []error{ctx.Errorf("integers do not match. Expected %v, got %v", expectedInt, actualInt)}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		actualUInt := ctx.ActualValue.Uint()
		// be careful, do not cast a negative expected value to uint
		if expectedInt < 0 || uint64(expectedInt) != actualUInt {
			return []error{ctx.Errorf("uints do not match. Expected %v, got %v", expectedInt, actualUInt)}
		}
	case reflect.Float32, reflect.Float64:
		actualFloat := ctx.ActualValue.Float()
		// classic comparison
		if float64(expectedInt) != actualFloat {
			return []error{ctx.Errorf("floats do not match. Expected %v, got %v", expectedInt, actualFloat)}
		}
	default:
		return []error{ctx.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got %v", ctx.ActualType.Kind())}
	}

	return nil
}

func (r *Rehapt) uintCompare(ctx compareCtx) []error {
	expectedUInt := ctx.ExpectedValue.Uint()

	if ctx.Actual == nil {
		return []error{ctx.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got <nil>")}
	}

	switch ctx.ActualType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		actualInt := ctx.ActualValue.Int()
		// be careful, do not cast a negative actual value to uint
		if actualInt < 0 || expectedUInt != uint64(actualInt) {
			return []error{ctx.Errorf("integers do not match. Expected %v, got %v", expectedUInt, actualInt)}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		actualUInt := ctx.ActualValue.Uint()
		// classic comparison
		if expectedUInt != actualUInt {
			return []error{ctx.Errorf("uints do not match. Expected %v, got %v", expectedUInt, actualUInt)}
		}
	case reflect.Float32, reflect.Float64:
		actualFloat := ctx.ActualValue.Float()
		// classic comparison
		if float64(expectedUInt) != actualFloat {
			return []error{ctx.Errorf("floats do not match. Expected %v, got %v", expectedUInt, actualFloat)}
		}
	default:
		return []error{ctx.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got %v", ctx.ActualType.Kind())}
	}

	return nil
}

func (r *Rehapt) floatCompare(ctx compareCtx) []error {
	expectedFloat := ctx.ExpectedValue.Float()

	if ctx.Actual == nil {
		return []error{ctx.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got <nil>")}
	}

	switch ctx.ActualType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		actualInt := ctx.ActualValue.Int()
		// classic comparison
		if expectedFloat != float64(actualInt) {
			return []error{ctx.Errorf("integers do not match. Expected %v, got %v", expectedFloat, actualInt)}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		actualInt := ctx.ActualValue.Uint()
		// classic comparison
		if expectedFloat != float64(actualInt) {
			return []error{ctx.Errorf("uints do not match. Expected %v, got %v", expectedFloat, actualInt)}
		}
	case reflect.Float32, reflect.Float64:
		actualFloat := ctx.ActualValue.Float()
		// classic comparison
		if expectedFloat != actualFloat {
			return []error{ctx.Errorf("floats do not match. Expected %v, got %v", expectedFloat, actualFloat)}
		}
	default:
		return []error{ctx.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got %v", ctx.ActualType.Kind())}
	}

	return nil
}

func (r *Rehapt) timeCompare(ctx compareCtx) []error {
	expectedTime := ctx.ExpectedValue.Interface().(time.Time)
	fn := TimeDelta(expectedTime, 0)
	return fn(r, ctx)
}
