package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/emyassine/ikaitla/framework/ui/output"
	"github.com/emyassine/ikaitla/framework/ui/theme"
)

func RenderKeyValue(out *output.Output, pairs map[string]string) {
	// stable order
	keys := make([]string, 0, len(pairs))
	maxKeyLen := 0
	for k := range pairs {
		keys = append(keys, k)
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := pairs[k]
		key := k + strings.Repeat(" ", maxKeyLen-len(k))
		if out.ColorsEnabled() {
			out.Printf("%s: %s", out.Stylize(key, theme.Slate900, theme.Bold), v)
		} else {
			out.Printf("%s: %s", key, v)
		}
	}
}

func RenderKeyValueAny(out *output.Output, pairs map[string]any) {
	keys := make([]string, 0, len(pairs))
	maxKeyLen := 0
	for k := range pairs {
		keys = append(keys, k)
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		key := k + strings.Repeat(" ", maxKeyLen-len(k))
		v := fmt.Sprintf("%v", pairs[k])
		if out.ColorsEnabled() {
			out.Printf("%s: %s", out.Stylize(key, theme.Slate900, theme.Bold), v)
		} else {
			out.Printf("%s: %s", key, v)
		}
	}
}
