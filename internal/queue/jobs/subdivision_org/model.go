package subdivision_org

type PackageData struct {
	SubdivisionOrg SubdivisionOrg `xml:"SubdivisionOrg"`
}

type SubdivisionOrg struct {
	UID  string `xml:"UID"`
	Name string `xml:"Name"`
}
