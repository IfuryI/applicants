package entrance_test_benefit

type PackageData struct {
	EntranceTestBenefit EntranceTestBenefit `xml:"EntranceTestBenefit"`
}

type EntranceTestBenefit struct {
	UID               string `xml:"UID"`
	UIDEntranceTest   string `xml:"UIDEntranceTest"`
	IDBenefit         int64  `xml:"IDBenefit"`
	IDDiplomaType     int64  `xml:"IDDiplomaType"`
	IDOlimpicClasses  int64  `xml:"IDOlimpicClasses"`
	IDOlimpicLevel    int64  `xml:"IDOlimpicLevel"`
	IDOlimpicProfiles int64  `xml:"IDOlimpicProfiles"`
	EgeMinValue       int64  `xml:"EgeMinValue"`
}
