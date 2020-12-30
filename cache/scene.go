package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
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
	baseInfo
	Type SceneType
	Status SceneStatus
	Location string
	Cover string
	Remark string
	Master string
	Entity string
	Address nosql.AddressInfo
	members []string
	groups []*GroupInfo
}

func (mine *cacheContext)CreateScene(info *SceneInfo) error {
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
	db.Entity = info.Entity
	db.Status = uint8(SceneStatusIdle)
	db.Location = info.Location
	db.Address = info.Address
	db.Members = make([]string, 0, 1)
	err := nosql.CreateScene(db)
	if err == nil {
		info.initInfo(db)
		cacheCtx.scenes = append(cacheCtx.scenes, info)
	}
	return err
}

func (mine *cacheContext)GetScene(uid string) *SceneInfo {
	if len(uid) < 2 {
		return nil
	}
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

func (mine *cacheContext)GetSceneByMember(uid string) *SceneInfo {
	for i := 0;i < len(cacheCtx.scenes);i += 1{
		if cacheCtx.scenes[i].HadMember(uid) {
			return cacheCtx.scenes[i]
		}
	}
	db,err := nosql.GetSceneByMaster(uid)
	if err == nil {
		info := new(SceneInfo)
		info.initInfo(db)
		cacheCtx.scenes = append(cacheCtx.scenes, info)
		return info
	}
	return nil
}

func (mine *cacheContext)GetScenes(number, page uint32) (uint32,uint32,[]*SceneInfo) {
	if number < 1 {
		number = 10
	}
	if len(mine.scenes) <1 {
		return 0, 0, make([]*SceneInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, mine.scenes)
	return total, maxPage, list.([]*SceneInfo)
}

func GetAllScenes() []*SceneInfo {
	return cacheCtx.scenes
}

func IsMasterUsed(uid string) bool {
	for i := 0;i < len(cacheCtx.scenes);i += 1{
		if cacheCtx.scenes[i].UID == uid {
			return true
		}
	}
	return false
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
	mine.Master = db.Master
	mine.Location = db.Location
	mine.Type = SceneType(db.Type)
	mine.Status = SceneStatus(db.Status)
	mine.members = db.Members
	mine.Address = db.Address
	groups,err := nosql.GetGroupsByScene(mine.UID)
	if err == nil {
		mine.groups = make([]*GroupInfo, 0, len(groups))
		for i := 0;i < len(groups);i += 1 {
			tmp := new(GroupInfo)
			tmp.initInfo(groups[i])
			mine.groups = append(mine.groups, tmp)
		}
	}else{
		mine.groups = make([]*GroupInfo, 0, 1)
	}
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
	if IsMasterUsed(master) {
		return errors.New("the master had used by other scene")
	}
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

func (mine *SceneInfo)UpdateAddress(country, province, city, zone, operator string) error {
	addr := nosql.AddressInfo{Country: country, Province: province, City: city, Zone: zone}
	err := nosql.UpdateSceneAddress(mine.UID, operator, addr)
	if err == nil {
		mine.Address = addr
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
	if mine.Master == member {
		return true
	}
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

func (mine *SceneInfo)CreateGroup(info *pb.ReqGroupAdd) (*GroupInfo, error) {
	db := new(nosql.Group)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetGroupNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = info.Operator
	db.Creator = info.Operator
	db.Name = info.Name
	db.Cover = info.Cover
	db.Remark = info.Remark
	db.Location = info.Location
	db.Contact = info.Contact
	db.Scene = info.Scene
	db.Address = nosql.AddressInfo{
		Country: info.Address.Country,
		Province: info.Address.Province,
		City: info.Address.City,
		Zone: info.Address.Zone,
	}
	db.Members = make([]string, 0, 1)
	err := nosql.CreateGroup(db)
	if err == nil {
		tmp := new(GroupInfo)
		tmp.initInfo(db)
		mine.groups = append(mine.groups, tmp)
		return tmp,nil
	}
	return nil,err
}

func (mine *SceneInfo)HadGroup(uid string) bool {
	for _, group := range mine.groups {
		if group.UID == uid {
			return true
		}
	}
	return false
}

func (mine *SceneInfo)HadGroupByName(name string) bool {
	for _, group := range mine.groups {
		if group.Name == name {
			return true
		}
	}
	return false
}

func (mine *SceneInfo)GetGroup(uid string) *GroupInfo {
	for _, group := range mine.groups {
		if group.UID == uid {
			return group
		}
	}
	return nil
}

func (mine *SceneInfo)RemoveGroup(uid, operator string) error {
	if !mine.HadGroup(uid) {
		return nil
	}
	err := nosql.RemoveGroup(uid, operator)
	if err == nil {
		for i := 0;i < len(mine.groups);i ++ {
			if mine.groups[i].UID == uid {
				mine.groups = append(mine.groups[:i], mine.groups[i+1:]...)
				break
			}
		}
	}
	return err
}

func (mine *SceneInfo)GetGroups(number, page uint32) (uint32,uint32,[]*GroupInfo) {
	if number < 1 {
		number = 10
	}
	if len(mine.groups) <1 {
		return 0, 0, make([]*GroupInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, mine.groups)
	return total, maxPage, list.([]*GroupInfo)
}
