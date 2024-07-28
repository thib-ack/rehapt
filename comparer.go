package rehapt

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func NoReplacement(s string) ReplaceFn {
	return func(r *Rehapt) (string, error) {
		return s, nil
	}
}

func TimeDeltaLayout(t time.Time, delta time.Duration, layout string) CompareFn {
	return func(r *Rehapt, ctx compareCtx) error {
		// TimeDelta can only compare with actual string values
		if ctx.ActualKind != reflect.String {
			return fmt.Errorf("different kinds. Expected string, got %v", ctx.ActualKind)
		}

		// Use specific time format or default one if not specified
		timeFmt := r.defaultTimeDeltaFormat
		if layout != "" {
			timeFmt = layout
		}

		// Parse the actual value given the format
		actualTime, err := time.Parse(timeFmt, ctx.ActualValue.String())
		if err != nil {
			return fmt.Errorf("invalid time. %v", err)
		}

		dt := t.Sub(actualTime)
		if dt < -delta || dt > delta {
			return fmt.Errorf("max difference between %v and %v allowed is %v, but difference was %v", t, actualTime, delta, dt)
		}
		return nil
	}
}

// TimeDelta allow to compare a time value with a given +/- delta.
// Delta is compared to math.Abs(expected - actual) which explain why
// if your expected time is time.Now() with a delta of 10sec,
// actual value will match from now-10sec to now+10sec
func TimeDelta(t time.Time, delta time.Duration) CompareFn {
	return TimeDeltaLayout(t, delta, "")
}

// NumberDelta allow to compare a number value with a given +/- delta.
// Delta is compared to math.Abs(expected - actual) which explain why
// if your expected value is 10 with a delta of 3, actual value will match from 7 to 13.
func NumberDelta(value float64, delta float64) CompareFn {
	return func(r *Rehapt, ctx compareCtx) error {
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

		dt := math.Abs(value - actualFloatValue)
		if dt > delta {
			return fmt.Errorf("max difference between %v and %v allowed is %v, but difference was %v", value, ctx.Actual, delta, dt)
		}
		return nil
	}
}

// Regexp allow to do advanced regexp expectation.
// If the regexp is invalid, an error is reported.
// If the actual value to compare with is not a string, an error is reported.
// If the actual value does not match the regexp, an error is reported
func Regexp(regex string) CompareFn {
	return func(r *Rehapt, ctx compareCtx) error {
		// Regexp can only compare with actual string values
		if ctx.ActualKind != reflect.String {
			return fmt.Errorf("different kinds. Expected string, got %v", ctx.ActualKind)
		}

		actualStr := ctx.ActualValue.String()

		// Make variable replacement
		var err error
		regex, err = r.replaceVars(regex)
		if err != nil {
			return err
		}

		re, err := regexp.Compile(regex)
		if err != nil {
			return err
		}
		if re.MatchString(actualStr) == false {
			return fmt.Errorf("regexp '%v' does not match '%v'", regex, actualStr)
		}
		return nil
	}
}

// RegexpVars is a mix between Regexp and StoreVar.
// It checks if the actual value matches the regexp.
// but all the groups defined in the regexp can be extracted to variables for later reuse
// The Vars hold the mapping groupid: varname.
// For example with Regexp: `^Hello (.*) !$` and Vars: map[int]string{0: "all", 1: "name"}
// then if the actual value is "Hello john !", it will match and 2 vars will be stored:
//
//	"all" = "Hello john !"  (group 0 is the full match)
//	"name" = "John"
func RegexpVars(regex string, vars map[int]string) CompareFn {
	return func(r *Rehapt, ctx compareCtx) error {
		// RegexpVars can only compare with actual string values
		if ctx.ActualKind != reflect.String {
			return fmt.Errorf("different kinds. Expected string, got %v", ctx.ActualKind)
		}

		actualStr := ctx.ActualValue.String()

		re, err := regexp.Compile(regex)
		if err != nil {
			return err
		}
		elements := re.FindStringSubmatch(actualStr)
		if len(elements) == 0 {
			return fmt.Errorf("regexp '%v' does not match '%v'", regex, actualStr)
		}

		for groupid, varname := range vars {
			if groupid >= len(elements) {
				return fmt.Errorf("expected variable index %d overflow regexp group count of %d", groupid, len(elements))
			}
			if err := r.SetVariable(varname, elements[groupid]); err != nil {
				return err
			}
		}
		return nil
	}
}

// StoreVar allow to store the actual value in a variable instead of checking its content
func StoreVar(name string) CompareFn {
	return func(r *Rehapt, ctx compareCtx) error {
		// Don't compare but store the actual value using the expectedStr as variable name
		if err := r.SetVariable(name, ctx.Actual); err != nil {
			return err
		}
		return nil
	}
}

// LoadVar allow to load the value of the variable and then compare with actual value
func LoadVar(name string) CompareFn {
	return func(r *Rehapt, ctx compareCtx) error {
		// Compare actual value with the loaded value (which might be a string or not)
		value := r.GetVariable(name)
		return r.compare(value, ctx.Actual)
	}
}

// Any allow you to ignore completely the value
func Any() CompareFn {
	return func(r *Rehapt, ctx compareCtx) error {
		return nil
	}
}

func And(cmp ...interface{}) CompareFn {
	return func(r *Rehapt, ctx compareCtx) error {
		for _, comparer := range cmp {
			err := r.compare(comparer, ctx.Actual)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func Or(cmp ...interface{}) CompareFn {
	return func(r *Rehapt, ctx compareCtx) error {
		errs := []string{}
		for _, comparer := range cmp {
			err := r.compare(comparer, ctx.Actual)
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
		// Return errors only if all failed
		if len(errs) == len(cmp) {
			return errors.New(strings.Join(errs, "\n"))
		}
		return nil
	}
}

// Not means we don't expect the given value
// it works as a boolean 'not' operator on the comparison
func Not(value interface{}) CompareFn {
	return func(r *Rehapt, ctx compareCtx) error {
		// Normal comparison, but error means ok and no error means error
		err := r.compare(value, ctx.Actual)
		if err == nil {
			return fmt.Errorf("expected not %v, got %v", value, ctx.Actual)
		}
		return nil
	}
}
