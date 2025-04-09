package main

import (
	"context"
	"file-service/api/proto"
	handldown "file-service/internal/handler/download"
	handlist "file-service/internal/handler/list"
	handlup "file-service/internal/handler/upload"
	domain "file-service/internal/models"
	"file-service/internal/repository/disk"
	"file-service/internal/repository/s3"
	"file-service/internal/usecase/download"
	"file-service/internal/usecase/list"
	"file-service/internal/usecase/upload"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"net/http"
	"os"
)

// fileServiceServer объединяет все контроллеры в один сервер
type fileServiceServer struct {
	proto.UnimplementedFileServiceServer
	uploadCtrl   *handlup.UploadController
	downloadCtrl *handldown.DownloadController
	listCtrl     *handlist.ListController
}

func newFileServiceServer(uploadCtrl *handlup.UploadController, downloadCtrl *handldown.DownloadController, listCtrl *handlist.ListController) *fileServiceServer {
	return &fileServiceServer{
		uploadCtrl:   uploadCtrl,
		downloadCtrl: downloadCtrl,
		listCtrl:     listCtrl,
	}
}

// Реализация методов FileService
func (s *fileServiceServer) UploadFile(stream proto.FileService_UploadFileServer) error {
	return s.uploadCtrl.UploadFile(stream)
}

func (s *fileServiceServer) DownloadFile(req *proto.DownloadRequest, stream proto.FileService_DownloadFileServer) error {
	return s.downloadCtrl.DownloadFile(req, stream)
}

func (s *fileServiceServer) ListFiles(ctx context.Context, req *emptypb.Empty) (*proto.ListResponse, error) {
	return s.listCtrl.ListFiles(ctx, req)
}

func main() {
	// Choose repository implementation
	var repo domain.FileRepository
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "disk" // default value
	}

	switch storageType {
	case "s3":
		s3Repo, err := s3.NewS3Repository("file-service-bucket")
		if err != nil {
			log.Fatalf("Failed to create S3 repository: %v", err)
		}
		repo = s3Repo
	case "disk":
		repo = disk.NewDiskRepository("./storage")
	default:
		log.Fatalf("Unknown storage type: %s", storageType)
	}

	// Initialize use cases
	uploadUC := upload.NewUploadUseCase(repo)
	downloadUC := download.NewDownloadUseCase(repo)
	listUC := list.NewListUseCase(repo)

	// Initialize controllers
	uploadCtrl := handlup.NewUploadController(uploadUC)
	downloadCtrl := handldown.NewDownloadController(downloadUC)
	listCtrl := handlist.NewListController(listUC)

	// Create combined gRPC server
	grpcServer := grpc.NewServer()
	fileServer := newFileServiceServer(uploadCtrl, downloadCtrl, listCtrl)
	proto.RegisterFileServiceServer(grpcServer, fileServer)

	// Start gRPC server in a goroutine
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	go func() {
		log.Println("gRPC Server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Setup gRPC-Gateway
	ctx := context.Background()
	mux := runtime.NewServeMux()

	// Serve Swagger UI and static files
	fs := http.FileServer(http.Dir("./swagger"))
	http.Handle("/swagger/", http.StripPrefix("/swagger", fs))

	// Register gRPC-Gateway handlers
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = proto.RegisterFileServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// Combine gRPC-Gateway and static file server
	http.Handle("/", mux)

	// Start HTTP server for REST/Swagger
	log.Println("HTTP Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
