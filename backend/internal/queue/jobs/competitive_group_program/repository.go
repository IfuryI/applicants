package competitive_group_program

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
	query := `insert into mdb.competitive_group_program (
			uid_competitive_group,
			uid_subdivision_org,
			uid,
			uid_education_program) values ($1, $2, $3, $4)`

	result, err := r.db.Exec(ctx, query,
		pd.CompetitiveGroupProgram.UID,
		pd.CompetitiveGroupProgram.UIDCompetitiveGroup,
		pd.CompetitiveGroupProgram.UIDSubdivisionOrg,
		pd.CompetitiveGroupProgram.UIDEducationProgram,
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
				uid_campaign,
				name,
				id_level_budget,
				id_education_level,
				id_education_source,
				id_education_form,
				admission_number,
				comment,
				idocso
			from mdb.competitive_group_program
			where uid = $1`

	row := r.db.QueryRow(ctx, query, pd.CompetitiveGroupProgram.UID)

	var pdNew PackageData

	err := row.Scan(
		&pdNew.CompetitiveGroupProgram.UID,
		&pdNew.CompetitiveGroupProgram.UIDCompetitiveGroup,
		&pdNew.CompetitiveGroupProgram.UIDSubdivisionOrg,
		&pdNew.CompetitiveGroupProgram.UIDEducationProgram,
	)
	if err != nil {
		return PackageData{}, fmt.Errorf("scan error: %w", err)
	}

	return pdNew, nil
}
