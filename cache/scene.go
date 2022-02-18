package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy"
	"omo.msa.organization/proxy/nosql"
	"time"
)

const (
	SceneTypeOther SceneType = 0 // 未知
	SceneTypeSchool SceneType = 1 // 学校
	SceneTypeMuseum SceneType = 2 // 博物馆
	SceneTypeYoung SceneType = 3 // 少年宫
	SceneTypeNursery SceneType = 4 // 幼儿园
	SceneTypeMaker SceneType = 5 //实践中心，创客
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
	Supporter string
	Bucket string
	ShortName string
	Address nosql.AddressInfo
	members []string
	parents []string
	//Exhibitions []proxy.ShowingInfo
	groups []*GroupInfo
	rooms  []*RoomInfo
	Domains []proxy.DomainInfo
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
	db.Short = info.ShortName
	db.Status = uint8(SceneStatusIdle)
	db.Location = info.Location
	db.Address = info.Address
	db.Bucket = info.Bucket
	db.Members = make([]string, 0, 1)
	db.Parents = make([]string, 0, 1)
	db.Domains = make([]proxy.DomainInfo, 0, 1)
	db.Supporter = ""
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

func (mine *cacheContext)GetScenes(page, number uint32) (uint32,uint32,[]*SceneInfo) {
	if number < 1 {
		number = 10
	}
	if len(mine.scenes) <1 {
		return 0, 0, make([]*SceneInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, mine.scenes)
	return total, maxPage, list.([]*SceneInfo)
}

func (mine *cacheContext)GetScenesByArray(array []string) []*SceneInfo {
	list := make([]*SceneInfo, 0, len(array))
	for _, s := range array {
		info := mine.GetScene(s)
		if info != nil {
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext)GetScenesByParent(parent string, page, number uint32) (uint32,uint32,[]*SceneInfo) {
	if number < 1 {
		number = 10
	}
	if len(mine.scenes) < 1 {
		return 0, 0, make([]*SceneInfo, 0, 1)
	}
	all := make([]*SceneInfo, 0, 100)
	for _, scene := range mine.scenes {
		if scene.HadParent(parent) {
			all = append(all, scene)
		}
	}
	total, maxPage, list := checkPage(page, number, all)
	return total, maxPage, list.([]*SceneInfo)
}

func (mine *cacheContext)GetScenesByType(tp uint8) []*SceneInfo {
	list := make([]*SceneInfo, 0, 10)
	for _, scene := range mine.scenes {
		if uint8(scene.Type) == tp {
			list = append(list, scene)
		}
	}
	return list
}

func GetAllScenes() []*SceneInfo {
	return cacheCtx.scenes
}

func IsMasterUsed(uid string) bool {
	for i := 0;i < len(cacheCtx.scenes);i += 1{
		if cacheCtx.scenes[i].Master == uid {
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
	mine.Entity = db.Entity

	mine.ShortName = db.Short
	mine.Type = SceneType(db.Type)
	mine.Status = SceneStatus(db.Status)
	mine.members = db.Members
	mine.Address = db.Address
	mine.Supporter = db.Supporter
	mine.Bucket = db.Bucket
	mine.parents = db.Parents
	if mine.parents == nil {
		mine.parents = make([]string, 0, 1)
	}
	mine.Domains = db.Domains
	if mine.Domains == nil {
		mine.Domains = make([]proxy.DomainInfo, 0, 1)
	}
	//mine.Exhibitions = db.Displays
	//if mine.Exhibitions == nil {
	//	mine.Exhibitions = make([]proxy.ShowingInfo, 0, 1)
	//}
}

func (mine *SceneInfo)initGroups() {
	if mine.groups != nil {
		return
	}
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

func (mine *SceneInfo) initRooms() {
	if mine.rooms != nil {
		return
	}
	list,err := nosql.GetRoomsByScene(mine.UID)
	if err == nil {
		mine.rooms = make([]*RoomInfo, 0, len(list))
		for i := 0;i < len(list);i += 1 {
			tmp := new(RoomInfo)
			tmp.initInfo(list[i])
			mine.rooms = append(mine.rooms, tmp)
		}
	}else{
		mine.rooms = make([]*RoomInfo, 0, 1)
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
	if mine.Master == master {
		return nil
	}
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
	if mine.Cover == cover {
		return nil
	}
	err := nosql.UpdateSceneCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo)UpdateType(operator string, tp uint8) error {
	if uint8(mine.Type) == tp {
		return nil
	}
	err := nosql.UpdateSceneType(mine.UID, operator, tp)
	if err == nil {
		mine.Type = SceneType(tp)
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

func (mine *SceneInfo)UpdateDomains(operator string, list []proxy.DomainInfo) error {
	err := nosql.UpdateSceneDomains(mine.UID, operator, list)
	if err == nil {
		mine.Domains = list
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo)UpdateBucket(msg, operator string) error {
	err := nosql.UpdateSceneBucket(mine.UID, msg, operator)
	if err == nil {
		mine.Bucket = msg
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo)UpdateShortName(name, operator string) error {
	err := nosql.UpdateSceneShort(mine.UID, operator, name)
	if err == nil {
		mine.ShortName = name
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

func (mine *SceneInfo)UpdateSupporter(supporter, operator string) error {
	err := nosql.UpdateSceneSupporter(mine.UID, supporter, operator)
	if err == nil {
		mine.Supporter = supporter
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
				if i == len(mine.members) - 1 {
					mine.members = append(mine.members[:i])
				}else{
					mine.members = append(mine.members[:i], mine.members[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

func (mine *SceneInfo)Parents() []string {
	return mine.parents
}

func (mine *SceneInfo) UpdateParents(operator string, list []string) error {
	if list == nil {
		return errors.New("the children is nil")
	}
	err := nosql.UpdateSceneParents(mine.UID,operator, list)
	if err == nil {
		mine.parents = list
	}
	return err
}

//func (mine *SceneInfo)UpdateDisplay(uid, effect, skin, operator string, slots []string) error {
//	//if !mine.HadExhibition(uid){
//	//	return nil
//	//}
//	//array := make([]proxy.ShowingInfo, 0, len(mine.Exhibitions))
//	//for _, info := range mine.Exhibitions {
//	//	if info.UID == uid {
//	//		info.Effect = effect
//	//		info.Skin = skin
//	//		info.Slots = slots
//	//	}
//	//	array = append(array, info)
//	//}
//	//err := nosql.UpdateSceneDisplay(mine.UID, operator, array)
//	//if err == nil {
//	//	mine.Exhibitions = array
//	//}
//	//return err
//	return nil
//}

//func (mine *SceneInfo)PutOnDisplay(uid string) error {
//	if mine.HadExhibition(uid){
//		return nil
//	}
//	info := proxy.ShowingInfo{UID: uid, Effect: ""}
//	err := nosql.AppendSceneDisplay(mine.UID, &info)
//	if err == nil {
//		mine.Exhibitions = append(mine.Exhibitions, info)
//	}
//	return err
//	return nil
//}

//func (mine *SceneInfo)CancelDisplay(uid string) error {
//	if !mine.HadExhibition(uid){
//		return nil
//	}
//	err := nosql.SubtractSceneDisplay(mine.UID, uid)
//	if err == nil {
//		for i := 0;i < len(mine.Exhibitions);i += 1 {
//			if mine.Exhibitions[i].UID == uid {
//				if i == len(mine.Exhibitions) - 1 {
//					mine.Exhibitions = append(mine.Exhibitions[:i])
//				}else{
//					mine.Exhibitions = append(mine.Exhibitions[:i], mine.Exhibitions[i+1:]...)
//				}
//				break
//			}
//		}
//	}
//	return err
//	return nil
//}

func (mine *SceneInfo)HadParent(uid string) bool {
	if mine.parents == nil {
		return false
	}
	for _, item := range mine.parents {
		if item == uid {
			return true
		}
	}
	return false
}

func (mine *SceneInfo)CreateGroup(info *pb.ReqGroupAdd) (*GroupInfo, error) {
	mine.initGroups()
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
	mine.initGroups()
	for _, group := range mine.groups {
		if group.UID == uid {
			return true
		}
	}
	return false
}

func (mine *SceneInfo)HadGroupByName(name string) bool {
	mine.initGroups()
	for _, group := range mine.groups {
		if group.Name == name {
			return true
		}
	}
	return false
}

func (mine *SceneInfo)GetGroup(uid string) *GroupInfo {
	mine.initGroups()
	for _, group := range mine.groups {
		if group.UID == uid {
			return group
		}
	}
	return nil
}

func (mine *SceneInfo)RemoveGroup(uid, operator string) error {
	mine.initGroups()
	if !mine.HadGroup(uid) {
		return nil
	}
	err := nosql.RemoveGroup(uid, operator)
	if err == nil {
		for i := 0;i < len(mine.groups);i ++ {
			if mine.groups[i].UID == uid {
				if i == len(mine.groups) - 1 {
					mine.groups = append(mine.groups[:i])
				}else{
					mine.groups = append(mine.groups[:i], mine.groups[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

func (mine *SceneInfo)GetGroups(number, page uint32) (uint32,uint32,[]*GroupInfo) {
	mine.initGroups()
	if number < 1 {
		number = 10
	}
	if len(mine.groups) <1 {
		return 0, 0, make([]*GroupInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, mine.groups)
	return total, maxPage, list.([]*GroupInfo)
}

func (mine *SceneInfo)CreateRoom(info *pb.ReqRoomAdd) (*RoomInfo, error) {
	mine.initRooms()
	db := new(nosql.Room)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRoomNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = info.Operator
	db.Creator = info.Operator
	db.Name = info.Name
	db.Remark = info.Remark
	db.Scene = info.Owner
	db.Quotes = make([]string, 0, 1)
	db.Devices = make([]proxy.DeviceInfo, 0, 1)
	err := nosql.CreateRoom(db)
	if err == nil {
		tmp := new(RoomInfo)
		tmp.initInfo(db)
		mine.rooms = append(mine.rooms, tmp)
		return tmp,nil
	}
	return nil,err
}

func (mine *SceneInfo)HadRoom(uid string) bool {
	mine.initRooms()
	for _, item := range mine.groups {
		if item.UID == uid {
			return true
		}
	}
	return false
}

func (mine *SceneInfo)HadRoomByName(name string) bool {
	mine.initRooms()
	for _, item := range mine.rooms {
		if item.Name == name {
			return true
		}
	}
	return false
}

func (mine *SceneInfo)GetRooms() []*RoomInfo {
	mine.initRooms()
	return mine.rooms
}

func (mine *SceneInfo)GetRoom(uid string) *RoomInfo {
	mine.initRooms()
	for _, item := range mine.rooms {
		if item.UID == uid {
			return item
		}
	}
	return nil
}

func (mine *SceneInfo)GetRoomsByType(tp uint8) []*RoomInfo {
	mine.initRooms()
	list := make([]*RoomInfo, 0, len(mine.rooms))
	for _, item := range mine.rooms {
		if item.HadDeviceByType(tp) {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SceneInfo)GetRoomsByQuote(quote string) []*RoomInfo {
	mine.initRooms()
	list := make([]*RoomInfo, 0, len(mine.rooms))
	for _, item := range mine.rooms {
		if item.HadQuote(quote) {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SceneInfo)RemoveRoom(uid, operator string) error {
	if !mine.HadRoom(uid) {
		return nil
	}
	err := nosql.RemoveRoom(uid, operator)
	if err == nil {
		for i := 0;i < len(mine.rooms);i ++ {
			if mine.rooms[i].UID == uid {
				if i == len(mine.rooms) - 1 {
					mine.rooms = append(mine.rooms[:i])
				}else{
					mine.rooms = append(mine.rooms[:i], mine.rooms[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

