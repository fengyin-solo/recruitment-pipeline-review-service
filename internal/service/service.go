package service

import (
	"recruit/internal/config"
	"recruit/internal/store"
	"recruit/pkg/logger"
)

type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}

// pageBounds 计算分页区间 [start, end)，对 page/size 做防御性归一化。
// page <= 0 视作第 1 页，避免 (page-1)*size 产生负数下标导致 slice 越界 panic
// （线上网关偶尔会把空页码转成 0）；size <= 0 视作空页。
// 对乘法/加法溢出同样兜底，返回值始终满足 0 <= start <= end <= total。
func pageBounds(page, size, total int) (start, end int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		return 0, 0
	}
	start = (page - 1) * size
	if start < 0 || start >= total {
		return 0, 0
	}
	end = start + size
	if end < 0 || end > total {
		end = total
	}
	return start, end
}
