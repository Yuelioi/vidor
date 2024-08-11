package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "bilibili/proto"

	"testing"
)

type PerRPCCredentials interface {
	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
	RequireTransportSecurity() bool
}

// 自定义 Token 认证结构体 需要实现PerRPCCredentials接口
type TokenAuth struct {
	Token string
}

func (ta *TokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"Authorization": "Bearer " + ta.Token}, nil
}

func (*TokenAuth) RequireTransportSecurity() bool {
	return false
}

// 创建 Token 认证实例
func NewTokenAuth(token string) *TokenAuth {
	return &TokenAuth{Token: token}
}

var (
	testUrl      = "https://www.bilibili.com/video/BV1BBvfeyEbA"
	testPagesUrl = "https://www.bilibili.com/video/BV1Ni421X7PT"
)

func TestShow(t *testing.T) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithPerRPCCredentials(NewTokenAuth("token 123456")))

	// 连接并加密
	conn, err := grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("链接失败%+v", err)
	}

	defer conn.Close()

	// 建立连接
	client := pb.NewSearchServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用远程服务的方法
	r, err := client.Show(ctx, &pb.RequestShow{Url: testUrl})
	if err != nil {
		log.Fatalf("could not get playback info: %v", err)
	}
	log.Printf("Playback Info: %v", r)

}
func TestParse(t *testing.T) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithPerRPCCredentials(NewTokenAuth("token 123456")))

	// 连接并加密
	conn, err := grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("链接失败%+v", err)
	}

	defer conn.Close()

	// 建立连接
	client := pb.NewSearchServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用远程服务的方法
	res, err := client.Show(ctx, &pb.RequestShow{Url: testPagesUrl})
	if err != nil {
		log.Fatalf("could not get playback info: %v", err)
	}

	res2, _ := client.Parse(ctx, &pb.RequestParse{Playlist: res})
	if err != nil {
		log.Fatalf("could not get playback info: %v", err)
	}
	log.Printf("Playback Info: %v", res2)

}
func TestDownload(t *testing.T) {

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithPerRPCCredentials(NewTokenAuth("token 123456")))

	// 连接并加密
	conn, err := grpc.NewClient("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("链接失败%+v", err)
	}

	defer conn.Close()

	// 建立连接
	client := pb.NewSearchServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用远程服务的方法
	r, err := client.Show(ctx, &pb.RequestShow{Url: "https://www.bilibili.com/video/BV1BBvfeyEbA"})
	if err != nil {
		log.Fatalf("could not get playback info: %v", err)
	}
	log.Printf("Playback Info: %v", r)

	// 创建下载流
	stream, err := client.Download(context.Background())
	if err != nil {
		log.Fatalf("could not create download stream: %v", err)
	}

	// 发送下载请求
	if err := stream.Send(&pb.PlaylistInfo{
		URL: "https://example.com/video",
	}); err != nil {
		log.Fatalf("could not send download request: %v", err)
	}

	// 接收下载进度
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Fatalf("could not receive download progress: %v", err)

			break
		}
		if err != nil {
			log.Fatalf("could not receive download progress: %v", err)
		}
		log.Printf("Download progress: %v%%", res.Progress)
	}

}
