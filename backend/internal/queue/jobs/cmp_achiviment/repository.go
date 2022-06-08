package cmp_achiviment

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
	query := `insert into mdb.cmp_achievement (uid_campaign, uid, id_category, name, max_value) values ($1, $2, $3, $4, $5)`

	result, err := r.db.Exec(ctx, query,
		pd.Achievement.UIDCampaign,
		pd.Achievement.UID,
		pd.Achievement.IDCategory,
		pd.Achievement.Name,
		pd.Achievement.MaxValue,
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
			select uid_campaign, uid, id_category, name, max_value
			from mdb.cmp_achievement
			where uid = $1`

	row := r.db.QueryRow(ctx, query, pd.Achievement.UID)

	var pdNew PackageData

	err := row.Scan(
		&pdNew.Achievement.UIDCampaign,
		&pdNew.Achievement.UID,
		&pdNew.Achievement.IDCategory,
		&pdNew.Achievement.Name,
		&pdNew.Achievement.MaxValue,
	)
	if err != nil {
		return PackageData{}, fmt.Errorf("scan error: %w", err)
	}

	return pdNew, nil
}
