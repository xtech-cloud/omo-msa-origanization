package nosql

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Scene struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Name   string                `json:"name" bson:"name"`
	Type   uint8 `json:"type" bson:"type"`
	Status uint8 `json:"status" bson:"status"`
	Cover string `json:"cover" bson:"cover"`
	Master string `json:"master" bson:"master"`
	Remark string `json:"remark" bson:"remark"`
	Location string `json:"location" bson:"location"`
	Members []string                `json:"members" bson:"members"`
}

func CreateScene(info *Scene) error {
	_, err := insertOne(TableScene, info)
	if err != nil {
		return err
	}
	return nil
}

func GetSceneNextID() uint64 {
	num, _ := getSequenceNext(TableScene)
	return num
}

func GetScene(uid string) (*Scene, error) {
	result, err := findOne(TableScene, uid)
	if err != nil {
		return nil, err
	}
	model := new(Scene)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllScenes(uid string) (*Scene, error) {
	result, err := findAll(TableScene, 0)
	if err != nil {
		return nil, err
	}
	model := new(Scene)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateSceneBase(uid, name, remark string) error {
	msg := bson.M{"name": name, "remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func UpdateSceneCover(uid string, icon string) error {
	msg := bson.M{"cover": icon, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func UpdateSceneLocal(uid string, local string) error {
	msg := bson.M{"location": local, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func UpdateSceneStatus(uid string, status uint8) error {
	msg := bson.M{"status": status, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func RemoveScene(uid string) error {
	_, err := removeOne(TableScene, uid)
	return err
}

func AppendSceneMember(uid string, member string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := appendElement(TableScene, uid, msg)
	return err
}

func SubtractSceneMember(uid , member string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := removeElement(TableScene, uid, msg)
	return err
}
