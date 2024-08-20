package app

import (
	"fmt"
	"net"
	"path/filepath"
	"runtime"

	pb "github.com/Yuelioi/vidor/internal/proto"
	"github.com/joho/godotenv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Port    string
	Service pb.DownloadServiceClient
}

func NewClient() (*Client, error) {

	c := &Client{}

	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.Service = pb.NewDownloadServiceClient(conn)

	return c, nil
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
	env := filepath.Join(filepath.Dir(filename), "..", ".env")

	// Attempt to load the .env file
	err := godotenv.Load(env)
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}
}
