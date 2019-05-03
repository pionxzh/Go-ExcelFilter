package main

import (
    "excelFilter/util"
	"strings"
)

// FilterByListRule will filter the array by list of rule
func FilterByListRule(arr, rules []string) []string {
    // define the filter function
	filterFn := func (str string) bool {
		for _, rule := range rules {
			if strings.Contains(str, rule) {
				return true
			}
		}
		return false
	}

    result := util.Filter(arr, filterFn)
	return result
}
