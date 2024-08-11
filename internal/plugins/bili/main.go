package main

import (
	pb "bilibili/proto"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedSearchServiceServer
	grpcServer *grpc.Server
}

func (s *Server) Show(ctx context.Context, req *pb.RequestShow) (*pb.PlaylistInfo, error) {

	downloader := New(context.Background(), req.Entries)
	pi, _ := downloader.Show(req.Url)

	return pi, nil
}

func (s *Server) Parse(ctx context.Context, req *pb.RequestParse) (*pb.PlaylistInfo, error) {

	downloader := New(context.Background(), req.Entries)
	newPi, _ := downloader.Parse(context.Background(), req.Playlist)
	return newPi, nil
}

func (s *Server) Download(stream grpc.BidiStreamingServer[pb.PlaylistInfo, pb.DownloadResponse]) error {

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Printf("req.URL: %v\n", req.URL)

		// 模拟下载过程
		for i := int64(0); i <= 100; i += 10 {
			res := &pb.DownloadResponse{
				Status:   "Downloading",
				Progress: i,
			}
			if err := stream.Send(res); err != nil {
				return err
			}
			time.Sleep(time.Second) // 模拟下载进度
		}

		res := &pb.DownloadResponse{
			Status:   "Completed",
			Progress: 100,
		}
		if err := stream.Send(res); err != nil {
			return err
		}

	}
}

func (s *Server) Stop(ctx context.Context, req *pb.RequestStop) (*pb.ResponseStop, error) {
	fmt.Printf("req.Id: %v\n", req.Id)

	return &pb.ResponseStop{
		Id:    req.Id,
		State: "Stopped",
	}, nil
}

func (s *Server) Shutdown(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	go func() {
		log.Println("Shutting down GRPC server...")
		s.grpcServer.GracefulStop()
	}()
	return &pb.Empty{}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	// 创建服务并注册
	grpcServer := grpc.NewServer()
	pb.RegisterSearchServiceServer(grpcServer, &Server{})

	// 启动
	log.Println("GRPC 启动")
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
		return
	}
}
