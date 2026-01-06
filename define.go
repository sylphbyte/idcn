package idcn

import "github.com/sylphbyte/idcn/src"

// Info 身份信息
type Info struct {
	Age      int    `json:"age"`
	Birthday string `json:"birthday"`
	Gender   int    `json:"gender"` // 1男 0女
}

// IDCardInfo 身份证完整信息
type IDCardInfo struct {
	CardNo string       `json:"card_no"`
	Area   src.AreaInfo `json:"area"`
	Info   Info         `json:"info"`
}

// parsedID 解析后的身份证信息
type parsedID struct {
	addressCode  string
	birthdayCode string
	orderCode    string
	checkBit     string
	length       int
}
