package formatter

import (
	"errors"
	"strconv"
	"strings"
	"time"

	dateparse "github.com/araddon/dateparse"
	valid "github.com/asaskevich/govalidator"
)

//FormatData returns type sepecific value according to value
func FormatData(value string) (string, interface{}, error) {
	if isValidString(value) {
		formattedData, err := StringToString(value)
		return "string", formattedData, err
	} else if isValidInt(value) {
		formattedData, err := StringToInteger(value)
		return "int", formattedData, err
	} else if isValidFloat(value) {
		formattedData, err := StringToFloat(value)
		return "float64", formattedData, err
	} else if isValidBool(value) {
		formattedData, err := StringToBool(value)
		return "bool", formattedData, err
	} else if isValidDateTime(value) {
		formattedData, err := StringToDateTime(value)
		return "time.Time", formattedData, err
	}
	return "unknown_type", nil, errors.New("unknown data type")
}

//FormatConstantDate returns string of standard format constants and error
func FormatConstantDate(s string) (string, error) {
	if dateFormat, err := dateparse.ParseFormat(s); err == nil {
		return dateFormat, nil
	}
	return s, errors.New("cannot parse to format constants")
}

//StringToString returns string corresponding to string value and error
func StringToString(s string) (string, error) {

	if (strings.Contains(string(s[0]), "'") && strings.Contains(string(s[len(s)-1]), "'")) || (strings.Contains(string(s[0]), "\"") && strings.Contains(string(s[len(s)-1]), "\"")) {
		return s, nil
	}
	return "", errors.New("string is not a string")
}

//StringToInteger returns integer value corresponding to string and error
func StringToInteger(s string) (int, error) {
	if s == "" {
		return 0, errors.New("string is not an integer")
	}

	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return int(i), nil
	}

	return 0, errors.New("string is not an integer")
}

//StringToFloat returns float64 value corresponding to string and error
func StringToFloat(s string) (float64, error) {
	if s == "" {
		return 0, errors.New("string is not float")
	}

	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}

	return 0, errors.New("string is not float")
}

//StringToBool returns boolean value corresponding to string and error
func StringToBool(s string) (bool, error) {
	if s == "" {
		return false, errors.New("string is not bool")
	}

	if b, err := strconv.ParseBool(s); err == nil {
		return b, nil
	}

	return false, errors.New("string is not bool")
}

//StringToDateTime returns date time value corresponding to string and error
func StringToDateTime(s string) (time.Time, error) {
	dateFormat, err := FormatConstantDate(s)

	if err != nil {
		return time.Now(), err
	}

	t, _ := time.Parse(dateFormat, s)

	return t, nil
}

func isValidString(s string) bool {
	if (strings.Contains(string(s[0]), "'") && strings.Contains(string(s[len(s)-1]), "'")) || (strings.Contains(string(s[0]), "\"") && strings.Contains(string(s[len(s)-1]), "\"")) {
		return true
	}
	return false
}

func isValidInt(s string) bool {
	return valid.IsInt(s)
}

func isValidFloat(s string) bool {
	return valid.IsFloat(s)
}

func isValidBool(s string) bool {
	if s == "true" || s == "false" {
		return true
	}
	return false
}

func isValidDateTime(s string) bool {
	if _, err := StringToDateTime(s); err == nil {
		return true
	}
	return false
}
