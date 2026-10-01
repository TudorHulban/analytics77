package appanalytics

import (
	"fmt"
	"net/url"
	"strings"
)

func extractDateSelectors(v url.Values) {
	fmt.Println(v)

	raw := v.Encode()

	for pair := range strings.SplitSeq(raw, "&") {
		key, value, couldCut := strings.Cut(pair, "=")
		if !couldCut {
			continue
		}

		fmt.Println("key:", key)
		fmt.Println("value:", value)
	}
}
