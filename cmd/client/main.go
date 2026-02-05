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
	port   string
	n      int
	m      int
)

func init() {
	flag.StringVar(&metric, "show", "la", "Show metrics. Possible values: la|cpu|disk|fs")
	flag.StringVar(&port, "port", ":8000", "Choose listen port. For example: :8000")
	flag.IntVar(&n, "n", 1, "Send stats every N seconds")
	flag.IntVar(&m, "m", 1, "Send stats for last M seconds")
}

func main() {
	flag.Parse()

	var err error
	switch metric {
	case "la":
		err = runClient(printHeaderLA, printLA)
	case "cpu":
		err = runClient(printHeaderCPU, printCPU)
	case "disk":
		err = runClient(printHeaderDisk, printDisks)
	case "fs":
		err = runClient(printHeaderFS, printFS)
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

	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func printHeaderLA() {
	fmt.Println("Load Average")
	fmt.Println("  time   | load1 | load5 | load15")
}

func printLA(stats *protobuf.SystemStatistics) {
	data := stats.LoadAvg
	if data != nil {
		fmt.Printf("%s | %5.2f | %5.2f | %5.2f\n", formatTime(stats), data.Load1, data.Load5, data.Load15)
	} else {
		fmt.Printf("%s |   -   |   -   |   -\n", formatTime(stats))
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

func printHeaderDisk() {
	fmt.Println("Load Disks")
	fmt.Println("  time   |   tps    |   read   |  write   |  name")
}

func printDisks(stats *protobuf.SystemStatistics) {
	fmt.Printf("%s |          |          |          |\n", formatTime(stats))
	for _, disk := range stats.LoadDisks {
		fmt.Printf("         | %8.2f | %8.2f | %8.2f | %s\n", disk.Tps, disk.KBRead, disk.KBWrite, disk.Name)
	}
}

func printHeaderFS() {
	fmt.Println("Used File Systems")
	fmt.Println("  time   | use%  | IUse% |  path")
}

func printFS(stats *protobuf.SystemStatistics) {
	fmt.Printf("%s |       |       |\n", formatTime(stats))
	for _, fs := range stats.UsedFs {
		fmt.Printf("         | %5.2f | %5.2f | %s\n", fs.UsedSpace, fs.UsedInode, fs.Path)
	}
}

func formatTime(stats *protobuf.SystemStatistics) string {
	return stats.Time.AsTime().Format("15:04:05")
}
