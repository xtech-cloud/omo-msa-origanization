package proxy

import "time"

//type ShowingInfo struct {
//	UID string `json:"uid" bson:"uid"`
//	Effect string `json:"effect" bson:"effect"`
//	Skin string `json:"skin" bson:"skin"`
//	Slots []string `json:"slots" bson:"slots"`
//	Updated time.Time `json:"updatedAt" bson:"updatedAt"`
//}

type DeviceInfo struct {
	Type uint8 `json:"type" bson:"type"`
	SN string `json:"sn" bson:"sn"`
	Remark string `json:"remark" bson:"remark"`
	Updated time.Time `json:"updatedAt" bson:"updatedAt"`
}

type DomainInfo struct {
	Type uint8 `json:"type" bson:"type"`
	UID string `json:"uid" bson:"uid"`
	Remark string `json:"remark" bson:"remark"`
	Updated time.Time `json:"updatedAt" bson:"updatedAt"`
}