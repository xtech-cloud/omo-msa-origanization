package cache

import (
	"errors"
	"omo.msa.organization/proxy/nosql"
)

type RegionInfo struct {
	baseInfo
	Entity string
	Remark string
	Scene string
	Parent string
	Master string
	Code string
	Location string
	Address nosql.AddressInfo
	Members []string
}

func (mine *cacheContext)GetRegion(uid string) (*RegionInfo,error) {
	db,err := nosql.GetRegion(uid)
	if err == nil {
		info := new(RegionInfo)
		info.initInfo(db)
		return info,nil
	}
	return nil, err
}

func (mine *cacheContext)GetRegionsByScene(scene string) []*RegionInfo {
	list := make([]*RegionInfo, 0, 10)
	dbs,err := nosql.GetRegionsByScene(scene)
	if err == nil {
		for _, db := range dbs {
			info := new(RegionInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext)GetRegionsByParent(parent string) []*RegionInfo {
	list := make([]*RegionInfo, 0, 10)
	dbs,err := nosql.GetRegionsByParent(parent)
	if err == nil {
		for _, db := range dbs {
			info := new(RegionInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *RegionInfo)initInfo(db *nosql.Region)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Remark = db.Remark
	mine.Code = db.Code
	mine.Scene = db.Scene
	mine.Parent = db.Parent
	mine.Master = db.Master
	mine.Entity = db.Entity
	mine.Location = db.Location
	mine.Address = db.Address
	mine.Members = db.Members
	if mine.Members == nil {
		mine.Members = make([]string, 0, 1)
	}
}

func (mine *RegionInfo)UpdateBase(name, remark, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	err := nosql.UpdateRegionBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *RegionInfo)Remove(operator string) error {
	if mine.HadChildren() {
		return errors.New("the region had children")
	}
	return nosql.RemoveRegion(mine.UID, operator)
}

func (mine *RegionInfo)HadChildren() bool {
	list := cacheCtx.GetRegionsByParent(mine.UID)
	if len(list) > 0 {
		return true
	}
	return false
}

func (mine *RegionInfo)UpdateMaster(master, operator string) error {
	err := nosql.UpdateRegionMaster(mine.UID, master, operator)
	if err == nil {
		mine.Master = master
		mine.Operator = operator
	}
	return err
}

func (mine *RegionInfo)UpdateParent(parent, operator string) error {
	err := nosql.UpdateRegionParent(mine.UID, parent, operator)
	if err == nil {
		mine.Parent = parent
		mine.Operator = operator
	}
	return err
}

func (mine *RegionInfo)UpdateLocation(local, operator string) error {
	err := nosql.UpdateRegionLocation(mine.UID, local, operator)
	if err == nil {
		mine.Location = local
		mine.Operator = operator
	}
	return err
}

func (mine *RegionInfo)UpdateEntity(entity, operator string) error {
	err := nosql.UpdateRegionEntity(mine.UID, entity, operator)
	if err == nil {
		mine.Entity = entity
		mine.Operator = operator
	}
	return err
}

func (mine *RegionInfo)UpdateAddress(country, province, city, zone, operator string) error {
	addr := nosql.AddressInfo{Country: country, Province: province, City: city, Zone: zone}
	err := nosql.UpdateRegionAddress(mine.UID, operator, addr)
	if err == nil {
		mine.Address = addr
		mine.Operator = operator
	}
	return err
}

func (mine *RegionInfo)HadMember(member string) bool {
	if mine.Master == member  {
		return true
	}
	for i := 0;i < len(mine.Members);i += 1 {
		if mine.Members[i] == member {
			return true
		}
	}
	return false
}

func (mine *RegionInfo)AppendMember(member string) error {
	if mine.HadMember(member){
		return nil
	}
	err := nosql.AppendRegionMember(mine.UID, member)
	if err == nil {
		mine.Members = append(mine.Members, member)
	}
	return err
}

func (mine *RegionInfo)SubtractMember(member string) error {
	if !mine.HadMember(member){
		return nil
	}
	err := nosql.SubtractRegionMember(mine.UID, member)
	if err == nil {
		for i := 0;i < len(mine.Members);i += 1 {
			if mine.Members[i] == member {
				if i == len(mine.Members) - 1 {
					mine.Members = append(mine.Members[:i])
				}else{
					mine.Members = append(mine.Members[:i], mine.Members[i+1:]...)
				}

				break
			}
		}
	}
	return err
}

