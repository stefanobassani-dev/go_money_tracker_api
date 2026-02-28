package tink

type amount struct {
	Value        value  `json:"value"`
	CurrencyCode string `json:"currencyCode"`
}

type financialInstitution struct {
	AccountNumber    string                 `json:"accountNumber"`
	ReferenceNumbers map[string]interface{} `json:"referenceNumbers"`
}
