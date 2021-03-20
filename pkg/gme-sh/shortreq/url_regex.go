package shortreq

import (
	"log"
	"regexp"
)

var UrlRegex *regexp.Regexp

func init() {
	var err error
	UrlRegex, err = regexp.Compile(`^(https?://)?((([\dA-Za-z.-]+)\.([a-z.]{2,6}))|[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})(:[0-9]+)?/?(.*)$`)
	if err != nil {
		log.Fatalln("error compiling regex:", err)
	}
}
