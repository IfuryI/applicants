package competitive_group_program

type PackageData struct {
	CompetitiveGroupProgram CompetitiveGroupProgram `xml:"CompetitiveGroupProgram"`
}

type CompetitiveGroupProgram struct {
	UID                 string `xml:"UID"`
	UIDCompetitiveGroup string `xml:"UIDCompetitiveGroup"`
	UIDSubdivisionOrg   string `xml:"UIDSubdivisionOrg"`
	UIDEducationProgram string `xml:"UIDEducationProgram"`
}
