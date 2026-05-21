package services

import (
	"fmt"
	"net"

	pb "mangahub/gen/manga/proto"
	grpcserver "mangahub/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartGRPC() {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		fmt.Println(err)
		return
	}

	s := grpc.NewServer()

	pb.RegisterMangaServiceServer(s, &grpcserver.MangaServer{})

	reflection.Register(s)

	fmt.Println("⚡ gRPC server running at :50051")

	s.Serve(lis)
}
