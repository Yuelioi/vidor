package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	pb "github.com/Yuelioi/vidor/internal/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Plugin struct {
	Name     string
	Type     string // System/ThirdPart
	Location string
	Enable   bool
	Port     int
	PID      int
	Regexs   []*regexp.Regexp
	Service  pb.DownloadServiceClient
}

// 加载插件基础信息
func LoadPlugin(name string, location string, _type string) (*Plugin, error) {
	p := &Plugin{
		Name:     name,
		Type:     _type,
		Location: location,
	}

	return p, nil
}

// 初始化插件 运行插件
func RunPlugin(p *Plugin) (*Plugin, error) {

	if p == nil {
		return nil, fmt.Errorf("插件不存在")
	}
	// availablePort, err := getAvailablePort()
	// exec.Command("server.exe", "--port", strconv.Itoa(availablePort))

	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	p.Service = pb.NewDownloadServiceClient(conn)
	// 插件设置
	LoadEnv()
	value := os.Getenv("SESSDATA")
	ctx := metadata.AppendToOutgoingContext(context.Background(), "plugin.sessdata", value, "host", "vidor")

	_, err = p.Service.Init(ctx, nil)
	if err != nil {
		return nil, err
	}

	logger.Infof("已成功加载插件%s", p.Name)
	return p, nil
}

// 停止插件
func StopPlugin(p *Plugin) (*Plugin, error) {
	if p == nil {
		return nil, fmt.Errorf("插件不存在")
	}

	_, err := p.Service.Shutdown(context.TODO(), &emptypb.Empty{})
	return nil, err
}

func getAvailablePort() (int, error) {
	// 监听 "localhost:0" 让系统分配一个可用端口
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("failed to find an available port: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	return port, nil
}

func LoadEnv() {

	_, filename, _, _ := runtime.Caller(0)
	env := filepath.Join(filepath.Dir(filename), "..", "..", ".env")

	// Attempt to load the .env file
	err := godotenv.Load(env)
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}
}
