/*
 * Copyright © 2016 Iury Braun
 * Copyright © 2017 Weyboo
 * 
 * number, string, email, date, username
 */

package sdtp

import (
    "fmt"
    "time"
    "reflect"
    "regexp"
    "strings"
    "net/http"
    "github.com/iurybraun/go-cfg_ini"
    "github.com/iurybraun/i18n_ini"
)

// Name of the struct tag used in examples.
const tagName = "validate"

// Regular expression to validate email address.
var mailRe = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)

// Regular expression to validate username.
//var usernameRe = regexp.MustCompile(`^[a-zA-Z0-9]([.]|[a-zA-Z0-9]){6,20}$`)
var usernameRe = regexp.MustCompile(`^[a-zA-Z0-9]+(?:[.\-_]?[a-zA-Z0-9])*$`)

// Generic data validator.
type Validator interface {
    // Validate method performs validation and returns result and optional error.
    Validate(interface{}, string) (bool, int, error)
}

// DefaultValidator does not perform any validations.
type DefaultValidator struct {
}

func (v DefaultValidator) Validate(val interface{}, lang string) (bool, int, error) {
    return true, 0, nil
}

// StringValidator validates string presence and/or its length.
type StringValidator struct {
    Min int
    Max int
}

func (v StringValidator) Validate(val interface{}, lang string) (bool, int, error) {
    l := len(val.(string))
    
    if l == 0 {
        //return false, fmt.Errorf("cannot be blank")
        return false, Required, fmt.Errorf(i18n_ini.LoadTr(lang, "parameter-errors.Required"))
    }
    
    if l < v.Min {
        //return false, fmt.Errorf("should be at least %v chars long", v.Min)
        return false, TooLow, fmt.Errorf(i18n_ini.LoadTr(lang, "parameter-errors.TooLow"))
    }
    
    if v.Max >= v.Min && l > v.Max {
        //return false, fmt.Errorf("should be less than %v chars long", v.Max)
        return false, LimitExceeded, fmt.Errorf(i18n_ini.LoadTr(lang, "parameter-errors.LimitExceeded"))
    }
    
    return true, 0, nil
}

// NumberValidator performs numerical value validation.
// Its limited to int type for simplicity.
type NumberValidator struct {
    Min int
    Max int
}

func (v NumberValidator) Validate(val interface{}, lang string) (bool, int, error) {
    num := val.(int)
    
    if num < v.Min {
        //return false, TooLow, fmt.Errorf("should be greater than %v", v.Min)
        return false, TooLow, fmt.Errorf(i18n_ini.LoadTr(lang, "parameter-errors.TooLow"))
    }
    
    if v.Max >= v.Min && num > v.Max {
        //return false, LimitExceeded, fmt.Errorf("should be less than %v", v.Max)
        return false, LimitExceeded, fmt.Errorf(i18n_ini.LoadTr(lang, "parameter-errors.LimitExceeded"))
    }
    
    return true, 0, nil
}

// DateValidator checks if date is a valid date.
type DateValidator struct {
}

func (v DateValidator) Validate(val interface{}, lang string) (bool, int, error) {
    
    date := val.(time.Time)
    
    if date.IsZero() {
        return false, Rejected, fmt.Errorf(i18n_ini.LoadTr(lang, "parameter-errors.Rejected"))
    }
    
    return true, 0, nil
    
}

// EmailValidator checks if string is a valid email address.
type EmailValidator struct {
}

func (v EmailValidator) Validate(val interface{}, lang string) (bool, int, error) {
    if !mailRe.MatchString(val.(string)) {
        return false, Rejected, fmt.Errorf(i18n_ini.LoadTr(lang, "parameter-errors.Rejected"))  //fmt.Errorf("is not a valid email address")
    }
    
    return true, 0, nil
}

type UsernameValidator struct {
    Min int
    Max int
}

func (v UsernameValidator) Validate(val interface{}, lang string) (bool, int, error) {
    
    l := len(val.(string))
    
    if l == 0 {
        return false, Required, fmt.Errorf(i18n_ini.LoadTr(lang, "parameter-errors.Required"))
    }
    
    if l < v.Min {
        return false, TooLow, fmt.Errorf(i18n_ini.LoadTr(lang, "parameter-errors.TooLow"))
    }
    
    if v.Max >= v.Min && l > v.Max {
        return false, LimitExceeded, fmt.Errorf(i18n_ini.LoadTr(lang, "parameter-errors.LimitExceeded"))
    }
    
    
    if usernameRe.FindString(val.(string)) == "" {
        return false, Rejected, fmt.Errorf(i18n_ini.LoadTr(lang, "parameter-errors.Rejected"))
    }
    
    return true, 0, nil
}

// Returns validator struct corresponding to validation type
func getValidatorFromTag(tag string) Validator {
    args := strings.Split(tag, ",")
    
    switch args[0] {
        case "number":
            validator := NumberValidator{}
            fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
            return validator
        case "string":
            validator := StringValidator{}
            fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
            return validator
        case "email":
            return EmailValidator{}
        case "date":
            return DateValidator{}
        case "username":
            validator := UsernameValidator{}
            fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
            return validator
    }
    
    return DefaultValidator{}
}

// Performs actual data validation using validator definitions on the struct
func ValidateStructFields(s interface{}, r *http.Request) E1 {  //[]error {
    var lang = cfg_ini.GetLang(r)
    
    err_map := NewErrorData()
    //errs := []error{}
    
    // ValueOf returns a Value representing the run-time data
    v := reflect.ValueOf(s)
    
    for i := 0; i < v.NumField(); i++ {
        // Get the field tag value
        tag := v.Type().Field(i).Tag.Get(tagName)

        // Skip if tag is not defined or ignored
        if tag == "" || tag == "-" {
            continue
        }

        // Get a validator that corresponds to a tag
        validator := getValidatorFromTag(tag)

        // Perform validation
        valid, code_error, err := validator.Validate(v.Field(i).Interface(), lang)

        // Append error to results
        if !valid && err != nil {
            //errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
            err_map = append(err_map, map[string]interface{}{
                            "code": code_error,
                            "location": "field:"+strings.ToLower(v.Type().Field(i).Name),
                            "message": err.Error(),  //i18n_ini.LoadTr(lang, "parameter-errors.Required"),
                        })
        }
    }
    
    if len(err_map) > 0 {
        return err_map  //errs
    } else {
        return nil
    }
}
