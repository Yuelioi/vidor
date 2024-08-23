package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"

	pb "github.com/Yuelioi/vidor/internal/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Plugin struct {
	PluginConfig
	Type     string `json:"type"` // System/ThirdPart
	Location string `json:"location"`
	Enable   bool
	State    int // 1.运行中 2.运行但是通信失败 3.未启动
	Port     int
	PID      int
	service  pb.DownloadServiceClient
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

	p.service = pb.NewDownloadServiceClient(conn)

	return p, nil
}

func InitPlugin(p *Plugin, c *Config) (*Plugin, error) {

	// 插件设置
	LoadEnv()
	value := os.Getenv("SESSDATA")
	ctx := metadata.AppendToOutgoingContext(context.Background(), "plugin.sessdata", value, "host", "vidor")

	_, err := p.service.Init(ctx, nil)
	if err != nil {
		return nil, err
	}

	logger.Infof("已成功加载插件%s", p.Name)
	return p, err
}

// 停止插件
func StopPlugin(p *Plugin) (*Plugin, error) {
	if p == nil {
		return nil, fmt.Errorf("插件不存在")
	}

	_, err := p.service.Shutdown(context.TODO(), &emptypb.Empty{})
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
