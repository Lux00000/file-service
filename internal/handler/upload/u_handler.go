package upload

import (
	"file-service/api/proto"
	domain "file-service/internal/models"
	"file-service/internal/usecase/upload"
	"file-service/internal/utils"
	"io"
)

type UploadController struct {
	useCase *upload.UploadUseCase
	proto.UnimplementedFileServiceServer
}

func NewUploadController(useCase *upload.UploadUseCase) *UploadController {
	return &UploadController{useCase: useCase}
}

func (c *UploadController) UploadFile(stream proto.FileService_UploadFileServer) error {
	ctx, cancel := utils.NewWithCancel(stream.Context())
	defer cancel()

	var fileInfo *proto.FileInfo
	var data []byte

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			req, err := stream.Recv()
			if err == io.EOF {
				file := &domain.File{
					Name: fileInfo.Name,
					Data: data,
				}

				if err := c.useCase.Upload(ctx, file); err != nil {
					return err
				}

				return stream.SendAndClose(&proto.UploadResponse{
					Name: file.Name,
					Size: uint32(len(file.Data)),
				})
			}
			if err != nil {
				return err
			}

			switch x := req.Data.(type) {
			case *proto.UploadRequest_Info:
				fileInfo = x.Info
			case *proto.UploadRequest_Chunk:
				data = append(data, x.Chunk...)
			}
		}
	}
}
