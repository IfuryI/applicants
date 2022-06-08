package competitive_group

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
	query := `insert into mdb.competitive_group (
		uid,
		uid_campaign,
		name,
		id_level_budget,
		id_education_level,
		id_education_source,
		id_education_form,
		admission_number,
		comment,
		idocso) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	result, err := r.db.Exec(ctx, query,
		pd.CompetitiveGroup.UID,
		pd.CompetitiveGroup.UIDCampaign,
		pd.CompetitiveGroup.Name,
		pd.CompetitiveGroup.IDLevelBudget,
		pd.CompetitiveGroup.IDEducationLevel,
		pd.CompetitiveGroup.IDEducationSource,
		pd.CompetitiveGroup.IDEducationForm,
		pd.CompetitiveGroup.AdmissionNumber,
		pd.CompetitiveGroup.Comment,
		pd.CompetitiveGroup.IDOCSO,
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
			from mdb.competitive_group
			where uid = $1`

	row := r.db.QueryRow(ctx, query, pd.CompetitiveGroup.UID)

	var pdNew PackageData

	err := row.Scan(
		&pdNew.CompetitiveGroup.UID,
		&pdNew.CompetitiveGroup.UIDCampaign,
		&pdNew.CompetitiveGroup.Name,
		&pdNew.CompetitiveGroup.IDLevelBudget,
		&pdNew.CompetitiveGroup.IDEducationLevel,
		&pdNew.CompetitiveGroup.IDEducationSource,
		&pdNew.CompetitiveGroup.IDEducationForm,
		&pdNew.CompetitiveGroup.AdmissionNumber,
		&pdNew.CompetitiveGroup.Comment,
		&pdNew.CompetitiveGroup.IDOCSO,
	)
	if err != nil {
		return PackageData{}, fmt.Errorf("scan error: %w", err)
	}

	return pdNew, nil
}
