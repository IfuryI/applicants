package competitive_group

type PackageData struct {
	CompetitiveGroup CompetitiveGroup `xml:"CompetitiveGroup"`
}

type CompetitiveGroup struct {
	UID               string `xml:"UID"`
	UIDCampaign       string `xml:"UIDCampaign"`
	Name              string `xml:"Name"`
	IDLevelBudget     int64  `xml:"IDLevelBudget"`
	IDEducationLevel  int64  `xml:"IDEducationLevel"`
	IDEducationSource int64  `xml:"IDEducationSource"`
	IDEducationForm   int64  `xml:"IDEducationForm"`
	AdmissionNumber   int64  `xml:"AdmissionNumber"`
	Comment           string `xml:"Comment"`
	IDOCSO            int64  `xml:"IDOCSO"`
}
