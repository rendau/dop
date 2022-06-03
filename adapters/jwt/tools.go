package jwt

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/rendau/dop/dopErrs"
)

func ParsePayload(token string, dst interface{}) error {
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) == 3 {
		if claimsRaw, err := base64.RawURLEncoding.DecodeString(tokenParts[1]); err == nil {
			if json.Unmarshal(claimsRaw, dst) == nil {
				return nil
			}
		}
	}

	return dopErrs.BadJwt
}