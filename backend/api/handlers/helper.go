package handlers

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrorTargetNotFound = errors.New("Target not found")
	boolMap             = map[string]bool{
		"true":  true,
		"1":     true,
		"false": false,
		"0":     false,
	}
)

func strListToBool(inpt []string) ([]bool, error) {
	capacity := len(inpt)
	if capacity == 0 {
		return []bool{}, nil
	}

	result := make([]bool, 0, capacity)
	for _, str := range inpt {
		boolean, ok := boolMap[str]
		if !ok {
			return []bool{}, fmt.Errorf("strListToBool: inpt can't be converted to bool")
		}
		result = append(result, boolean)
	}

	return result, nil
}

func strListToInt(inpt []string) ([]int, error) {
	capacity := len(inpt)
	if capacity == 0 {
		return []int{}, nil
	}

	result := make([]int, 0, capacity)
	for _, str := range inpt {
		number, err := strconv.Atoi(str)
		if err != nil {
			return []int{}, fmt.Errorf("strListToInt: inpt can't be converted to int")
		}
		result = append(result, number)
	}

	return result, nil
}

func removeDuplicate[T comparable](inpt []T) []T {
	record := map[T]bool{}
	list := []T{}
	for _, item := range inpt {
		if _, ok := record[item]; !ok {
			record[item] = true
			list = append(list, item)
		}
	}
	return list
}
