package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"

	"github.com/gin-gonic/gin"
)

type PrefixID string

const (
	ITPrefix    PrefixID = "IT"
	NursePrefix PrefixID = "NS"
)

func UuidGenerator(prefix PrefixID) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	randStr := make([]byte, 30)
	for i := range randStr {
		randStr[i] = chars[rand.Intn(len(chars))]
	}

	return string(prefix) + string(randStr)
}

func JsonBinding(ctx *gin.Context, in interface{}) (string, error) {
	if err := ctx.ShouldBindJSON(in); err != nil {
		var errMsg string
		switch e := err.(type) {
		case *json.SyntaxError:
			errMsg = fmt.Sprintf("Invalid JSON syntax at position %d", e.Offset)
		case *json.UnmarshalTypeError:
			errMsg = fmt.Sprintf("Invalid type for JSON value: expected %s but got %s", e.Type, e.Value)
		default:
			errMsg = "JSON binding error"
		}

		return errMsg, err
	}

	return "", nil
}

func IsValidUrl(in string) bool {
	u, err := url.ParseRequestURI(in)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	re := regexp.MustCompile(`\.[a-zA-Z]+$`)
	return re.MatchString(u.Host)
}
