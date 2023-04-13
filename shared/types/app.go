package types

type Error struct {
	Message string `json:"message"`
	Field   string `json:"field"`
	Tag     string `json:"tag"`
}

type GoodResponse struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
}

type BadResponse struct {
	Code  int    `json:"code"`
	Error *Error `json:"error"`
}