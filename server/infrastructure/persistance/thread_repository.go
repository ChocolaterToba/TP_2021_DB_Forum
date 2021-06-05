package persistance

import (
	"dbforum/domain/entity"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ThreadRepo struct {
	postgresDB *pgxpool.Pool
}

func NewThreadRepo(postgresDB *pgxpool.Pool) *ThreadRepo {
	return &ThreadRepo{postgresDB}
}

func (threadRepo *ThreadRepo) CreateThread(thread *entity.Thread) (int, error) {
	return -1, nil
}

func (threadRepo *ThreadRepo) GetThreadByID(threadID int) (*entity.Thread, error) {
	return nil, nil
}

func (threadRepo *ThreadRepo) GetThreadByIDString(threadIDString string) (*entity.Thread, error) {
	return nil, nil
}

func (threadRepo *ThreadRepo) EditThread(thread *entity.Thread) error {
	return nil
}

func (threadRepo *ThreadRepo) GetPostsByThreadID(threadID int) ([]*entity.Post, error) {
	return nil, nil
}

func (threadRepo *ThreadRepo) GetPostsByThreadIDString(threadIDstring string) ([]*entity.Post, error) {
	return nil, nil
}
