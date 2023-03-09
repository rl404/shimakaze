package utils

import (
	"crypto/md5"
	"encoding/json"
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

// QueryToKey to convert query to key string.
func QueryToKey(query interface{}) string {
	j, _ := json.Marshal(query)
	return fmt.Sprintf("%x", md5.Sum(j))
}
