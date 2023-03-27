package cache

import (
	"errors"
	"omo.msa.organization/proxy/nosql"
)

type GroupInfo struct {
	baseInfo
	Remark    string
	Contact   string
	Cover     string
	Master    string
	Assistant string
	Address   nosql.AddressInfo
	Location  string
	Scene     string
	members   []string
}

func (mine *cacheContext) GetGroup(uid string) *GroupInfo {
	for _, scene := range mine.scenes {
		scene.initGroups()
		group := scene.GetGroup(uid)
		if group != nil {
			return group
		}
	}
	return nil
}

func (mine *cacheContext) GetGroupByMember(uid string) []*GroupInfo {
	list := make([]*GroupInfo, 0, 5)
	for _, scene := range mine.scenes {
		scene.initGroups()
		for _, group := range scene.groups {
			if group.HadMember(uid) {
				list = append(list, group)
			}
		}
	}
	return list
}

func (mine *cacheContext) GetGroupByContact(phone string) []*GroupInfo {
	list := make([]*GroupInfo, 0, 5)
	for _, scene := range mine.scenes {
		scene.initGroups()
		for _, group := range scene.groups {
			if group.Contact == phone {
				list = append(list, group)
			}
		}
	}
	return list
}

func (mine *cacheContext) RemoveGroup(uid, operator string) error {
	for _, scene := range mine.scenes {
		scene.initGroups()
		if scene.HadGroup(uid) {
			return scene.RemoveGroup(uid, operator)
		}
	}
	return nil
}

func (mine *GroupInfo) initInfo(db *nosql.Group) {
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
	mine.Assistant = db.Assistant
	mine.Contact = db.Contact
	mine.members = db.Members
	mine.Address = db.Address
	mine.Scene = db.Scene
}

func (mine *GroupInfo) UpdateBase(name, remark, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	err := nosql.UpdateGroupBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *GroupInfo) UpdateContact(phone, operator string) error {
	err := nosql.UpdateGroupContact(mine.UID, phone, operator)
	if err == nil {
		mine.Contact = phone
		mine.Operator = operator
	}
	return err
}

func (mine *GroupInfo) UpdateMaster(master, operator string) error {
	err := nosql.UpdateGroupMaster(mine.UID, master, operator)
	if err == nil {
		mine.Master = master
		mine.Operator = operator
	}
	return err
}

func (mine *GroupInfo) UpdateAssistant(uid, operator string) error {
	err := nosql.UpdateGroupAssistant(mine.UID, uid, operator)
	if err == nil {
		mine.Assistant = uid
		mine.Operator = operator
	}
	return err
}

func (mine *GroupInfo) UpdateCover(cover, operator string) error {
	err := nosql.UpdateGroupCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *GroupInfo) UpdateLocation(local, operator string) error {
	err := nosql.UpdateGroupLocation(mine.UID, local, operator)
	if err == nil {
		mine.Location = local
		mine.Operator = operator
	}
	return err
}

func (mine *GroupInfo) UpdateAddress(country, province, city, zone, operator string) error {
	addr := nosql.AddressInfo{Country: country, Province: province, City: city, Zone: zone}
	err := nosql.UpdateGroupAddress(mine.UID, operator, addr)
	if err == nil {
		mine.Address = addr
		mine.Operator = operator
	}
	return err
}

func (mine *GroupInfo) HadMember(member string) bool {
	if mine.Master == member || mine.Assistant == member {
		return true
	}
	for i := 0; i < len(mine.members); i += 1 {
		if mine.members[i] == member {
			return true
		}
	}
	return false
}

func (mine *GroupInfo) AllMembers() []string {
	return mine.members
}

func (mine *GroupInfo) AppendMember(member string) error {
	if mine.HadMember(member) {
		return errors.New("the member had existed")
	}
	err := nosql.AppendGroupMember(mine.UID, member)
	if err == nil {
		mine.members = append(mine.members, member)
	}
	return err
}

func (mine *GroupInfo) SubtractMember(member string) error {
	if !mine.HadMember(member) {
		return errors.New("the member not existed")
	}
	err := nosql.SubtractGroupMember(mine.UID, member)
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
