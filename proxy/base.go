package proxy

type ShowingInfo struct {
	UID string `json:"uid" bson:"uid"`
	Effect string `json:"effect" bson:"effect"`
	Skin string `json:"skin" bson:"skin"`
	Slots []string `json:"slots" bson:"slots"`
}