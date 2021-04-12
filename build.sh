rm -rf ./releases/*
cd web
npm run build
cd ../

# build for mac
cd server
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build
cd ../releases
mkdir gshark_darwin_amd64
cd gshark_darwin_amd64
cp ../../server/gshark .
cp -rf ../../server/resource .
cp ../../server/config.yaml .
cd ../../
cp -rf ./web/dist ./releases/gshark_darwin_amd64

# build for windows
cd server
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
cd ../releases
mkdir gshark_windows_amd64
cd gshark_windows_amd64
cp ../../server/gshark.exe .
cp -rf ../../server/resource .
cp ../../server/config.yaml .
cd ../../
cp -rf ./web/dist ./releases/gshark_windows_amd64

# build for linux
cd server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
cd ../releases
mkdir gshark_linux_amd64
cd gshark_linux_amd64
cp ../../server/gshark .
cp -rf ../../server/resource .
cp ../../server/config.yaml .
cd ../../
cp -rf ./web/dist ./releases/gshark_linux_amd64

7z a -r ./releases/gshark_windows_amd64.zip ./releases/gshark_windows_amd64/
7z a -r ./releases/gshark_darwin_amd64.zip ./releases/gshark_darwin_amd64/
7z a -r ./releases/gshark_linux_amd64.zip ./releases/gshark_linux_amd64

rm -rf ./releases/gshark*/