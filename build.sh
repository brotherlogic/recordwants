protoc --proto_path ../../../ -I=./proto --go_out=plugins=grpc:./proto proto/recordwants.proto
mv proto/github.com/brotherlogic/recordwants/proto/* ./proto
