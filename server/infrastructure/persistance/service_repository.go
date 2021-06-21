package persistance

import (
	"context"
	"dbforum/domain/entity"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ServiceRepo struct {
	postgresDB *pgxpool.Pool
}

func NewServiceRepo(postgresDB *pgxpool.Pool) *ServiceRepo {
	return &ServiceRepo{postgresDB}
}

const getStatsQuery string = "SELECT MAX(case when seqnum = 1 then estimate end) as forums_count,\n" + // ransposition column into row
	"max(case when seqnum = 2 then estimate end) as posts_count,\n" +
	"max(case when seqnum = 3 then estimate end) as threads_count,\n" +
	"max(case when seqnum = 4 then estimate end) as users_count\n" +
	"from (select stats_col.*, row_number() over (order by relname) as seqnum\n" +
	"  from (\n" +
	"    SELECT reltuples AS estimate, relname\n" + // Actual stats request
	"    FROM pg_class\n" +
	"    WHERE relname in('forums', 'posts', 'threads', 'users')" +
	"  ) as stats_col\n" +
	") stats_with_rownumber"

func (serviceRepo *ServiceRepo) GetStats() (*entity.StatsInfo, error) {
	tx, err := serviceRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return nil, entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	row := tx.QueryRow(context.Background(), getStatsQuery)
	statsInfo := entity.StatsInfo{}
	err = row.Scan(&statsInfo.ForumsCount, &statsInfo.PostsCount, &statsInfo.ThreadsCount, &statsInfo.UsersCount)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, entity.TransactionCommitError
	}
	return &statsInfo, nil
}

const truncateAllQuery string = "TRUNCATE TABLE users RESTART IDENTITY CASCADE;" // Clears all tables, they all depend on users

func (serviceRepo *ServiceRepo) TruncateAll() error {
	tx, err := serviceRepo.postgresDB.Begin(context.Background())
	if err != nil {
		return entity.TransactionBeginError
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), truncateAllQuery)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return entity.TransactionCommitError
	}

	return nil
}

const vacuumPostsQuery string = "VACUUM (VERBOSE,ANALYZE) posts"

func (serviceRepo *ServiceRepo) VacuumPosts() error {
	_, err := serviceRepo.postgresDB.Exec(context.Background(), vacuumPostsQuery)
	if err != nil {
		return err
	}

	return nil
}
