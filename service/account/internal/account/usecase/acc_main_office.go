package usecase

const (
	AccMainOfficeUnknown uint8 = 0
	AccMainOfficeMoscow  uint8 = 25
	AccMainOfficeKursk   uint8 = 38
)

var PossibleAccMainOffice = map[uint8]uint8{
	AccCountryRegionMoscow: AccMainOfficeMoscow,
	//AccCountryRegionKursk:  AccMainOfficeKursk,
}
