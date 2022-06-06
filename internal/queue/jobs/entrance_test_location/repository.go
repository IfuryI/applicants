package entrance_test_location

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
	query := `insert into mdb.entrance_test_location (
		id_choice,
		uid_choice_epgu,
		uid_entrance_test,
		test_date,
		test_location,
		entrance_count) values ($1, $2, $3, $4, $5, $6)`

	result, err := r.db.Exec(ctx, query,
		pd.EntranceTestLocation.IDChoice,
		pd.EntranceTestLocation.UIDChoice,
		pd.EntranceTestLocation.UIDEntranceTest,
		pd.EntranceTestLocation.TestDate,
		pd.EntranceTestLocation.TestLocation,
		pd.EntranceTestLocation.EntranceCount,
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
				id_choice,
				uid_choice_epgu,
				uid_entrance_test,
				test_date,
				test_location,
				entrance_count
			from mdb.entrance_test_location
			where id_choice = $1`

	row := r.db.QueryRow(ctx, query, pd.EntranceTestLocation.IDChoice)

	var pdNew PackageData

	err := row.Scan(
		&pdNew.EntranceTestLocation.IDChoice,
		&pdNew.EntranceTestLocation.UIDChoice,
		&pdNew.EntranceTestLocation.UIDEntranceTest,
		&pdNew.EntranceTestLocation.TestDate,
		&pdNew.EntranceTestLocation.TestLocation,
		&pdNew.EntranceTestLocation.EntranceCount,
	)
	if err != nil {
		return PackageData{}, fmt.Errorf("scan error: %w", err)
	}

	return pdNew, nil
}
