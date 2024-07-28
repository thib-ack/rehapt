package rehapt

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

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
	// We build a list of all the indexes (0,1,2,...,N-1)
	// So each time we find a matching element, we can remove its index from this list
	// and ignore it on next search
	actualIndexes := make([]int, actualLen)
	for i := range actualIndexes {
		actualIndexes[i] = i
	}

	var errs []string

nextExpected:
	for i := 0; i < expectedLen; i++ {
		expectedElement := ctx.ExpectedValue.Index(i)

		// Now find a matching element in actual object.
		// Once found, ignore the index.
		for j := 0; j < len(actualIndexes); j++ {
			idx := actualIndexes[j]
			actualElement := ctx.ActualValue.Index(idx)

			if err := r.compare(expectedElement.Interface(), actualElement.Interface()); err == nil {
				// That's a match, ignore this index now, and continue to next expected.
				actualIndexes = append(actualIndexes[:j], actualIndexes[j+1:]...)
				continue nextExpected
			}
		}

		// If we arrive here, we have an expected not matching any actual
		errs = append(errs, fmt.Sprintf("expected element %v at index %v not found", expectedElement, i))
	}

	// If here we still have actual index, it means unmatched element
	if len(actualIndexes) > 0 {
		errs = append(errs, fmt.Sprintf("actual elements at indexes %v not found", actualIndexes))
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
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
		return fmt.Errorf("different slice sizes. Expected %d, got %d. Expected %v got %v", expectedLen, actualLen, ctx.Expected, ctx.Actual)
	}

	var errs []string

	// ordered comparison
	for i := 0; i < expectedLen; i++ {
		expectedElement := ctx.ExpectedValue.Index(i)
		actualElement := ctx.ActualValue.Index(i)
		if err := r.compare(expectedElement.Interface(), actualElement.Interface()); err != nil {
			errs = append(errs, fmt.Sprintf("slice element %v does not match. %v", i, err))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
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

	var errs []string

	// Partial match. Ignore the keys not listed in expected map
	// to do this we just have to skip the map size comparison
	keys := ctx.ExpectedValue.MapKeys()
	for _, key := range keys {
		expectedElement := ctx.ExpectedValue.MapIndex(key)
		actualElement := ctx.ActualValue.MapIndex(key)

		if actualElement.IsValid() == false {
			errs = append(errs, fmt.Sprintf("expected key %v not found", key))
			continue
		}

		if err := r.compare(expectedElement.Interface(), actualElement.Interface()); err != nil {
			errs = append(errs, fmt.Sprintf("map element [%v] does not match. %v", key, err))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
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
		return fmt.Errorf("different map sizes. Expected %d, got %d. Expected %v got %v", ctx.ExpectedValue.Len(), ctx.ActualValue.Len(), ctx.Expected, ctx.Actual)
	}

	var errs []string
	keys := ctx.ExpectedValue.MapKeys()
	for _, key := range keys {
		expectedElement := ctx.ExpectedValue.MapIndex(key)
		actualElement := ctx.ActualValue.MapIndex(key)

		if actualElement.IsValid() == false {
			errs = append(errs, fmt.Sprintf("expected key %v not found in actual %v", key, ctx.Actual))
			continue
		}

		if err := r.compare(expectedElement.Interface(), actualElement.Interface()); err != nil {
			errs = append(errs, fmt.Sprintf("map element [%v] does not match. %v", key, err))
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
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
