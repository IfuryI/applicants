package cmp_achiviment

type PackageData struct {
	Achievement Achievement `xml:"Achievement"`
}

type Achievement struct {
	UIDCampaign string `xml:"UIDCampaign"`
	UID         string `xml:"UID"`
	IDCategory  int64  `xml:"IDCategory"`
	Name        string `xml:"Name"`
	MaxValue    int64  `xml:"MaxValue"`
}
