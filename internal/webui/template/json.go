package template

import "encoding/json"

func toJSON(value any) string {
	buf, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	return string(buf)
}
