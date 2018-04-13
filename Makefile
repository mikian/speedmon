.PHONY:

erx:
	GOOS="linux" GOARCH="mipsle" go build -o build/speedmon_mipsle

usg:
	GOOS="linux" GOARCH="mips64" go build -o build/speedmon_mips64

all: erx usg
