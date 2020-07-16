package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy/nosql"
	"time"
)

const (
	SceneTypeOther SceneType = 0
	SceneTypeSchool SceneType = 1
	SceneTypeMuseum SceneType = 2
)

const (
	SceneStatusIdle SceneStatus = 0
	SceneStatusFroze  SceneStatus = 1
)

type SceneType uint8

type SceneStatus uint8

type SceneInfo struct {
	BaseInfo
	Type SceneType
	Status SceneStatus
	Location string
	Cover string
	Remark string
	Master string
	members []string
}

func CreateScene(info *SceneInfo) error {
	db := new(nosql.Scene)
	db.UID = primitive.NewObjectID()
	db.Type = uint8(info.Type)
	db.ID = nosql.GetSceneNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = info.Operator
	db.Creator = info.Creator
	db.Name = info.Name
	db.Cover = info.Cover
	db.Remark = info.Remark
	db.Status = uint8(SceneStatusIdle)
	db.Location = info.Location
	db.Members = make([]string, 0, 1)
	err := nosql.CreateScene(db)
	if err == nil {
		info.initInfo(db)
		cacheCtx.scenes = append(cacheCtx.scenes, info)
	}
	return err
}

func GetScene(uid string) *SceneInfo {
	for i := 0;i < len(cacheCtx.scenes);i += 1{
		if cacheCtx.scenes[i].UID == uid {
			return cacheCtx.scenes[i]
		}
	}
	db,err := nosql.GetScene(uid)
	if err == nil {
		info := new(SceneInfo)
		info.initInfo(db)
		cacheCtx.scenes = append(cacheCtx.scenes, info)
		return info
	}
	return nil
}

func GetScenes(number, page int32) (int32,int32,[]*SceneInfo) {
	list := make([]*SceneInfo, 0, number)
	length := int32(len(cacheCtx.scenes))
	max := length / number + 1
	for i := 0;i < len(cacheCtx.scenes);i += 1{
		t := int32(i) / number + 1
		if t == page {
			list = append(list, cacheCtx.scenes[i])
		}
	}

	return length, max, list
}

func GetAllScenes() []*SceneInfo {
	return cacheCtx.scenes
}

func RemoveScene(uid, operator string) error {
	if len(uid) < 1 {
		return errors.New("the scene uid is empty")
	}
	err := nosql.RemoveScene(uid, operator)
	if err == nil {
		for i := 0;i < len(cacheCtx.scenes);i += 1 {
			if cacheCtx.scenes[i].UID == uid {
				cacheCtx.scenes = append(cacheCtx.scenes[:i], cacheCtx.scenes[i+1:]...)
				break
			}
		}
	}
	return err
}

func (mine *SceneInfo)initInfo(db *nosql.Scene)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Cover = db.Cover
	mine.Remark = db.Remark
	mine.Location = db.Location
	mine.Type = SceneType(db.Type)
	mine.Status = SceneStatus(db.Status)
	mine.members = db.Members
}

func (mine *SceneInfo)UpdateBase(name, remark, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	err := nosql.UpdateSceneBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo)UpdateMaster(master, operator string) error {
	err := nosql.UpdateSceneMaster(mine.UID, master, operator)
	if err == nil {
		mine.Master = master
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo)UpdateCover(cover, operator string) error {
	err := nosql.UpdateSceneCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo)UpdateLocation(local, operator string) error {
	err := nosql.UpdateSceneLocal(mine.UID, local, operator)
	if err == nil {
		mine.Location = local
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo)UpdateStatus(st SceneStatus, operator string) error {
	err := nosql.UpdateSceneStatus(mine.UID, uint8(st), operator)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo)HadMember(member string) bool {
	for i := 0;i < len(mine.members);i += 1 {
		if mine.members[i] == member {
			return true
		}
	}
	return false
}

func (mine *SceneInfo)AllMembers() []string {
	return mine.members
}

func (mine *SceneInfo)AppendMember(member string) error {
	if mine.HadMember(member){
		return errors.New("the member had existed")
	}
	err := nosql.AppendSceneMember(mine.UID, member)
	if err == nil {
		mine.members = append(mine.members, member)
	}
	return err
}

func (mine *SceneInfo)SubtractMember(member string) error {
	if !mine.HadMember(member){
		return errors.New("the member not existed")
	}
	err := nosql.SubtractSceneMember(mine.UID, member)
	if err == nil {
		for i := 0;i < len(mine.members);i += 1 {
			if mine.members[i] == member {
				mine.members = append(mine.members[:i], mine.members[i+1:]...)
				break
			}
		}
	}
	return err
}
