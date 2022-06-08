package admission_volume

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
	query := `insert into mdb.admission_volume (
		uid,
		uid_campaign,
		id_direction,
		id_education_level,
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
		pd.AdmissionVolume.UID,
		pd.AdmissionVolume.UIDCampaign,
		pd.AdmissionVolume.IDDirection,
		pd.AdmissionVolume.IDEducationLevel,
		pd.AdmissionVolume.BudgetO,
		pd.AdmissionVolume.BudgetOZ,
		pd.AdmissionVolume.BudgetZ,
		pd.AdmissionVolume.QuotaO,
		pd.AdmissionVolume.QuotaOZ,
		pd.AdmissionVolume.QuotaZ,
		pd.AdmissionVolume.PaidO,
		pd.AdmissionVolume.PaidOZ,
		pd.AdmissionVolume.PaidZ,
		pd.AdmissionVolume.TargetO,
		pd.AdmissionVolume.TargetOZ,
		pd.AdmissionVolume.TargetZ,
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
					id_direction,
					id_education_level,
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
			from mdb.admission_volume
			where uid = $1`

	row := r.db.QueryRow(ctx, query, pd.AdmissionVolume.UID)

	var pdNew PackageData

	err := row.Scan(
		&pdNew.AdmissionVolume.UID,
		&pdNew.AdmissionVolume.UIDCampaign,
		&pdNew.AdmissionVolume.IDDirection,
		&pdNew.AdmissionVolume.IDEducationLevel,
		&pdNew.AdmissionVolume.BudgetO,
		&pdNew.AdmissionVolume.BudgetOZ,
		&pdNew.AdmissionVolume.BudgetZ,
		&pdNew.AdmissionVolume.QuotaO,
		&pdNew.AdmissionVolume.QuotaOZ,
		&pdNew.AdmissionVolume.QuotaZ,
		&pdNew.AdmissionVolume.PaidO,
		&pdNew.AdmissionVolume.PaidOZ,
		&pdNew.AdmissionVolume.PaidZ,
		&pdNew.AdmissionVolume.TargetO,
		&pdNew.AdmissionVolume.TargetOZ,
		&pdNew.AdmissionVolume.TargetZ,
	)
	if err != nil {
		return PackageData{}, fmt.Errorf("scan error: %w", err)
	}

	return pdNew, nil
}
