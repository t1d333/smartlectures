package repository

import (
	"context"

	"github.com/t1d333/smartlectures/internal/models"
)

type Repository interface {
	GetDir(dirId int, ctx context.Context) (models.Dir, error)
	CreateDir(note models.Dir, ctx context.Context) (int, error)
	DeleteDir(dirId int, ctx context.Context) error
	UpdateDir(dir models.Dir, ctx context.Context) error
	GetDirsOverview(userId int, ctx context.Context) (models.DirsOverview, error)
}
