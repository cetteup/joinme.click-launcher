package internal

import (
	"net"
	"net/url"
	"strconv"
)

const (
	urlQueryKeyMod = "mod"
	portMin        = 1
	portMax        = 65535
)

func QueryHasMod(query url.Values) bool {
	return query != nil && query.Has(urlQueryKeyMod)
}

func GetModFromQuery(query url.Values) string {
	return query.Get(urlQueryKeyMod)
}

func IsValidIPv4(input string) bool {
	ip := net.ParseIP(input)
	if ip == nil {
		return false
	}
	return ip.To4() != nil && ip.IsGlobalUnicast()
}

func IsValidPort(input string) bool {
	portAsInt, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return false
	}
	return portAsInt >= portMin && portAsInt <= portMax
}
