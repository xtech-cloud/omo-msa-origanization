package proxy

import "time"

type DisplayInfo struct {
	Type uint32 `json:"type" bson:"type"` //产品类型
	Group string `json:"group" bson:"group"` //所在组
	Updated time.Time `json:"updatedAt" bson:"updatedAt"`
	Showings []string `json:"showings" bson:"showings"` //预备展
}

type DomainInfo struct {
	Type uint8 `json:"type" bson:"type"`
	UID string `json:"uid" bson:"uid"`
	Name string `json:"name" bson:"name"`
	Remark string `json:"remark" bson:"remark"`
	Keywords string `json:"keywords" bson:"keywords"`
	Updated time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (mine *DisplayInfo)Clone() *DisplayInfo {
	tmp := new(DisplayInfo)
	tmp.Type = mine.Type
	tmp.Group = mine.Group
	tmp.Showings = mine.Showings
	tmp.Updated = mine.Updated
	return tmp
}