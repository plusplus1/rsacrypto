package impl

type rsaParams struct {
	Data   string `json:"data" form:"data"`
	PriKey string `json:"pri_key" form:"pri_key"`
	PubKey string `json:"pub_key" form:"pub_key"`
}
