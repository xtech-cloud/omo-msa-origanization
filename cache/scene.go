package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.organization/proxy"
	"omo.msa.organization/proxy/nosql"
	"omo.msa.organization/tool"
	"time"
)

const (
	SceneTypeOther   SceneType = 0 // 未知
	SceneTypeSchool  SceneType = 1 // 学校
	SceneTypeMuseum  SceneType = 2 // 博物馆
	SceneTypeYoung   SceneType = 3 // 少年宫
	SceneTypeNursery SceneType = 4 // 幼儿园
	SceneTypeMaker   SceneType = 5 //实践中心，创客
)

const (
	SceneStatusIdle  SceneStatus = 0
	SceneStatusFroze SceneStatus = 1
)

type SceneType uint8

type SceneStatus uint8

type SceneInfo struct {
	baseInfo
	Type      SceneType
	Status    SceneStatus
	Limit     uint16
	Location  string
	Cover     string
	Remark    string
	Master    string
	Entity    string
	Supporter string
	Bucket    string
	ShortName string
	Address   nosql.AddressInfo
	members   []string
	parents   []string
	Questions []string
	Domains   []proxy.DomainInfo
	groups    []*GroupInfo
	rooms     []*RoomInfo
}

func (mine *cacheContext) CreateScene(info *SceneInfo) error {
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
	db.Limit = info.Limit
	db.Members = make([]string, 0, 1)
	db.Parents = make([]string, 0, 1)
	db.Questions = make([]string, 0, 1)
	db.Domains = make([]proxy.DomainInfo, 0, 1)
	db.Supporter = ""
	err := nosql.CreateScene(db)
	if err == nil {
		info.initInfo(db)
		cacheCtx.scenes = append(cacheCtx.scenes, info)
	}
	return err
}

func (mine *cacheContext) GetScene(uid string) *SceneInfo {
	if len(uid) < 2 {
		return nil
	}
	for i := 0; i < len(cacheCtx.scenes); i += 1 {
		if cacheCtx.scenes[i].UID == uid {
			return cacheCtx.scenes[i]
		}
	}
	db, err := nosql.GetScene(uid)
	if err == nil {
		info := new(SceneInfo)
		info.initInfo(db)
		cacheCtx.scenes = append(cacheCtx.scenes, info)
		return info
	}
	return nil
}

func (mine *cacheContext) GetSceneByMember(uid string) *SceneInfo {
	for i := 0; i < len(cacheCtx.scenes); i += 1 {
		if cacheCtx.scenes[i].HadMember(uid) {
			return cacheCtx.scenes[i]
		}
	}
	db, err := nosql.GetSceneByMaster(uid)
	if err == nil {
		info := new(SceneInfo)
		info.initInfo(db)
		cacheCtx.scenes = append(cacheCtx.scenes, info)
		return info
	}
	return nil
}

func (mine *cacheContext) GetScenes(page, number uint32) (uint32, uint32, []*SceneInfo) {
	if number < 1 {
		number = 10
	}
	if len(mine.scenes) < 1 {
		return 0, 0, make([]*SceneInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, mine.scenes)
	return total, maxPage, list.([]*SceneInfo)
}

func (mine *cacheContext) GetScenesByArray(array []string) []*SceneInfo {
	list := make([]*SceneInfo, 0, len(array))
	for _, s := range array {
		info := mine.GetScene(s)
		if info != nil {
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetScenesByParent(parent string, page, number uint32) (uint32, uint32, []*SceneInfo) {
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

func (mine *cacheContext) GetScenesByType(tp uint8) []*SceneInfo {
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
	for i := 0; i < len(cacheCtx.scenes); i += 1 {
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
		for i := 0; i < len(cacheCtx.scenes); i += 1 {
			if cacheCtx.scenes[i].UID == uid {
				cacheCtx.scenes = append(cacheCtx.scenes[:i], cacheCtx.scenes[i+1:]...)
				break
			}
		}
	}
	return err
}

func (mine *SceneInfo) initInfo(db *nosql.Scene) {
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
	mine.Limit = db.Limit

	mine.ShortName = db.Short
	mine.Type = SceneType(db.Type)
	mine.Status = SceneStatus(db.Status)
	mine.members = db.Members
	mine.Address = db.Address
	mine.Supporter = db.Supporter
	mine.Bucket = db.Bucket
	mine.Questions = db.Questions
	mine.parents = db.Parents
	if mine.parents == nil {
		mine.parents = make([]string, 0, 1)
	}
	mine.Domains = db.Domains
	if mine.Domains == nil {
		mine.Domains = make([]proxy.DomainInfo, 0, 1)
	}
}

func (mine *SceneInfo) initGroups() {
	if mine.groups != nil {
		return
	}
	groups, err := nosql.GetGroupsByScene(mine.UID)
	if err == nil {
		mine.groups = make([]*GroupInfo, 0, len(groups))
		for i := 0; i < len(groups); i += 1 {
			tmp := new(GroupInfo)
			tmp.initInfo(groups[i])
			mine.groups = append(mine.groups, tmp)
		}
	} else {
		mine.groups = make([]*GroupInfo, 0, 1)
	}
}

func (mine *SceneInfo) initRooms() {
	if mine.rooms != nil {
		return
	}
	list, err := nosql.GetRoomsByScene(mine.UID)
	if err == nil {
		mine.rooms = make([]*RoomInfo, 0, len(list))
		for i := 0; i < len(list); i += 1 {
			tmp := new(RoomInfo)
			tmp.initInfo(list[i])
			mine.rooms = append(mine.rooms, tmp)
		}
	} else {
		mine.rooms = make([]*RoomInfo, 0, 1)
	}
}

func (mine *SceneInfo) UpdateBase(name, remark, operator string) error {
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

func (mine *SceneInfo) UpdateMaster(master, operator string) error {
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

func (mine *SceneInfo) UpdateCover(cover, operator string) error {
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

func (mine *SceneInfo) UpdateType(operator string, tp uint8) error {
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

func (mine *SceneInfo) UpdateLocation(local, operator string) error {
	err := nosql.UpdateSceneLocal(mine.UID, local, operator)
	if err == nil {
		mine.Location = local
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo) UpdateDomains(operator string, list []proxy.DomainInfo) error {
	err := nosql.UpdateSceneDomains(mine.UID, operator, list)
	if err == nil {
		mine.Domains = list
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo) UpdateQuestions(operator string, arr []string) error {
	err := nosql.UpdateSceneQuestions(mine.UID, operator, arr)
	if err == nil {
		mine.Questions = arr
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo) UpdateLimit(operator string, limit int) error {
	err := nosql.UpdateSceneLimit(mine.UID, operator, uint16(limit))
	if err == nil {
		mine.Limit = uint16(limit)
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo) UpdateBucket(msg, operator string) error {
	err := nosql.UpdateSceneBucket(mine.UID, msg, operator)
	if err == nil {
		mine.Bucket = msg
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo) UpdateShortName(name, operator string) error {
	err := nosql.UpdateSceneShort(mine.UID, operator, name)
	if err == nil {
		mine.ShortName = name
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo) UpdateAddress(country, province, city, zone, operator string) error {
	addr := nosql.AddressInfo{Country: country, Province: province, City: city, Zone: zone}
	err := nosql.UpdateSceneAddress(mine.UID, operator, addr)
	if err == nil {
		mine.Address = addr
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo) UpdateStatus(st SceneStatus, operator string) error {
	err := nosql.UpdateSceneStatus(mine.UID, uint8(st), operator)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo) UpdateSupporter(supporter, operator string) error {
	err := nosql.UpdateSceneSupporter(mine.UID, supporter, operator)
	if err == nil {
		mine.Supporter = supporter
		mine.Operator = operator
	}
	return err
}

func (mine *SceneInfo) HadMember(member string) bool {
	if mine.Master == member {
		return true
	}
	for i := 0; i < len(mine.members); i += 1 {
		if mine.members[i] == member {
			return true
		}
	}
	return false
}

func (mine *SceneInfo) AllMembers() []string {
	return mine.members
}

func (mine *SceneInfo) AppendMember(member string) error {
	if mine.HadMember(member) {
		return errors.New("the member had existed")
	}
	err := nosql.AppendSceneMember(mine.UID, member)
	if err == nil {
		mine.members = append(mine.members, member)
	}
	return err
}

func (mine *SceneInfo) SubtractMember(member string) error {
	if !mine.HadMember(member) {
		return errors.New("the member not existed")
	}
	err := nosql.SubtractSceneMember(mine.UID, member)
	if err == nil {
		for i := 0; i < len(mine.members); i += 1 {
			if mine.members[i] == member {
				if i == len(mine.members)-1 {
					mine.members = append(mine.members[:i])
				} else {
					mine.members = append(mine.members[:i], mine.members[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

func (mine *SceneInfo) Parents() []string {
	return mine.parents
}

func (mine *SceneInfo) UpdateParents(operator string, list []string) error {
	if list == nil {
		return errors.New("the children is nil")
	}
	err := nosql.UpdateSceneParents(mine.UID, operator, list)
	if err == nil {
		mine.parents = list
	}
	return err
}

func (mine *SceneInfo) HadParent(uid string) bool {
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

//region Group Fun
func (mine *SceneInfo) CreateGroup(info *pb.ReqGroupAdd) (*GroupInfo, error) {
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
		Country:  info.Address.Country,
		Province: info.Address.Province,
		City:     info.Address.City,
		Zone:     info.Address.Zone,
	}
	db.Members = make([]string, 0, 1)
	err := nosql.CreateGroup(db)
	if err == nil {
		tmp := new(GroupInfo)
		tmp.initInfo(db)
		mine.groups = append(mine.groups, tmp)
		return tmp, nil
	}
	return nil, err
}

func (mine *SceneInfo) HadGroup(uid string) bool {
	mine.initGroups()
	for _, group := range mine.groups {
		if group.UID == uid {
			return true
		}
	}
	return false
}

func (mine *SceneInfo) HadGroupByName(name string) bool {
	mine.initGroups()
	for _, group := range mine.groups {
		if group.Name == name {
			return true
		}
	}
	return false
}

func (mine *SceneInfo) GetGroup(uid string) *GroupInfo {
	mine.initGroups()
	for _, group := range mine.groups {
		if group.UID == uid {
			return group
		}
	}
	return nil
}

func (mine *SceneInfo) RemoveGroup(uid, operator string) error {
	mine.initGroups()
	if !mine.HadGroup(uid) {
		return nil
	}
	err := nosql.RemoveGroup(uid, operator)
	if err == nil {
		for i := 0; i < len(mine.groups); i++ {
			if mine.groups[i].UID == uid {
				if i == len(mine.groups)-1 {
					mine.groups = append(mine.groups[:i])
				} else {
					mine.groups = append(mine.groups[:i], mine.groups[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

func (mine *SceneInfo) GetGroups(number, page uint32) (uint32, uint32, []*GroupInfo) {
	mine.initGroups()
	if number < 1 {
		number = 10
	}
	if len(mine.groups) < 1 {
		return 0, 0, make([]*GroupInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, mine.groups)
	return total, maxPage, list.([]*GroupInfo)
}

//endregion

//region Region
func (mine *SceneInfo) CreateRegion(info *pb.ReqRegionAdd) (*RegionInfo, error) {
	db := new(nosql.Region)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRegionNextID()
	db.CreatedTime = time.Now()
	db.UpdatedTime = time.Now()
	db.Operator = info.Operator
	db.Creator = info.Operator
	db.Name = info.Name
	db.Remark = info.Remark
	db.Location = info.Location
	db.Scene = info.Scene
	db.Parent = info.Parent
	db.Master = ""
	if info.Address != nil {
		db.Address = nosql.AddressInfo{
			Country:  info.Address.Country,
			Province: info.Address.Province,
			City:     info.Address.City,
			Zone:     info.Address.Zone,
		}
	} else {
		db.Address = nosql.AddressInfo{
			Country:  "",
			Province: "",
			City:     "",
			Zone:     "",
		}
	}

	db.Members = make([]string, 0, 1)
	err := nosql.CreateRegion(db)
	if err == nil {
		tmp := new(RegionInfo)
		tmp.initInfo(db)
		return tmp, nil
	}
	return nil, err
}

func (mine *SceneInfo) HadRegion(uid string) bool {
	info, _ := mine.GetRegion(uid)
	if info != nil {
		return true
	}
	return false
}

func (mine *SceneInfo) HadRegionByName(name string) bool {
	arr := cacheCtx.GetRegionsByScene(mine.UID)
	for _, group := range arr {
		if group.Name == name {
			return true
		}
	}
	return false
}

func (mine *SceneInfo) GetRegion(uid string) (*RegionInfo, error) {
	return cacheCtx.GetRegion(uid)
}

func (mine *SceneInfo) RemoveRegion(uid, operator string) error {
	info, er := mine.GetRegion(uid)
	if er != nil {
		return er
	}
	if info == nil {
		return nil
	}
	return info.Remove(operator)
}

func (mine *SceneInfo) GetRegions(number, page uint32) (uint32, uint32, []*RegionInfo) {
	array := cacheCtx.GetRegionsByScene(mine.UID)
	if number < 1 {
		number = 10
	}
	if len(array) < 1 {
		return 0, 0, make([]*RegionInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, array)
	return total, maxPage, list.([]*RegionInfo)
}

//endregion

//region Room Fun
func (mine *SceneInfo) CreateRoom(info *pb.ReqRoomAdd) (*RoomInfo, error) {
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
	//db.Displays = make([]proxy.DisplayInfo, 0, 1)
	err := nosql.CreateRoom(db)
	if err == nil {
		tmp := new(RoomInfo)
		tmp.initInfo(db)
		mine.rooms = append(mine.rooms, tmp)
		return tmp, nil
	}
	return nil, err
}

func (mine *SceneInfo) HadRoom(uid string) bool {
	mine.initRooms()
	for _, item := range mine.rooms {
		if item.UID == uid {
			return true
		}
	}
	return false
}

func (mine *SceneInfo) HadRoomByName(name string) bool {
	mine.initRooms()
	for _, item := range mine.rooms {
		if item.Name == name {
			return true
		}
	}
	return false
}

func (mine *SceneInfo) GetRooms() []*RoomInfo {
	mine.initRooms()
	return mine.rooms
}

func (mine *SceneInfo) GetRoom(uid string) *RoomInfo {
	mine.initRooms()
	for _, item := range mine.rooms {
		if item.UID == uid {
			return item
		}
	}
	return nil
}

func (mine *SceneInfo) GetRoomsByType(tp uint8) []*RoomInfo {
	mine.initRooms()
	list := make([]*RoomInfo, 0, len(mine.rooms))
	for _, item := range mine.rooms {
		if item.HadDeviceByType(tp) {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SceneInfo) GetRoomsByQuote(quote string) []*RoomInfo {
	mine.initRooms()
	list := make([]*RoomInfo, 0, len(mine.rooms))
	for _, item := range mine.rooms {
		if item.HadQuote(quote) {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SceneInfo) GetRoomsByDevice(sn string) []*RoomInfo {
	mine.initRooms()
	list := make([]*RoomInfo, 0, len(mine.rooms))
	for _, item := range mine.rooms {
		if item.HadDevice(sn) {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SceneInfo) GetDevice(sn string) *DeviceInfo {
	mine.initRooms()
	for _, item := range mine.rooms {
		device := item.GetDevice(sn)
		if device != nil {
			return device
		}
	}
	return nil
}

func (mine *SceneInfo) GetDevices(arr []string) ([]*DeviceInfo, error) {
	all, err := nosql.GetDevicesByScene(mine.UID)
	if err != nil {
		return nil, err
	}
	list := make([]*DeviceInfo, 0, len(all))
	for _, device := range all {
		if tool.HasItem(arr, device.SN) {
			info := new(DeviceInfo)
			info.initInfo(device)
			list = append(list, info)
		}
	}
	return list, nil
}

func (mine *SceneInfo) GetDevicesByArea(area string) ([]*DeviceInfo, error) {
	all, err := nosql.GetDevicesByArea(mine.UID, area)
	if err != nil {
		return nil, err
	}
	list := make([]*DeviceInfo, 0, len(all))
	for _, device := range all {
		info := new(DeviceInfo)
		info.initInfo(device)
		list = append(list, info)
	}
	return list, nil
}

func (mine *SceneInfo) RemoveRoom(uid, operator string) error {
	if !mine.HadRoom(uid) {
		return nil
	}
	err := nosql.RemoveRoom(uid, operator)
	if err == nil {
		for i := 0; i < len(mine.rooms); i++ {
			if mine.rooms[i].UID == uid {
				if i == len(mine.rooms)-1 {
					mine.rooms = append(mine.rooms[:i])
				} else {
					mine.rooms = append(mine.rooms[:i], mine.rooms[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

func (mine *SceneInfo) ClearQuotes(operator string, list []string) {
	mine.initRooms()
	for _, room := range mine.rooms {
		if room.HadQuotes(list) {
			_ = room.UpdateQuotes(operator, make([]string, 0, 1))
			break
		}
	}
}

//endregion
