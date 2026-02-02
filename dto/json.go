package dto

import "encoding/json"

type JSON struct{}

func (j *JSON) ToObject(
	data []byte,
	obj any,
) error {
	return json.Unmarshal(data, obj)
}
