package proxy

import "net/url"

var defaultUrl, _ = url.Parse("https://github.com")
var Default Proxy = Proxy{
	Target: defaultUrl,
	Auth:   nil,
}
