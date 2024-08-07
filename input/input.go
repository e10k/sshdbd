package input

import (
	"fmt"
	"strings"
)

func ParseInput(s string) (string, string, []string, error) {
	s = strings.Trim(s, " ")
	slice := strings.Split(s, ":")

	// a valid input will have one of these forms:
	//    main:sakila (connection id + database name)
	//    main:sakila:table_1,table_2,table_3 (connection id + database name + comma separated list of table names)
	l := len(slice)
	if l == 3 {
		var tables []string
		for _, t := range strings.Split(slice[2], ",") {
			t = strings.Trim(t, " ")
			if len(t) > 0 {
				tables = append(tables, t)
			}
		}
		return slice[0], slice[1], tables, nil
	} else if l == 2 {
		return slice[0], slice[1], nil, nil
	}

	return "", "", nil, fmt.Errorf("unexpected input data format: %s", s)
}
