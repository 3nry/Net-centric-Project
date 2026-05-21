package grpcserver

import (
	"context"
	"fmt"

	"mangahub/config"
	pb "mangahub/gen/manga/proto"
	"mangahub/models"
)

// MangaServer implement interface được generate từ proto
type MangaServer struct {
	pb.UnimplementedMangaServiceServer
}

// GetManga — lấy 1 manga theo ID
func (s *MangaServer) GetManga(ctx context.Context, req *pb.MangaRequest) (*pb.MangaResponse, error) {
	var m models.Manga
	if err := config.DB.First(&m, req.Id).Error; err != nil {
		return nil, fmt.Errorf("manga not found: %w", err)
	}
	return &pb.MangaResponse{
		Id:          int64(m.ID),
		Title:       m.Title,
		CoverImage:  m.CoverImage,
		Author:      m.Author,
		Description: m.Description,
		Genre:       m.Genre,
		Status:      m.Status,
	}, nil
}

// SearchManga — tìm manga theo tên
func (s *MangaServer) SearchManga(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	var mangas []models.Manga
	config.DB.Where("title LIKE ?", "%"+req.Query+"%").Limit(20).Find(&mangas)

	var results []*pb.MangaResponse
	for _, m := range mangas {
		results = append(results, &pb.MangaResponse{
			Id:          int64(m.ID),
			Title:       m.Title,
			CoverImage:  m.CoverImage,
			Author:      m.Author,
			Description: m.Description,
			Genre:       m.Genre,
			Status:      m.Status,
		})
	}
	return &pb.SearchResponse{Results: results}, nil
}

// UpdateProgress — cập nhật tiến độ đọc của user
func (s *MangaServer) UpdateProgress(ctx context.Context, req *pb.ProgressRequest) (*pb.ProgressResponse, error) {
	var progress models.UserProgress
	result := config.DB.Where("user_id = ? AND manga_id = ?", req.UserId, req.MangaId).First(&progress)

	if result.Error != nil {
		// Chưa có record → tạo mới
		progress = models.UserProgress{
			UserID:         uint(req.UserId),
			MangaID:        uint(req.MangaId),
			CurrentChapter: int(req.CurrentChapter),
			Status:         req.Status,
		}
		config.DB.Create(&progress)
	} else {
		// Đã có → update
		config.DB.Model(&progress).Updates(map[string]interface{}{
			"current_chapter": req.CurrentChapter,
			"status":          req.Status,
		})
	}

	return &pb.ProgressResponse{
		Success: true,
		Message: "Progress updated successfully",
	}, nil
}
