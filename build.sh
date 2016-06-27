#!/bin/bash
begin=$(date +%s)
echo "Building for mac"
env GOOS=darwin GOARCH=amd64 go build -v -o easy_download-mac easy_download.go &


echo "Building for linux"
env GOOS=linux GOARCH=i386 go build -v -o easy_download-linux easy_download.go &


echo "Building for windows"
env GOOS=windows GOARCH=amd64 go build -v -o easy_download-win.exe easy_download.go &

wait

end=$(date +%s)

total=$((end-begin))

echo "Built-in "$total"s"
