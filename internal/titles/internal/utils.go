package internal

import (
	"net/url"
	"strings"
)

func BuildOriginURL(offerIDs []string, args []string) string {
	params := url.Values{
		"offerIds":  {strings.Join(offerIDs, ",")},
		"authCode":  {},
		"cmdParams": {url.PathEscape(strings.Join(args, " "))},
	}
	u := url.URL{
		Scheme:   "origin2",
		Path:     "game/launch",
		RawQuery: params.Encode(),
	}
	return u.String()
}
