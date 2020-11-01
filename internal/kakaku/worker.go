package kakaku

var liquidC *liquid

func init() {
	liquidC = newClient()
}

func UpdateAssetPrice() error {
	product, err := liquidC.getTicker(BTC_JPY_PAIR)
	if err != nil {
		return err
	}

	asset, err := FirstOrCreate(product.Symbol, Conn())
	if err != nil {
		return err
	}

	err = Conn().Model(asset).Update("price_jpy", product.LastTradedPrice).Error
	if err != nil {
		return err
	}

	return nil
}
