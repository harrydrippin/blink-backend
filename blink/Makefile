SRC_PATH = ./
DIST_PATH = ./

PROTOC = protoc
OPT_I = -I=$(SRC_PATH)
GO_OUT = --go_out=plugins=grpc:$(DIST_PATH)

FILES = ./blink.proto

all:
	make go

go: 
	$(PROTOC) $(OPT_I) $(FILES) $(GO_OUT)