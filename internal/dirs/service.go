package dirs

import (
	"context"

	"github.com/t1d333/smartlectures/internal/models"
)

type Service interface {
	GetDir(ctx context.Context, dirId int) (models.Dir, error)
	CreateDir(ctx context.Context, dir models.Dir) (int, error)
	DeleteDir(ctx context.Context, dirId int) error
	UpdateDir(ctx context.Context, dir models.Dir) error
	GetDirsOverview(ctx context.Context, userId int) (models.DirsOverview, error)
}
