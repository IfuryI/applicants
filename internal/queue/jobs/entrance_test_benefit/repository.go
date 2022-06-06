package entrance_test_benefit

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
	query := `insert into mdb.entrance_test_benefit (
		uid,
		uid_entrance_test,
		id_benefit,
		id_olimpic_diploma_type,
		id_olympic_classes,
		id_olympic_level,
		id_olympic_profiles,
		ege_min_value) values ($1, $2, $3, $4, $5, $6, $7, $8)`

	result, err := r.db.Exec(ctx, query,
		pd.EntranceTestBenefit.UID,
		pd.EntranceTestBenefit.UIDEntranceTest,
		pd.EntranceTestBenefit.IDBenefit,
		pd.EntranceTestBenefit.IDDiplomaType,
		pd.EntranceTestBenefit.IDOlimpicClasses,
		pd.EntranceTestBenefit.IDOlimpicLevel,
		pd.EntranceTestBenefit.IDOlimpicProfiles,
		pd.EntranceTestBenefit.EgeMinValue,
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
				uid_entrance_test,
				id_benefit,
				id_olimpic_diploma_type,
				id_olympic_classes,
				id_olympic_level,
				id_olympic_profiles,
				ege_min_value
			from mdb.entrance_test_benefit
			where uid = $1`

	row := r.db.QueryRow(ctx, query, pd.EntranceTestBenefit.UID)

	var pdNew PackageData

	err := row.Scan(
		&pdNew.EntranceTestBenefit.UID,
		&pdNew.EntranceTestBenefit.UIDEntranceTest,
		&pdNew.EntranceTestBenefit.IDBenefit,
		&pdNew.EntranceTestBenefit.IDDiplomaType,
		&pdNew.EntranceTestBenefit.IDOlimpicClasses,
		&pdNew.EntranceTestBenefit.IDOlimpicLevel,
		&pdNew.EntranceTestBenefit.IDOlimpicProfiles,
		&pdNew.EntranceTestBenefit.EgeMinValue,
	)
	if err != nil {
		return PackageData{}, fmt.Errorf("scan error: %w", err)
	}

	return pdNew, nil
}
