package campaign

type PackageData struct {
	Campaign Campaign `xml:"Campaign"`
}

// EducationForms ...
type EducationForms struct {
	IDEducationForm []int `xml:"IDEducationForm"`
}

// EducationLevels ...
type EducationLevels struct {
	IDEducationLevel []int `xml:"IDEducationLevel"`
}

// CompetitiveGroupList ...
type CompetitiveGroupList struct {
	UID string `xml:"UID"`
}

// AchievementList ...
type AchievementList struct {
	UID string `xml:"UID"`
}

// AdmissionVolumeList ...
type AdmissionVolumeList struct {
	UID string `xml:"UID"`
}

// EndApplicationList ...
type EndApplicationList struct {
	UID string `xml:"UID"`
}

// TermsAdmissionList ...
type TermsAdmissionList struct {
	UID string `xml:"UID"`
}

// Campaign ...
type Campaign struct {
	UID                  string                `xml:"UID"`
	Name                 string                `xml:"Name"`
	YearStart            int                   `xml:"YearStart"`
	YearEnd              int                   `xml:"YearEnd"`
	EducationForms       *EducationForms       `xml:"EducationForms"`
	EducationLevels      *EducationLevels      `xml:"EducationLevels"`
	IDCampaignType       int                   `xml:"IDCampaignType"`
	IDCampaignStatus     int                   `xml:"IDCampaignStatus"`
	NumberAgree          int                   `xml:"NumberAgree"`
	EndDate              string                `xml:"EndDate"`
	CountDirections      int                   `xml:"CountDirections"`
	CompetitiveGroupList *CompetitiveGroupList `xml:"CompetitiveGroupList"`
	AchievementList      *AchievementList      `xml:"AchievementList"`
	AdmissionVolumeList  *AdmissionVolumeList  `xml:"AdmissionVolumeList"`
	EndApplicationList   *EndApplicationList   `xml:"EndApplicationList"`
	TermsAdmissionList   *TermsAdmissionList   `xml:"TermsAdmissionList"`
}
