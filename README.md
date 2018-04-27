# Instruction

## 1. setup environment

./build_lhls_server.sh

## 2. launch the server by "go run main.go $streamfolder"

for example, go run main.go "/home/lei/streamlinelhls/www"
After the server is running, following upload endpoints will be open: 
Manifest: "http://{yourdomain}:8888/upload/{videoid}/{resolution}/{manifest}"
Segment: "http://{yourdomain}:8888/upload/{videoid}/{resolution}/{segment}"

You can test those endpoint with 
curl http://localhost:8888/upload/2550/720p/123.ts --upload-file $anyfile

## 3. open the player, type "http://{yourdomain}/lhls/{videoid}/index.html"

Or you can also check the stream in safari, type "http://{yourdomain}/lhls/{videoid}/{manifestname}", for example
"http://{yourdomain}:8888/lhls/2550/2550.m3u8"

Beside the main the manifest file, you can also check one specified resolution by 
"http://{yourdomain}/lhls/{videoid}/{resolution}/{manifestname}"

