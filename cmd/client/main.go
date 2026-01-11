package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"

	protobuf "github.com/Av4JhG/system-monitoring/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	metric string
	n      int
	m      int
)

func init() {
	flag.StringVar(&metric, "show", "la", "Show metrics. Possible values: la|cpu|disk|fs")
	flag.IntVar(&n, "n", 1, "Send stats every N seconds")
	flag.IntVar(&m, "m", 1, "Send stats for last M seconds")
}

func main() {
	flag.Parse()

	var err error
	switch metric {
	case "la":

	case "cpu":
		err = runClient(printHeaderCPU, printCPU)
	case "disk":

	case "fs":

	default:
		flag.Usage()
	}

	if err != nil {
		log.Print(err)
	}
}

type (
	printHeader func()
	printStats  func(stats *protobuf.SystemStatistics)
)

func runClient(ph printHeader, ps printStats) error {
	ph()

	conn, err := grpc.NewClient(":8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx := context.Background()

	client := protobuf.NewSmClient(conn)
	req := &protobuf.SystemStatisticsRequest{
		N: int32(n),
		M: int32(m),
	}
	reqClient, err := client.GetStats(ctx, req)
	if err != nil {
		return fmt.Errorf("client request fail: %w", err)
	}

	for {
		stats, err := reqClient.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		ps(stats)
	}
}

func printHeaderCPU() {
	fmt.Println("Load CPU")
	fmt.Println("  time   | user  | system| idle")
}

func printCPU(stats *protobuf.SystemStatistics) {
	data := stats.Cpu
	if data != nil {
		fmt.Printf("%s | %5.2f | %5.2f | %5.2f\n", formatTime(stats), data.User, data.System, data.Idle)
	} else {
		fmt.Printf("%s |   -   |   -   |   -\n", formatTime(stats))
	}
}

func formatTime(stats *protobuf.SystemStatistics) string {
	return stats.Time.AsTime().Format("15:04:05")
}
