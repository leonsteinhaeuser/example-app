package lib

import "net/url"

func URLSchemeOrDefault(url *url.URL, defaultScheme string) string {
	if url == nil {
		return defaultScheme
	}
	// this should not be necessary, but just in case
	// the url is not parsed correctly
	// we return the default scheme
	if url.Scheme == "" {
		return defaultScheme
	}
	return url.Scheme
}
