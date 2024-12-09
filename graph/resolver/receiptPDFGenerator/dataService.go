package receiptPDFGenerator

import "github.com/bishal-dd/receipt-generator-backend/graph/model"

func (r *ReceiptPDFGeneratorResolver) GetProfileByUserID(userId string) (*model.Profile, error) {
	var profile *model.Profile
	if err := r.db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}