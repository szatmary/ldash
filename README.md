# Instruction

## 1. setup environment

./build_ldash_server.sh

## 2. launch the server by "go run main.go $streamfolder"

for example, go run main.go "/home/lei/www"
After the server is running, following upload endpoints will be open: 
Manifest: "http://{yourdomain}:8080/ldash/upload/{filepath}"
Segment: "http://{yourdomain}:8080/ldash/upload/{filepath}"

You can test those endpoint with 
curl http://localhost:8080/ldash/upload/2550/720p/123.ts --upload-file $anyfile

## 3. check download
you can also check the download endpoint
"http://{yourdomain}:8080/ldash/download/{filepath}"

