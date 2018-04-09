# Speedmon

Simple tool for running `speedtest.py` and then reporting results to CloudWatch.
Designed for running via Cron on EdgerouterX to monitor home ISP network speed.
As it happens, my provider advertises 350M connection, but reality is far from it...

## Compile


	GOOS="linux" GOARCH="mipsle" go build -o speedmon

## Setup

Create AWS credential file as `$HOME/.aws/credentials` with content:

```
[default]
region = eu-west-1
aws_access_key_id = YOURKEY
aws_secret_access_key = YOURSECRET
```

Run simply via `./speedmon`
