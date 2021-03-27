package internal

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/DiscoOrg/disgo/api"
)

// IDFromToken returns the applicationID from the token
func IDFromToken(token string) (*api.Snowflake, error) {
	strs := strings.Split(token, ".")
	if len(strs) == 0 {
		return nil, errors.New("token is not in a valid format")
	}
	byteID, err := base64.StdEncoding.DecodeString(strs[0])
	if err != nil {
		return nil, err
	}
	strID := api.Snowflake(byteID)
	return &strID, nil
}
