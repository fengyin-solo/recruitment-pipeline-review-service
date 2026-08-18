// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"recruit/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	CreateJob(j *model.Job) error
	GetJob(id string) (*model.Job, error)
	ListJobs() []*model.Job
	UpdateJob(j *model.Job) error
	DeleteJob(id string) error

	CreateCandidate(c *model.Candidate) error
	GetCandidate(id string) (*model.Candidate, error)
	GetCandidateByEmail(email string) (*model.Candidate, error)
	ListCandidates() []*model.Candidate
	UpdateCandidate(c *model.Candidate) error
	DeleteCandidate(id string) error

	CreateResume(r *model.Resume) error
	GetResume(id string) (*model.Resume, error)
	ListResumes() []*model.Resume
	UpdateResume(r *model.Resume) error
	DeleteResume(id string) error

	CreateInterview(i *model.Interview) error
	GetInterview(id string) (*model.Interview, error)
	ListInterviews() []*model.Interview
	UpdateInterview(i *model.Interview) error
	DeleteInterview(id string) error

	CreateOffer(o *model.Offer) error
	GetOffer(id string) (*model.Offer, error)
	ListOffers() []*model.Offer
	UpdateOffer(o *model.Offer) error
	DeleteOffer(id string) error
}
