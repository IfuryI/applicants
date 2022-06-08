package competitive_benefit

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
	query := `insert into mdb.competitive_benefit (
		uid,
		uid_competitive_group,
		id_olimpic_type,
		id_olimpic_levels,
		id_benefit,
		id_olimpic_diploma_type,
		ege_min_value,
		olimpic_profile) values ($1, $2, $3, $4, $5, $6, $7, $8)`

	result, err := r.db.Exec(ctx, query,
		pd.CompetitiveBenefit.UID,
		pd.CompetitiveBenefit.UIDCompetitiveGroup,
		pd.CompetitiveBenefit.IDOlimpicType,
		pd.CompetitiveBenefit.IDOlimpicLevels,
		pd.CompetitiveBenefit.IDBenefit,
		pd.CompetitiveBenefit.IDOlimpicDiplomaType,
		pd.CompetitiveBenefit.EgeMinValue,
		pd.CompetitiveBenefit.OlimpicProfile,
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
				id_olimpic_type,
				id_olimpic_levels,
				id_benefit,
				id_olimpic_diploma_type,
				ege_min_value,
				olimpic_profile
			from mdb.competitive_benefit
			where uid = $1`

	row := r.db.QueryRow(ctx, query, pd.CompetitiveBenefit.UID)

	var pdNew PackageData

	err := row.Scan(
		&pdNew.CompetitiveBenefit.UID,
		&pdNew.CompetitiveBenefit.UIDCompetitiveGroup,
		&pdNew.CompetitiveBenefit.IDOlimpicType,
		&pdNew.CompetitiveBenefit.IDOlimpicLevels,
		&pdNew.CompetitiveBenefit.IDBenefit,
		&pdNew.CompetitiveBenefit.IDOlimpicDiplomaType,
		&pdNew.CompetitiveBenefit.EgeMinValue,
		&pdNew.CompetitiveBenefit.OlimpicProfile,
	)
	if err != nil {
		return PackageData{}, fmt.Errorf("scan error: %w", err)
	}

	return pdNew, nil
}
