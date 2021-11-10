package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

func main() {
	t := map[string]interface{}{
		"Hello": "world",
		"Say": map[string]interface{}{
			"China": map[string]interface{}{
				"Shenzhen": "你好",
				"Beijing":  "你好",
			},
			"US": "Hello",
		},
		"Alibaba": "Hello",
	}
	buf := new(bytes.Buffer)
	treestr(t, buf, make([]string, 0))

	s := buf.String()

	fmt.Println(s)
}

//alibaba=Hello&&hello=world&&say.china.beijing=你好&&say.china.shenzhen=你好&&say.us=Hello
func treestr(t map[string]interface{}, buf *bytes.Buffer, ps []string) {
	if len(t) == 0 {
		return
	}

	ks := make([]string, 0, len(t))
	for k := range t {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	for _, k := range ks {
		ps = append(ps, k)
		v := t[k]

		switch vt := v.(type) {
		case map[string]interface{}:
			treestr(vt, buf, ps)
		case string:
			if buf.Len() > 0 {
				buf.WriteString("&&")
			}
			for i, s := range ps {
				if i > 0 {
					buf.WriteByte(byte('.'))
				}
				buf.WriteString(strings.ToLower(s))
			}
			buf.WriteByte(byte('='))
			buf.WriteString(vt)
		default:
			fmt.Println(v, vt)
		}

		ps = ps[0 : len(ps)-1]
	}
}
