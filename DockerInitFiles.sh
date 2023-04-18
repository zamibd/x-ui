#!/bin/sh
if [ $1 == "amd64" ]; then
    ARCH="64";
    FNAME="amd64";
elif [ $1 == "arm64" ]; then
    ARCH="arm64-v8a"
    FNAME="arm64";
else
    ARCH="64";
    FNAME="amd64";
fi
mkdir -p build/bin
wget "https://github.com/XTLS/Xray-core/releases/download/v1.8.1/Xray-linux-${ARCH}.zip"
wget "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat" -o build/bin/geoip.dat
wget "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat" -o build/bin/geosite.dat
mv xray "build/bin/xray-linux-${FNAME}"
mv main build/x-ui