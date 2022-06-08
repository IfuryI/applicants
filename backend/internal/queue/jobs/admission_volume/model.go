package admission_volume

type PackageData struct {
	AdmissionVolume AdmissionVolume `xml:"AdmissionVolume"`
}

type AdmissionVolume struct {
	UID              string `xml:"UID"`
	UIDCampaign      string `xml:"UIDCampaign"`
	IDDirection      int64  `xml:"IDCategory"`
	IDEducationLevel int64  `xml:"IDEducationLevel"`
	BudgetO          int64  `xml:"BudgetO"`
	BudgetOZ         int64  `xml:"BudgetOZ"`
	BudgetZ          int64  `xml:"BudgetZ"`
	QuotaO           int64  `xml:"QuotaO"`
	QuotaOZ          int64  `xml:"QuotaOZ"`
	QuotaZ           int64  `xml:"QuotaZ"`
	PaidO            int64  `xml:"PaidO"`
	PaidOZ           int64  `xml:"PaidOZ"`
	PaidZ            int64  `xml:"PaidZ"`
	TargetO          int64  `xml:"TargetO"`
	TargetOZ         int64  `xml:"TargetOZ"`
	TargetZ          int64  `xml:"TargetZ"`
}
