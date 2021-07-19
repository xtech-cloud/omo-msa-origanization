package nosql

import (
	"context"
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
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name     string   `json:"name" bson:"name"`
	Type     uint8    `json:"type" bson:"type"`
	Status   uint8    `json:"status" bson:"status"`
	Cover    string   `json:"cover" bson:"cover"`
	Master   string   `json:"master" bson:"master"`
	Remark   string   `json:"remark" bson:"remark"`
	Entity   string   `json:"entity" bson:"entity"`
	Location string   `json:"location" bson:"location"`
	Address AddressInfo `json:"address" bson:"address"`
	Members  []string `json:"members" bson:"members"`
	Exhibitions []string `json:"exhibitions" bson:"exhibitions"`
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

func GetSceneByMaster(user string) (*Scene, error) {
	msg := bson.M{"master":user}
	result, err := findOneBy(TableScene, msg)
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

func GetAllScenes() ([]*Scene, error) {
	var items = make([]*Scene, 0, 100)
	cursor, err1 := findAll(TableScene, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Scene)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateSceneBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func UpdateSceneMaster(uid, master, operator string) error {
	msg := bson.M{"master": master, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func UpdateSceneCover(uid string, icon, operator string) error {
	msg := bson.M{"cover": icon, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func UpdateSceneLocal(uid, local, operator string) error {
	msg := bson.M{"location": local, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func UpdateSceneAddress(uid, operator string, address AddressInfo) error {
	msg := bson.M{"address": address, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func UpdateSceneStatus(uid string, status uint8, operator string) error {
	msg := bson.M{"status": status, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func RemoveScene(uid, operator string) error {
	_, err := removeOne(TableScene, uid, operator)
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

func SubtractSceneMember(uid, member string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := removeElement(TableScene, uid, msg)
	return err
}

func AppendSceneDisplay(uid, display string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"exhibitions": display}
	_, err := appendElement(TableScene, uid, msg)
	return err
}

func SubtractSceneDisplay(uid, display string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"exhibitions": display}
	_, err := removeElement(TableScene, uid, msg)
	return err
}
