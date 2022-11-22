package repo

import (
	"log"
	"strings"
	"time"

	"github.com/Dcarbon/iott-cloud/domain"
	"github.com/Dcarbon/iott-cloud/libs/dbutils"
	"github.com/Dcarbon/iott-cloud/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type iotRepo struct {
	db *gorm.DB
}

func NewIOTRepo(dbUrl string) (domain.IIot, error) {
	var db, err = dbutils.NewDB(dbUrl)
	if nil != err {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.IOTDevice{},
		&models.Metric{},
		&models.MintSign{},
	)
	if nil != err {
		return nil, err
	}

	var ip = &iotRepo{
		db: db,
	}
	return ip, nil
}

func (ip *iotRepo) Create(iot *models.IOTDevice) error {
	iot.ID = 0
	iot.Address = strings.ToLower(iot.Address)
	var err = ip.tblIOT().Create(iot).Error
	if nil != err {
		return models.ParsePostgresError("IOT", err)
	}
	return nil
}

func (ip *iotRepo) ChangeStatus(iotAddr string, status models.IOTStatus,
) (*models.IOTDevice, error) {
	var iot = &models.IOTDevice{}
	var err = ip.tblIOT().
		Model(iot).
		Clauses(clause.Returning{}).
		Where("address = ?", iotAddr).
		Update("status", status).
		Error

	return iot, models.ParsePostgresError("IOT", err)
}

func (ip *iotRepo) GetByBB(min, max *models.Point4326,
) ([]*models.IOTDevice, error) {
	var iots = make([]*models.IOTDevice, 0)
	var err = ip.tblIOT().
		Where(
			"ST_WITHIN(pos, ST_MakeEnvelope(?, ?, ?, ?, 4326))",
			min.Lng, min.Lat, max.Lng, max.Lat).
		Find(&iots).Error
	return iots, models.ParsePostgresError("IOT", err)
}

func (ip *iotRepo) CreateMetric(m *models.Metric) error {
	log.Println("Create metric for ", m.Address)
	m.ID = uuid.NewV4().String()
	m.CreatedAt = time.Now()
	m.Address = strings.ToLower(m.Address)
	var status = ip.GetIOTStatus(m.Address)
	if status < 0 {
		return models.NewError(models.ECodeIOTNotAllowed, "Iot status is not valid")
	}

	var err = ip.tblMetrics().Create(m).Error
	return models.ParsePostgresError("Metrics", err)
}

func (ip *iotRepo) GetMetrics(iot string, from, to int64,
) ([]*models.Metric, error) {
	iot = strings.ToLower(iot)

	var ftime = time.Unix(from, 0)
	var fto = time.Unix(to, 0)
	var metrics = make([]*models.Metric, 0)
	var err = ip.tblMetrics().
		Select("id, is_result, warning, metrics, created_at").
		Where(
			"address = ? AND created_at >= ? AND created_at <= ?",
			iot, ftime, fto).
		Find(&metrics).
		Error
	return metrics, models.ParsePostgresError("Metrics", err)
}

func (ip *iotRepo) GetRawMetric(metricId string,
) (*models.Metric, error) {

	var metric = &models.Metric{}
	var err = ip.tblMetrics().
		Where("id = ?", metricId).
		First(metric).Error
	return metric, models.ParsePostgresError("Metrics", err)
}

func (ip *iotRepo) GetIOTStatus(iotAddr string) models.IOTStatus {
	var device = &models.IOTDevice{}
	var err = ip.tblIOT().
		Select("status").
		Where("address = ?", strings.ToLower(iotAddr)).
		First(&device).Error
	if nil != err {
		device.Status = models.IOTStatusReject
	}
	return device.Status
}

func (ip *iotRepo) CreateMint(mint *models.MintSign,
) error {
	mint.IOT = strings.ToLower(mint.IOT)

	err := mint.Verify()
	if nil != err {
		return err
	}

	var latest = make([]*models.MintSign, 0, 1)
	err = ip.tblSign().
		Where("iot = ?", mint.IOT).
		Order("created_at desc").
		Limit(1).
		Find(&latest).Error
	if nil != err {
		return models.ParsePostgresError("", err)
	}

	if len(latest) == 0 {
		if mint.Nonce != 1 {
			return models.NewError(
				models.ECodeIOTInvalidNonce,
				"Nonce is not valid",
			)
		}
		err = models.ParsePostgresError(
			"",
			ip.tblSign().Create(mint).Error,
		)
	} else if latest[0].Nonce == mint.Nonce {
		err = ip.tblSign().
			Where("id = ?", latest[0].ID).
			Updates(map[string]interface{}{
				"nonce":      mint.Nonce,
				"amount":     mint.Amount,
				"r":          mint.R,
				"s":          mint.S,
				"v":          mint.V,
				"updated_at": time.Now(),
			}).Error
		err = models.ParsePostgresError("", err)
	} else if latest[0].Nonce+1 == mint.Nonce {
		err = models.ParsePostgresError(
			"",
			ip.tblSign().Create(mint).Error,
		)
	} else {
		err = models.NewError(models.ECodeIOTInvalidNonce, "")
	}

	return err
}

func (ip *iotRepo) GetMintSigns(iotAddr string, fromNonce int,
) ([]*models.MintSign, error) {
	iotAddr = strings.ToLower(iotAddr)

	var signeds = make([]*models.MintSign, 0)
	var err = ip.tblSign().
		Where("iot = ? AND nonce >= ?", iotAddr, fromNonce).
		Find(&signeds).
		Order("id asc").
		Error
	if nil != err {
		return nil, models.ParsePostgresError("", err)
	}
	return signeds, nil
}

func (ip *iotRepo) tblIOT() *gorm.DB {
	return ip.db.Table(models.TableNameIOT)
}

func (ip *iotRepo) tblMetrics() *gorm.DB {
	return ip.db.Table(models.TableNameMetrics)
}

func (ip *iotRepo) tblSign() *gorm.DB {
	return ip.db.Table(models.TableNameMintSign)
}
