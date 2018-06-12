package main

type (
	ReqCryptOuter struct {
		Data   string `json:"data" form:"data"`
		PriKey string `json:"pri_key" form:"pri_key"`
	}
)
