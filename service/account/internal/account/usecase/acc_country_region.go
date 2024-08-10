package usecase

const (
	AccCountryRegionUnknown uint8 = 0
	AccCountryRegionKursk   uint8 = 46
	AccCountryRegionMoscow  uint8 = 50
)

var PossibleAccCountryRegion = map[uint8][...]uint8{
	AccCountryRus: {
		AccCountryRegionKursk,
		AccCountryRegionMoscow,
	},
}
