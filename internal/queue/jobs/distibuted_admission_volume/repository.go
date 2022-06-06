package distibuted_admission_volume

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
	query := `insert into mdb.distibuted_admission_volume (
		uid,
		uid_admission_volume,
		id_direction,
		id_level_budget,
		budget_o,
		budget_oz,
		budget_z,
		quota_o,
		quota_oz,
		quota_z,
		paid_o,
		paid_oz,
		paid_z,
		target_o,
		target_oz,
		target_z) values (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	result, err := r.db.Exec(ctx, query,
		pd.DistributedAdmissionVolume.UID,
		pd.DistributedAdmissionVolume.UIDAdmissionVolume,
		pd.DistributedAdmissionVolume.IDDirection,
		pd.DistributedAdmissionVolume.IDLevelBudget,
		pd.DistributedAdmissionVolume.BudgetO,
		pd.DistributedAdmissionVolume.BudgetOZ,
		pd.DistributedAdmissionVolume.BudgetZ,
		pd.DistributedAdmissionVolume.QuotaO,
		pd.DistributedAdmissionVolume.QuotaOZ,
		pd.DistributedAdmissionVolume.QuotaZ,
		pd.DistributedAdmissionVolume.PaidO,
		pd.DistributedAdmissionVolume.PaidOZ,
		pd.DistributedAdmissionVolume.PaidZ,
		pd.DistributedAdmissionVolume.TargetO,
		pd.DistributedAdmissionVolume.TargetOZ,
		pd.DistributedAdmissionVolume.TargetZ,
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
				uid_admission_volume,
				id_direction,
				id_level_budget,
				budget_o,
				budget_oz,
				budget_z,
				quota_o,
				quota_oz,
				quota_z,
				paid_o,
				paid_oz,
				paid_z,
				target_o,
				target_oz,
				target_z
			from mdb.distibuted_admission_volume
			where uid = $1`

	row := r.db.QueryRow(ctx, query, pd.DistributedAdmissionVolume.UID)

	var pdNew PackageData

	err := row.Scan(
		&pdNew.DistributedAdmissionVolume.UID,
		&pdNew.DistributedAdmissionVolume.UIDAdmissionVolume,
		&pdNew.DistributedAdmissionVolume.IDDirection,
		&pdNew.DistributedAdmissionVolume.IDLevelBudget,
		&pdNew.DistributedAdmissionVolume.BudgetO,
		&pdNew.DistributedAdmissionVolume.BudgetOZ,
		&pdNew.DistributedAdmissionVolume.BudgetZ,
		&pdNew.DistributedAdmissionVolume.QuotaO,
		&pdNew.DistributedAdmissionVolume.QuotaOZ,
		&pdNew.DistributedAdmissionVolume.QuotaZ,
		&pdNew.DistributedAdmissionVolume.PaidO,
		&pdNew.DistributedAdmissionVolume.PaidOZ,
		&pdNew.DistributedAdmissionVolume.PaidZ,
		&pdNew.DistributedAdmissionVolume.TargetO,
		&pdNew.DistributedAdmissionVolume.TargetOZ,
		&pdNew.DistributedAdmissionVolume.TargetZ,
	)
	if err != nil {
		return PackageData{}, fmt.Errorf("scan error: %w", err)
	}

	return pdNew, nil
}
