package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudwatch"

    "fmt"
    "log"
    "os/exec"
    "strings"
    "strconv"
    "time"
)

func main() {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    
    // Create new cloudwatch client.
    svc := cloudwatch.New(sess)


    out, err := exec.Command("./speedtest.py", "--no-pre-allocate", "--csv").Output()
    if err != nil {
        log.Fatal(err)
    }

    s := strings.Split(string(out), ",")
    // 0		  1		   2			3		   4		 5	   6	     7	     8	    9
    // server_id   := s[0]
    // sponsor 	   := s[1]
    server_name  := s[2]
    timestamp, _ := time.Parse(time.RFC3339, s[3])
    // distance 	 := s[4]
    ping, _      := strconv.ParseFloat(s[5], 64)
    download, _  := strconv.ParseFloat(s[6], 64)
    upload, _	 := strconv.ParseFloat(s[7], 64)
    // share 		 := s[8]
    // ip_address 	 := s[9]

    fmt.Printf("DL: %s, UL: %s\n", download, upload)

    result, err := svc.PutMetricData(&cloudwatch.PutMetricDataInput{
        MetricData: []*cloudwatch.MetricDatum{
            &cloudwatch.MetricDatum{
                MetricName: aws.String("Ping"),
                Timestamp:  aws.Time(timestamp),
                Unit:       aws.String(cloudwatch.StandardUnitMilliseconds),
                Value:      aws.Float64(ping),
                Dimensions: []*cloudwatch.Dimension{
                    &cloudwatch.Dimension{
                        Name:  aws.String("Server"),
                        Value: aws.String(server_name),
                    },
                },
            },
            &cloudwatch.MetricDatum{
                MetricName: aws.String("DownloadSpeed"),
                Timestamp:  aws.Time(timestamp),
                Unit:       aws.String(cloudwatch.StandardUnitBitsSecond),
                Value:      aws.Float64(download),
                Dimensions: []*cloudwatch.Dimension{
                    &cloudwatch.Dimension{
                        Name:  aws.String("Server"),
                        Value: aws.String(server_name),
                    },
                },
            },
            &cloudwatch.MetricDatum{
                MetricName: aws.String("UploadSpeed"),
                Timestamp:  aws.Time(timestamp),
                Unit:       aws.String(cloudwatch.StandardUnitBitsSecond),
                Value:      aws.Float64(upload),
                Dimensions: []*cloudwatch.Dimension{
                    &cloudwatch.Dimension{
                        Name:  aws.String("Server"),
                        Value: aws.String(server_name),
                    },
                },
            },
            &cloudwatch.MetricDatum{
                MetricName: aws.String("Ping"),
                Timestamp:  aws.Time(timestamp),
                Unit:       aws.String(cloudwatch.StandardUnitMilliseconds),
                Value:      aws.Float64(ping),
                Dimensions: []*cloudwatch.Dimension{
                    &cloudwatch.Dimension{
                        Name:  aws.String("Server"),
                        Value: aws.String("Common"),
                    },
                },
            },
            &cloudwatch.MetricDatum{
                MetricName: aws.String("DownloadSpeed"),
                Timestamp:  aws.Time(timestamp),
                Unit:       aws.String(cloudwatch.StandardUnitBitsSecond),
                Value:      aws.Float64(download),
                Dimensions: []*cloudwatch.Dimension{
                    &cloudwatch.Dimension{
                        Name:  aws.String("Server"),
                        Value: aws.String("Common"),
                    },
                },
            },
            &cloudwatch.MetricDatum{
                MetricName: aws.String("UploadSpeed"),
                Timestamp:  aws.Time(timestamp),
                Unit:       aws.String(cloudwatch.StandardUnitBitsSecond),
                Value:      aws.Float64(upload),
                Dimensions: []*cloudwatch.Dimension{
                    &cloudwatch.Dimension{
                        Name:  aws.String("Server"),
                        Value: aws.String("Common"),
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