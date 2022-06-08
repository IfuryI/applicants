package campaign

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
	query := `insert into mdb.education_program (uid, name, id_education_form, idocso) values ($1, $2, $3, $4)`

	result, err := r.db.Exec(ctx, query,
		pd.Campaign.UID,
		pd.Campaign.Name,
		pd.Campaign.IDEducationForm,
		pd.Campaign.IDOCSO,
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
			select uid, name, id_education_form, idocso
			from mdb.education_program
			where uid = $1`

	row := r.db.QueryRow(ctx, query, pd.EducationProgram.UID)

	var pdNew PackageData

	err := row.Scan(&pdNew.EducationProgram.UID,
		&pdNew.EducationProgram.Name,
		&pdNew.EducationProgram.IDEducationForm,
		&pdNew.EducationProgram.IDOCSO,
	)
	if err != nil {
		return PackageData{}, fmt.Errorf("scan error: %w", err)
	}

	return pdNew, nil
}
