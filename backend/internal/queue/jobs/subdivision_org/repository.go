package subdivision_org

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
	query := `insert into mdb.subdivision_org (uid, name) values ($1, $2)`

	result, err := r.db.Exec(ctx, query, pd.SubdivisionOrg.UID, pd.SubdivisionOrg.Name)
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
			select uid, name
			from mdb.subdivision_org
			where uid = $1`

	row := r.db.QueryRow(ctx, query, pd.SubdivisionOrg.UID)

	var pdNew PackageData

	err := row.Scan(&pdNew.SubdivisionOrg.UID, &pdNew.SubdivisionOrg.Name)
	if err != nil {
		return PackageData{}, fmt.Errorf("scan error: %w", err)
	}

	return pdNew, nil
}
