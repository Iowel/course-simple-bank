package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	RUB = "RUB"
	CAD = "CAD"
)

// IsSupportedCurrency return true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, RUB:
		return true
	}
	return false
}
