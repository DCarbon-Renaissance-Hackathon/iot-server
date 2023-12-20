package repo

import (
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type XSMImpl struct {
	db *gorm.DB
}

func NewXSMImpl(db *gorm.DB) (*XSMImpl, error) {
	err := db.AutoMigrate(&models.XSMetric{})
	if nil != err {
		return nil, err
	}

	var impl = &XSMImpl{
		db: db,
	}
	return impl, nil
}

func (impl *XSMImpl) Create(req *domain.RXSMCreate,
) (*models.XSMetric, error) {
	var data = &models.XSMetric{
		Id:         uuid.NewV4().String(),
		IotAddress: dmodels.EthAddress(req.Address.String()),
		SensorType: req.SensorType,
		Metric:     req.Metric,
		CreatedAt:  time.Now(),
	}
	var err = impl.tblXSM().Create(data).Error
	if nil != err {
		return nil, err
	}
	return data, nil
}

func (impl *XSMImpl) GetList(req *domain.RXSMGetList,
) ([]*models.XSMetric, error) {
	var data = make([]*models.XSMetric, 0)
	var query = impl.tblXSM().Offset(req.Skip).
		Select("sensor_type, metric, created_at")
	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}

	if req.From > 0 {
		query = query.Where("created_at > ?", time.Unix(req.From, 0))
	}

	if req.To > 0 {
		query = query.Where("created_at < ?", time.Unix(req.To, 0))
	}

	if req.Address != "" {
		query = query.Where("iot_address = ?", req.Address)
	}

	var err = query.Find(&data).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("", err)
	}
	return data, nil
}

func (impl *XSMImpl) tblXSM() *gorm.DB {
	return impl.db.Table(models.TableNameXM)
}
