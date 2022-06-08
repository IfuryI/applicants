package entrance_test_location

type PackageData struct {
	EntranceTestLocation EntranceTestLocation `xml:"EntranceTestLocation"`
}

type EntranceTestLocation struct {
	IDChoice        int64  `xml:"IDLevelBudget"`
	UIDChoice       string `xml:"UIDCampaign"`
	UIDEntranceTest string `xml:"UIDCampaign"`
	TestDate        string `xml:"UIDCampaign"`
	TestLocation    string `xml:"UIDCampaign"`
	EntranceCount   int64  `xml:"IDEducationLevel"`
}
