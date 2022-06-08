package entrance_test

type PackageData struct {
	EntranceTest EntranceTest `xml:"EntranceTest"`
}

type EntranceTest struct {
	UID                    string `xml:"UID"`
	UIDCompetitiveGroup    string `xml:"UIDCompetitiveGroup"`
	IDEntranceTestType     int64  `xml:"IDEntranceTestType"`
	TestName               string `xml:"TestName"`
	IsEge                  bool   `xml:"IsEge"`
	MinScore               int64  `xml:"MinScore"`
	Priority               int64  `xml:"Priority"`
	IDSubject              int64  `xml:"IDSubject"`
	UIDReplaceEntranceTest string `xml:"UIDReplaceEntranceTest"`
}
