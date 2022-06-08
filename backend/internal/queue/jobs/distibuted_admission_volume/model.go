package distibuted_admission_volume

type PackageData struct {
	DistributedAdmissionVolume DistributedAdmissionVolume `xml:"DistributedAdmissionVolume"`
}

type DistributedAdmissionVolume struct {
	UID                string `xml:"UID"`
	UIDAdmissionVolume string `xml:"UIDAdmissionVolume"`
	IDDirection        int64  `xml:"IDCategory"`
	IDLevelBudget      int64  `xml:"IDEducationLevel"`
	BudgetO            int64  `xml:"BudgetO"`
	BudgetOZ           int64  `xml:"BudgetOZ"`
	BudgetZ            int64  `xml:"BudgetZ"`
	QuotaO             int64  `xml:"QuotaO"`
	QuotaOZ            int64  `xml:"QuotaOZ"`
	QuotaZ             int64  `xml:"QuotaZ"`
	PaidO              int64  `xml:"PaidO"`
	PaidOZ             int64  `xml:"PaidOZ"`
	PaidZ              int64  `xml:"PaidZ"`
	TargetO            int64  `xml:"TargetO"`
	TargetOZ           int64  `xml:"TargetOZ"`
	TargetZ            int64  `xml:"TargetZ"`
}
