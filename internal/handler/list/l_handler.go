package list

import (
	"context"
	"file-service/api/proto"
	"file-service/internal/usecase/list"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListController struct {
	useCase *list.ListUseCase
	proto.UnimplementedFileServiceServer
}

func NewListController(useCase *list.ListUseCase) *ListController {
	return &ListController{useCase: useCase}
}

func (c *ListController) ListFiles(ctx context.Context, _ *emptypb.Empty) (*proto.ListResponse, error) {
	files, err := c.useCase.List(ctx)
	if err != nil {
		return nil, err
	}

	response := &proto.ListResponse{}
	for _, file := range files {
		response.Files = append(response.Files, &proto.FileInfo{
			Name:      file.Name,
			CreatedAt: timestamppb.New(file.CreatedAt),
			UpdatedAt: timestamppb.New(file.UpdatedAt),
		})
	}

	return response, nil
}
