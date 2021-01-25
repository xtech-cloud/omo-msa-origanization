export GO111MODULE=on
export GOSUMDB=off
export GOPROXY=https://mirrors.aliyun.com/goproxy/
go install omo.msa.organization
mkdir _build
mkdir _build/bin

cp -rf /root/go/bin/omo.msa.organization _build/bin/
cd _build
tar -zcf msa.organization.tar.gz ./*
mv msa.organization.tar.gz ../
cd ../
rm -rf _build
