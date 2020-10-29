package kakaku

type Asset struct {
	gorm.Model
	Symbol string
	PriceJPY decimal.Decimal
}