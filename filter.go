// -*- tab-width: 2 -*-

package main

import "strings"

var literalBadStrings = []string{"\"X-PAYPAL-OPERATION-NAME=-\" \"X-PAYPAL-API-RC=-\" \"ORIG-URL=/cgi-bin/webscr\"",
	"\"ORIG-URL=/akamai/sureroute-test-object.html\"",
	"\"url_path\":\"/tealeaftarget\"",
	"\"url_path\":\"/tealeaftarget/\""}

type doubleString struct {
	needsOne string
	needsTwo string
}

var literalDoubleBadStrings = []doubleString{doubleString{"apache_acces",
	"slingshotsiloapi"},
	doubleString{"apache_access",
		"slcsb"}}
var literalAndNotBadStrings = []doubleString{doubleString{"apache_access",
	"Paypal-Debug-Id"},
	doubleString{"apache_access",
		"shotapi"}}
