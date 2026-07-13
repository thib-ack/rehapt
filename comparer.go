package rehapt

import (
	"math"
	"reflect"
	"regexp"
	"time"
)

func NoReplacement(s string) ReplaceFn {
	return func(r *Rehapt) (string, error) {
		return s, nil
	}
}

func TimeDeltaLayout(t time.Time, delta time.Duration, layout string) CompareFn {
	return func(r *Rehapt, ctx compareCtx) []error {
		// TimeDelta can only compare with actual string values
		err := r.requireKind(ctx, reflect.String)
		if err != nil {
			return []error{err}
		}

		// Use specific time format or default one if not specified
		timeFmt := r.defaultTimeDeltaFormat
		if layout != "" {
			timeFmt = layout
		}

		// Parse the actual value given the format
		actualTime, err := time.Parse(timeFmt, ctx.ActualValue.String())
		if err != nil {
			return []error{ctx.Errorf("invalid time. %v", err)}
		}

		dt := t.Sub(actualTime)
		if dt < -delta || dt > delta {
			return []error{ctx.Errorf("max difference between %v and %v allowed is %v, but difference was %v", t, actualTime, delta, dt)}
		}
		return nil
	}
}

// TimeDelta allows comparing a time value with a given +/- delta.
// Delta is compared to math.Abs(expected - actual) which explains why
// if your expected time is time.Now() with a delta of 10sec,
// actual value will match from now-10sec to now+10sec
func TimeDelta(t time.Time, delta time.Duration) CompareFn {
	return TimeDeltaLayout(t, delta, "")
}

// NumberDelta allows comparing a number value with a given +/- delta.
// Delta is compared to math.Abs(expected - actual) which explains why
// if your expected value is 10 with a delta of 3, actual value will match from 7 to 13.
func NumberDelta(value float64, delta float64) CompareFn {
	return func(r *Rehapt, ctx compareCtx) []error {
		if ctx.Actual == nil {
			return []error{ctx.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got <nil>")}
		}
		actualFloatValue := 0.0
		switch ctx.ActualType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			actualFloatValue = float64(ctx.ActualValue.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			actualFloatValue = float64(ctx.ActualValue.Uint())
		case reflect.Float32, reflect.Float64:
			actualFloatValue = ctx.ActualValue.Float()
		default:
			return []error{ctx.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got %v", ctx.ActualType.Kind())}
		}

		dt := math.Abs(value - actualFloatValue)
		if dt > delta {
			return []error{ctx.Errorf("max difference between %v and %v allowed is %v, but difference was %v", value, ctx.Actual, delta, dt)}
		}
		return nil
	}
}

// NumberRange allows comparing a number value within a given [min-max] inclusive range
func NumberRange(min float64, max float64) CompareFn {
	return func(r *Rehapt, ctx compareCtx) []error {
		if min > max {
			return []error{ctx.Errorf("range [%v,%v] is invalid", min, max)}
		}

		if ctx.Actual == nil {
			return []error{ctx.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got <nil>")}
		}
		actualFloatValue := 0.0
		switch ctx.ActualType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			actualFloatValue = float64(ctx.ActualValue.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			actualFloatValue = float64(ctx.ActualValue.Uint())
		case reflect.Float32, reflect.Float64:
			actualFloatValue = ctx.ActualValue.Float()
		default:
			return []error{ctx.Errorf("different kinds. Expected int{8,16,32,64}, uint{8,16,32,64} or float{32,64}, got %v", ctx.ActualType.Kind())}
		}

		if actualFloatValue < min || actualFloatValue > max {
			return []error{ctx.Errorf("value %v is not within the range [%v,%v]", ctx.Actual, min, max)}
		}
		return nil
	}
}

// Regexp allows doing advanced regexp expectations.
// If the regexp is invalid, an error is reported.
// If the actual value to compare with is not a string, an error is reported.
// If the actual value does not match the regexp, an error is reported
// Note that Regexp uses unanchored semantics, which means you have to
// use ^ and $ to make sure you compare the full string
func Regexp(regex string) CompareFn {
	return func(r *Rehapt, ctx compareCtx) []error {
		// Regexp can only compare with actual string values
		err := r.requireKind(ctx, reflect.String)
		if err != nil {
			return []error{err}
		}

		actualStr := ctx.ActualValue.String()

		// Make variable replacement
		regex, err = r.replaceVars(ctx, regex)
		if err != nil {
			return []error{err}
		}

		re, err := regexp.Compile(regex)
		if err != nil {
			return []error{ctx.Errorf("%v", err)}
		}
		if re.MatchString(actualStr) == false {
			return []error{ctx.Errorf("regexp '%v' does not match '%v'", regex, actualStr)}
		}
		return nil
	}
}

// RegexpVars is a mix between Regexp and StoreVar.
// It checks if the actual value matches the regexp,
// but all the groups defined in the regexp can be extracted to variables for later reuse.
// The Vars hold the mapping groupid: varname.
// For example with Regexp: `^Hello (.*) !$` and Vars: map[int]string{0: "all", 1: "name"}
// then if the actual value is "Hello John !", it will match and 2 vars will be stored:
//
//	"all" = "Hello john !"  (group 0 is the full match)
//	"name" = "John"
//
// Note that RegexpVars uses unanchored semantics, which means you have to
// use ^ and $ to make sure you compare the full string
func RegexpVars(regex string, vars map[int]string) CompareFn {
	return func(r *Rehapt, ctx compareCtx) []error {
		// RegexpVars can only compare with actual string values
		err := r.requireKind(ctx, reflect.String)
		if err != nil {
			return []error{err}
		}

		actualStr := ctx.ActualValue.String()

		// Make variable replacement
		regex, err = r.replaceVars(ctx, regex)
		if err != nil {
			return []error{err}
		}

		re, err := regexp.Compile(regex)
		if err != nil {
			return []error{ctx.Errorf("%v", err)}
		}
		elements := re.FindStringSubmatch(actualStr)
		if len(elements) == 0 {
			return []error{ctx.Errorf("regexp '%v' does not match '%v'", regex, actualStr)}
		}

		for groupid, varname := range vars {
			if groupid >= len(elements) {
				return []error{ctx.Errorf("expected variable index %d overflow regexp group count of %d", groupid, len(elements))}
			}
			if err2 := r.setVariable(ctx, varname, elements[groupid]); err2 != nil {
				return []error{err2}
			}
		}
		return nil
	}
}

// StoreVar allows storing the actual value in a variable instead of checking its content
func StoreVar(name string) CompareFn {
	return func(r *Rehapt, ctx compareCtx) []error {
		// Don't compare but store the actual value using the expectedStr as variable name
		if err := r.setVariable(ctx, name, ctx.Actual); err != nil {
			return []error{err}
		}
		return nil
	}
}

// LoadVar allows loading the value of the variable and then comparing it with the actual value
func LoadVar(name string) CompareFn {
	return func(r *Rehapt, ctx compareCtx) []error {
		// Compare actual value with the loaded value (which might be a string or not)
		value := r.GetVariable(name)
		return r.compare(value, ctx.Actual, ctx.Path)
	}
}

// Any allows you to completely ignore the value
func Any() CompareFn {
	return func(r *Rehapt, ctx compareCtx) []error {
		return nil
	}
}

// And allows you to cumulate multiple checks.
// All the comparisons have to be valid to be considered as valid.
// The comparisons are evaluated in order and stop on the first failure
func And(cmp ...interface{}) CompareFn {
	return func(r *Rehapt, ctx compareCtx) []error {
		for _, comparer := range cmp {
			errs := r.compare(comparer, ctx.Actual, ctx.Path)
			if len(errs) > 0 {
				return errs
			}
		}
		return nil
	}
}

// Or allows you to support optional checks.
// Only one of the comparisons has to be valid to be considered as valid.
// The comparisons are evaluated in order and do not stop on the first success
// this way you can use Or(StoreVar("myvar"), ...)
func Or(cmp ...interface{}) CompareFn {
	return func(r *Rehapt, ctx compareCtx) []error {
		if len(cmp) == 0 {
			return nil
		}

		allErrs := []error{}
		failed := 0
		for _, comparer := range cmp {
			errs := r.compare(comparer, ctx.Actual, ctx.Path)
			if len(errs) > 0 {
				failed++
				allErrs = append(allErrs, errs...)
			}
		}
		// Return errors only if all failed
		if failed == len(cmp) {
			return allErrs
		}
		return nil
	}
}

// Not means we don't expect the given value.
// It works as a boolean 'not' operator on the comparison
func Not(value interface{}) CompareFn {
	return func(r *Rehapt, ctx compareCtx) []error {
		// Normal comparison, but error means ok and no error means error
		errs := r.compare(value, ctx.Actual, ctx.Path)
		if len(errs) == 0 {
			return []error{ctx.Errorf("expected not %v, got %v", value, ctx.Actual)}
		}
		return nil
	}
}
