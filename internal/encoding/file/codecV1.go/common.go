package codecV1

import (
	"regexp"
)

const (
	v1MaxContentTotalSize = 4194304 // 4MB
	v1MaxContentChunkSize = 262144  // 512
)

var base64UrlSafeRegex = regexp.MustCompile("^(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}==|[A-Za-z0-9-_]{3}=|[A-Za-z0-9-_]{4})$")
