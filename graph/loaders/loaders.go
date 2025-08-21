package loaders

// import vikstrous/dataloadgen with your other imports
import (
	"time"

	"github.com/bishal-dd/receipt-generator-backend/graph/loaders/encryptedServiceLoader"
	"github.com/bishal-dd/receipt-generator-backend/graph/loaders/profileLoader"
	"github.com/bishal-dd/receipt-generator-backend/graph/loaders/receiptLoader"
	"github.com/bishal-dd/receipt-generator-backend/graph/loaders/serviceLoader"
	"github.com/bishal-dd/receipt-generator-backend/graph/loaders/userLoader"
	"github.com/bishal-dd/receipt-generator-backend/graph/model"
	"github.com/vikstrous/dataloadgen"
	"gorm.io/gorm"
)

type Loaders struct {
	UserLoader             *dataloadgen.Loader[string, *model.User]
	ReceiptLoader          *dataloadgen.Loader[string, []*model.Receipt]
	ServiceLoader          *dataloadgen.Loader[string, []*model.Service]
	EncryptedServiceLoader *dataloadgen.Loader[string, []*model.EncryptedService]
	ProfileLoader          *dataloadgen.Loader[string, []*model.Profile]
}

func NewLoaders(conn *gorm.DB) *Loaders {
	return &Loaders{
		UserLoader:             dataloadgen.NewLoader(userLoader.NewUserReader(conn).GetUsers, dataloadgen.WithWait(time.Millisecond)),
		ReceiptLoader:          dataloadgen.NewLoader(receiptLoader.NewReceiptReader(conn).GetReceiptsByUserIds, dataloadgen.WithWait(time.Millisecond)),
		ServiceLoader:          dataloadgen.NewLoader(serviceLoader.NewServiceReader(conn).GetServicesByReceiptIds, dataloadgen.WithWait(time.Millisecond)),
		EncryptedServiceLoader: dataloadgen.NewLoader(encryptedServiceLoader.NewEncryptedServiceReader(conn).GetEncryptedServicesByReceiptIds, dataloadgen.WithWait(time.Millisecond)),
		ProfileLoader:          dataloadgen.NewLoader(profileLoader.NewProfileReader(conn).GetProfileByUserIds, dataloadgen.WithWait(time.Millisecond)),
	}
}
