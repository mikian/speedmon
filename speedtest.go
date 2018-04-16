package main

import (
	"github.com/surol/speedtest-cli/speedtest"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"fmt"
	"log"
	"os"
	"time"
	"flag"
)

func version() {
	fmt.Print(speedtest.Version)
}

func usage() {
	fmt.Fprint(os.Stderr, "Command line interface for testing internet bandwidth using speedtest.net.\n\n")
	flag.PrintDefaults()
}

func main() {
	opts := speedtest.ParseOpts()

	switch {
	case opts.Help:
		usage()
		return
	case opts.Version:
		version()
		return
	}

	client := speedtest.NewClient(opts)

	if opts.List {
		servers, err := client.AllServers()
		if err != nil {
			log.Fatalf("Failed to load server list: %v\n", err)
		}
		fmt.Println(servers)
		return
	}

	config, err := client.Config()
	if err != nil {
		log.Fatal(err)
	}

	client.Log("Testing from %s (%s)...\n", config.Client.ISP, config.Client.IP)

	server := selectServer(opts, client);

	downloadSpeed := server.DownloadSpeed()
	// reportSpeed(opts, "Download", downloadSpeed)

	uploadSpeed := server.UploadSpeed()
	// reportSpeed(opts, "Upload", uploadSpeed)

	reportMeasurement(server, downloadSpeed, uploadSpeed)
}

func reportSpeed(opts *speedtest.Opts, prefix string, speed int) {
	if opts.SpeedInBytes {
		fmt.Printf("%s: %.2f MiB/s\n", prefix, float64(speed) / (1 << 20))
	} else {
		fmt.Printf("%s: %.2f Mib/s\n", prefix, float64(speed) / (1 << 17))
	}
}

func reportMeasurement(server *speedtest.Server, download int, upload int) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
	}))

	// Create new cloudwatch client.
	svc := cloudwatch.New(sess)

	fmt.Printf("Download: %.2f Mib/s\tUpload: %.2f Mib/s\tPing: %dms\n",
		float64(download) / (1 << 17),
		float64(upload) / (1 << 17),
		server.Latency / time.Millisecond)

	timestamp := time.Now()
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	result, err := svc.PutMetricData(&cloudwatch.PutMetricDataInput{
			MetricData: []*cloudwatch.MetricDatum{
					&cloudwatch.MetricDatum{
							MetricName: aws.String("Ping"),
							Timestamp:  aws.Time(timestamp),
							Unit:       aws.String(cloudwatch.StandardUnitMilliseconds),
							Value:      aws.Float64(float64(server.Latency / time.Millisecond)),
							Dimensions: []*cloudwatch.Dimension{
									&cloudwatch.Dimension{
											Name:  aws.String("ServerName"),
											Value: aws.String(server.Name),
									},
							},
					},
					&cloudwatch.MetricDatum{
							MetricName: aws.String("DownloadSpeed"),
							Timestamp:  aws.Time(timestamp),
							Unit:       aws.String(cloudwatch.StandardUnitMegabitsSecond),
							Value:      aws.Float64(float64(download) / (1 << 17)),
							Dimensions: []*cloudwatch.Dimension{
									&cloudwatch.Dimension{
											Name:  aws.String("ServerName"),
											Value: aws.String(server.Name),
									},
							},
					},
					&cloudwatch.MetricDatum{
							MetricName: aws.String("UploadSpeed"),
							Timestamp:  aws.Time(timestamp),
							Unit:       aws.String(cloudwatch.StandardUnitMegabitsSecond),
							Value:      aws.Float64(float64(upload) / (1 << 17)),
							Dimensions: []*cloudwatch.Dimension{
									&cloudwatch.Dimension{
											Name:  aws.String("ServerName"),
											Value: aws.String(server.Name),
									},
							},
					},
					&cloudwatch.MetricDatum{
							MetricName: aws.String("Ping"),
							Timestamp:  aws.Time(timestamp),
							Unit:       aws.String(cloudwatch.StandardUnitMilliseconds),
							Value:      aws.Float64(float64(server.Latency / time.Millisecond)),
							Dimensions: []*cloudwatch.Dimension{
									&cloudwatch.Dimension{
											Name:  aws.String("Host"),
											Value: aws.String(hostname),
									},
							},
					},
					&cloudwatch.MetricDatum{
							MetricName: aws.String("DownloadSpeed"),
							Timestamp:  aws.Time(timestamp),
							Unit:       aws.String(cloudwatch.StandardUnitMegabitsSecond),
							Value:      aws.Float64(float64(download) / (1 << 17)),
							Dimensions: []*cloudwatch.Dimension{
									&cloudwatch.Dimension{
											Name:  aws.String("Host"),
											Value: aws.String(hostname),
									},
							},
					},
					&cloudwatch.MetricDatum{
							MetricName: aws.String("UploadSpeed"),
							Timestamp:  aws.Time(timestamp),
							Unit:       aws.String(cloudwatch.StandardUnitMegabitsSecond),
							Value:      aws.Float64(float64(upload) / (1 << 17)),
							Dimensions: []*cloudwatch.Dimension{
									&cloudwatch.Dimension{
											Name:  aws.String("Host"),
											Value: aws.String(hostname),
									},
							},
					},
			},
			Namespace: aws.String("Home"),
	})
	if err != nil {
			fmt.Println("Error", err)
			return
	}

	fmt.Println("Success", result)
}

func selectServer(opts *speedtest.Opts, client *speedtest.Client) (selected *speedtest.Server) {
	if opts.Server != 0 {
		servers, err := client.AllServers()
		if err != nil {
			log.Fatal("Failed to load server list: %v\n", err)
			return nil
		}
		selected = servers.Find(opts.Server)
		if selected == nil {
			log.Fatalf("Server not found: %d\n", opts.Server)
			return nil
		}
		selected.MeasureLatency(speedtest.DefaultLatencyMeasureTimes, speedtest.DefaultErrorLatency)
	} else {
		servers, err := client.ClosestServers()
		if err != nil {
			log.Fatal("Failed to load server list: %v\n", err)
			return nil
		}
		selected = servers.MeasureLatencies(
			speedtest.DefaultLatencyMeasureTimes,
			speedtest.DefaultErrorLatency).First()
	}

	if opts.Quiet {
		log.Printf("Ping: %d ms\n", selected.Latency / time.Millisecond)
	} else {
		client.Log("Hosted by %s (%s) [%.2f km]: %d ms\n",
			selected.Sponsor,
			selected.Name,
			selected.Distance,
			selected.Latency / time.Millisecond)
	}

	return selected
}
