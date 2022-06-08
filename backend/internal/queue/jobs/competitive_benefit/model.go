package competitive_benefit

type PackageData struct {
	CompetitiveBenefit CompetitiveBenefit `xml:"CompetitiveBenefit"`
}

type CompetitiveBenefit struct {
	UID                  string `xml:"UID"`
	UIDCompetitiveGroup  string `xml:"UIDCompetitiveGroup"`
	IDOlimpicType        int64  `xml:"IDOlimpicType"`
	IDOlimpicLevels      int64  `xml:"IDOlimpicLevels"`
	IDBenefit            int64  `xml:"IDBenefit"`
	IDOlimpicDiplomaType int64  `xml:"IDOlimpicDiplomaType"`
	EgeMinValue          int64  `xml:"EgeMinValue"`
	OlimpicProfile       string `xml:"OlimpicProfile"`
}
