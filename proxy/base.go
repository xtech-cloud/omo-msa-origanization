package proxy

type ShowingInfo struct {
	UID string `json:"uid" bson:"uid"`
	Effect string `json:"effect" bson:"effect"`
	Slots []string `json:"slots" bson:"slots"`
}