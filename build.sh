rm -rf ./releases/*
cd web
npm run build
cd ../

# build for mac
cd server
GOOS=darwin GOARCH=amd64 go build 
cd ../releases
mkdir gshark_darwin_amd64
cd gshark_darwin_amd64
mv ../../server/gshark .
cp -rf ../../server/resource .
cp ../../server/config-temp.yaml config.yaml
cd ../../
cp -rf ./web/dist ./releases/gshark_darwin_amd64
7z a -r ./releases/gshark_darwin_amd64.zip ./releases/gshark_darwin_amd64/

# build for windows
cd server
GOOS=windows GOARCH=amd64 go build
cd ../releases
mkdir gshark_windows_amd64
cd gshark_windows_amd64
mv ../../server/gshark.exe .
cp -rf ../../server/resource .
cp ../../server/config-temp.yaml config.yaml
cd ../../
cp -rf ./web/dist ./releases/gshark_windows_amd64
7z a -r ./releases/gshark_windows_amd64.zip ./releases/gshark_windows_amd64/

# build for linux
cd server
GOOS=linux GOARCH=amd64 go build -o gshark
cd ../releases
mkdir gshark_linux_amd64
cd gshark_linux_amd64
mv ../../server/gshark .
cp -rf ../../server/resource .
cp ../../server/config-temp.yaml config.yaml
cd ../../
cp -rf ./web/dist ./releases/gshark_linux_amd64
7z a -r ./releases/gshark_linux_amd64.zip ./releases/gshark_linux_amd64


rm -rf ./releases/gshark*/