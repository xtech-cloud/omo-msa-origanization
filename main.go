package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/registry/etcdv3/v2"
	proto "github.com/xtech-cloud/omo-msp-organization/proto/organization"
	"io"
	"omo.msa.organization/cache"
	"omo.msa.organization/config"
	"omo.msa.organization/grpc"
	"os"
	"path/filepath"
	"time"
)

var (
	BuildVersion string
	BuildTime    string
	CommitID     string
)

func main() {
	config.Setup()
	err := cache.InitData()
	if err != nil {
		panic(err)
	}
	// New Service
	service := micro.NewService(
		micro.Name("omo.msa.organization"),
		micro.Version(BuildVersion),
		micro.RegisterTTL(time.Second*time.Duration(config.Schema.Service.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Schema.Service.Interval)),
		micro.Address(config.Schema.Service.Address),
	)
	// Initialise service
	service.Init()
	// Register Handler
	_ = proto.RegisterSceneServiceHandler(service.Server(), new(grpc.SceneService))
	_ = proto.RegisterGroupServiceHandler(service.Server(), new(grpc.GroupService))
	_ = proto.RegisterRoomServiceHandler(service.Server(), new(grpc.RoomService))
	_ = proto.RegisterAreaServiceHandler(service.Server(), new(grpc.AreaService))
	_ = proto.RegisterRegionServiceHandler(service.Server(), new(grpc.RegionService))
	_ = proto.RegisterDeviceServiceHandler(service.Server(), new(grpc.DeviceService))

	app, _ := filepath.Abs(os.Args[0])

	logger.Info("-------------------------------------------------------------")
	logger.Info("- Micro Service Agent -> Run")
	logger.Info("-------------------------------------------------------------")
	logger.Infof("- version      : %s", BuildVersion)
	logger.Infof("- application  : %s", app)
	logger.Infof("- md5          : %s", md5hex(app))
	logger.Infof("- build        : %s", BuildTime)
	logger.Infof("- commit       : %s", CommitID)
	logger.Info("-------------------------------------------------------------")
	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

func md5hex(_file string) string {
	h := md5.New()

	f, err := os.Open(_file)
	if err != nil {
		return ""
	}
	defer f.Close()

	io.Copy(h, f)

	return hex.EncodeToString(h.Sum(nil))
}
