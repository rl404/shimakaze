package utils

import (
	"fmt"
	"strings"
)

// GetKey to generate cache key.
func GetKey(params ...interface{}) string {
	strParams := []string{"shimakaze"}
	for _, p := range params {
		if tmp := fmt.Sprintf("%v", p); tmp != "" {
			strParams = append(strParams, tmp)
		}
	}
	return strings.Join(strParams, ":")
}
