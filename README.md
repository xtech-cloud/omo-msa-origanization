# omo-msa-user
Micro Service Agent - organization

生成proto:
protoc -I ./grpc/proto --go_out=plugins=grpc:./grpc/proto ./grpc/proto/*.proto
