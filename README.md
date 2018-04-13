# Speedmon

This is simple speedtest monitoring tool that is forked from https://github.com/surol/speedtest-cli.
Only addittions is reporting, which instead of reporting results to console, sends them to CloudWatch.

## Compile

	GOOS="linux" GOARCH="mipsle" go build -o build/speedmon  # For EdgeRouter X
	GOOS="linux" GOARCH="mips64" go build -o build/speedmon  # For Unifi USG

## Setup

Create AWS credential file as `$HOME/.aws/credentials` with content:

```
[default]
region = eu-west-1
aws_access_key_id = YOURKEY
aws_secret_access_key = YOURSECRET
```

Run simply via `./speedmon`
