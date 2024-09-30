package sadykhanModels

type Catalog struct {
	Company  string  `xml:"company" json:"company"`
	Merchant string  `xml:"merchantid" json:"merchant_id"`
	Offers   []Offer `xml:"offers>offer" json:"offers"`
}

type Offer struct {
	SKU            string         `xml:"sku,attr" json:"sku"`
	Model          string         `xml:"model" json:"model"`
	Brand          string         `xml:"brand" json:"brand"`
	ExpirationDate string         `xml:"-" json:"expiration_date,omitempty"`
	Barcodes       Barcodes       `xml:"barcodes" json:"barcodes"`
	Availabilities Availabilities `xml:"availabilities" json:"availabilities"`
	CityPrices     CityPrices     `xml:"city_prices" json:"city_prices"`
}

type Barcodes struct {
	Codes []string `xml:"barcode" json:"barcodes"`
}

type Availabilities struct {
	Available    string `xml:"availability>available,attr" json:"available"`
	StoreID      string `xml:"availability>storeId,attr" json:"store_id"`
	Availability int    `xml:"availability" json:"availability"`
}

type CityPrices struct {
	CityID    int `xml:"city_price>city_id" json:"city_id"`
	CityPrice int `xml:"city_price" json:"city_price"`
}
