.PHONY: all

mipsle:
	GOOS="linux" GOARCH="mipsle" go build -o build/speedmon_mipsle

mips64:
	GOOS="linux" GOARCH="mips64" go build -o build/speedmon_mips64

all: mipsle mips64

install: all
	hub release create -a build/speedmon_mips64 -a build/speedmon_mipsle $(date +%Y%m%d)
