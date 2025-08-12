package encryptedReceipt

import (
	"errors"
	"time"
)

func searchDataRangeValidation(dateRange []string) error {
	if dateRange != nil {
		if len(dateRange) != 2 {
			return errors.New("dateRange must have 2 elements")
		}
		if dateRange[0] == "" || dateRange[1] == "" {
			return errors.New("dateRange elements must not be empty")
		}
		if _, err := time.Parse("2006-01-02", dateRange[0]); err != nil {
			return errors.New("dateRange elements must be in the format YYYY-MM-DD")
		}
		if _, err := time.Parse("2006-01-02", dateRange[1]); err != nil {
			return errors.New("dateRange elements must be in the format YYYY-MM-DD")
		}

	}
	return nil
}
