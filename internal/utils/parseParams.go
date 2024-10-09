package utils

import "strconv"

func ParseQueryParam(param string, defaultValue int) (int, error) {
	if param == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(param)
}
