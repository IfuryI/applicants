package education_program

type PackageData struct {
	EducationProgram EducationProgram `xml:"EducationProgram"`
}

type EducationProgram struct {
	UID             string `xml:"UID"`
	Name            string `xml:"Name"`
	IDEducationForm int64  `xml:"IDEducationForm"`
	IDOCSO          int64  `xml:"IDOCSO"`
}
