PB_PROTO_FILES=$(shell find pb -name *.proto)

.PHONY: pb
# generate pb proto
pb:
	protoc --proto_path=./pb \
 	       --go_out=paths=source_relative:./pb \
	       $(PB_PROTO_FILES)


.DEFAULT_GOAL := pb
