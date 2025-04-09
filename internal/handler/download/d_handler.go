package download

import (
	"file-service/api/proto"
	"file-service/internal/usecase/download"
	"file-service/internal/utils"
)

// DownloadController ...
type DownloadController struct {
	useCase *download.DownloadUseCase
	proto.UnimplementedFileServiceServer
}

// NewDownloadController ...
func NewDownloadController(useCase *download.DownloadUseCase) *DownloadController {
	return &DownloadController{useCase: useCase}
}

// DownloadFile ...
func (c *DownloadController) DownloadFile(req *proto.DownloadRequest, stream proto.FileService_DownloadFileServer) error {
	ctx, cancel := utils.NewWithCancel(stream.Context())
	defer cancel()

	file, err := c.useCase.Download(ctx, req.Name)
	if err != nil {
		return err
	}

	chunkSize := 64 * 1024
	for currentByte := 0; currentByte < len(file.Data); currentByte += chunkSize {
		end := currentByte + chunkSize
		if end > len(file.Data) {
			end = len(file.Data)
		}

		if err := stream.Send(&proto.DownloadResponse{
			Chunk: file.Data[currentByte:end],
		}); err != nil {
			return err
		}
	}

	return nil
}
