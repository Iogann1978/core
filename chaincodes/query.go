package chaincodes

import "encoding/json"

type Query struct {
	Selector map[string]interface{} `json:"selector,omitempty"`
	Limit    *int                   `json:"limit,omitempty"`
	Offset   *int                   `json:"skip,omitempty"`
}

func (q Query) String() string {
	j, _ := json.Marshal(q)
	return string(j)
}
