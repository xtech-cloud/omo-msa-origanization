go install omo.msa.organization
mkdir _build
mkdir _build/bin

cp -rf /root/go/bin/omo.msa.organization _build/bin/
cp -rf conf _build/
cd _build
tar -zcf msa.organization.tar.gz ./*
mv msa.organization.tar.gz ../
cd ../
rm -rf _build
