package entrance_test

import (
	"bitbucket.org/projectiu7/backend/src/master/internal/utils"
	"context"
	"fmt"
)

type Repository struct {
	db utils.PgxPoolIface
}

func NewRepository(db utils.PgxPoolIface) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(ctx context.Context, pd PackageData) (string, error) {
	query := `insert into mdb.entrance_test (
			uid,
			uid_competitive_group,
			id_entrance_test_type,
			test_name,
			is_ege,
			min_score,
			priority,
			id_subject,
			uid_replace_entrance_test) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	result, err := r.db.Exec(ctx, query,
		pd.EntranceTest.UID,
		pd.EntranceTest.UIDCompetitiveGroup,
		pd.EntranceTest.IDEntranceTestType,
		pd.EntranceTest.TestName,
		pd.EntranceTest.IsEge,
		pd.EntranceTest.MinScore,
		pd.EntranceTest.Priority,
		pd.EntranceTest.IDSubject,
		pd.EntranceTest.UIDReplaceEntranceTest,
	)
	if err != nil {
		return "", err
	}

	values := result.RowsAffected()

	if values != 1 {
		return "", fmt.Errorf("%w", err)
	}

	return "", nil
}

func (r *Repository) Get(ctx context.Context, pd PackageData) (PackageData, error) {
	query := `
			select 
				uid,
				uid_competitive_group,
				id_entrance_test_type,
				test_name,
				is_ege,
				min_score,
				priority,
				id_subject,
				uid_replace_entrance_test
			from mdb.entrance_test
			where uid = $1`

	row := r.db.QueryRow(ctx, query, pd.EntranceTest.UID)

	var pdNew PackageData

	err := row.Scan(
		&pdNew.EntranceTest.UID,
		&pdNew.EntranceTest.UIDCompetitiveGroup,
		&pdNew.EntranceTest.IDEntranceTestType,
		&pdNew.EntranceTest.TestName,
		&pdNew.EntranceTest.IsEge,
		&pdNew.EntranceTest.MinScore,
		&pdNew.EntranceTest.Priority,
		&pdNew.EntranceTest.IDSubject,
		&pdNew.EntranceTest.UIDReplaceEntranceTest,
	)
	if err != nil {
		return PackageData{}, fmt.Errorf("scan error: %w", err)
	}

	return pdNew, nil
}
