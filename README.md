# omo-msa-organization
Micro Service Agent - organization

生成proto:
protoc -I ./grpc/proto --go_out=plugins=grpc:./grpc/proto ./grpc/proto/*.proto

make call
MICRO_REGISTRY=consul micro call omo.msa.organization SceneService.AddOne '{"name":"school-1", "type":1, "cover":"", "master":"111111", "remark":"test-1", "location":"ddd", "operator":"dddd"}'
MICRO_REGISTRY=consul micro call omo.msa.organization SceneService.GetOne '{"uid":"5f0fbf01b780dd269d83eb79"}'
MICRO_REGISTRY=consul micro call omo.msa.organization SceneService.RemoveOne '{"uid":"5f0fbf01b780dd269d83eb79"}'