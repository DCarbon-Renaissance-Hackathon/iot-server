package repo

import (
	"time"

	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SensorRepo struct {
	db *gorm.DB
}

func NewSensorRepo() (*SensorRepo, error) {
	var db = getSingletonDB()
	err := db.AutoMigrate(
		&models.Sensor{},
		&models.SM{},
		&models.SMSignature{},
	)
	if nil != err {
		return nil, err
	}

	var impl = &SensorRepo{
		db: db,
	}

	return impl, nil
}

func (impl *SensorRepo) CreateSensor(req *domain.RCreateSensor,
) (*models.Sensor, error) {
	var sensor = &models.Sensor{
		ID:        0,
		IotID:     req.IotID,
		Address:   &req.Address,
		Type:      req.Type,
		Status:    models.SensorStatusRegister,
		CreatedAt: time.Now(),
	}
	var err = impl.tblSensor().Create(sensor).Error
	if nil != err {
		return nil, models.ParsePostgresError("Create sensor", err)
	}

	return sensor, nil
}

func (impl *SensorRepo) ChangeSensorStatus(req *domain.RChangeSensorStatus,
) (*models.Sensor, error) {
	var sensor = &models.Sensor{}
	var err = impl.tblSensor().
		Model(sensor).
		Clauses(clause.Returning{}).
		Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"status": req.Status,
		}).
		Error
	if nil != err {
		return nil, models.ParsePostgresError("Change sensor status", err)
	}

	return sensor, nil
}

func (impl *SensorRepo) GetSensor(req *domain.RGetSensor,
) (*models.Sensor, error) {
	var sensor = &models.Sensor{}
	var query = impl.tblSensor()
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	} else {
		query = query.Where("address = ?", req.Address.String())
	}
	var err = query.First(sensor).Error
	if nil != err {
		return nil, models.ParsePostgresError("Get sensor", err)
	}
	return sensor, nil
}

func (impl *SensorRepo) GetSensors(req *domain.RGetSensors,
) ([]*models.Sensor, error) {
	var sensors = make([]*models.Sensor, 0, req.Limit)
	var query = impl.tblSensor().Offset(req.Skip).Limit(req.Limit)
	if req.IotId != 0 {
		query = query.Where("iod_id = ?", req.IotId)
	}
	var err = query.Find(&sensors).Error
	if nil != err {
		return nil, models.ParsePostgresError("Get sensors", err)
	}
	return sensors, nil
}

func (impl *SensorRepo) CreateSM(req *domain.RCreateSM,
) (*models.SM, error) {
	sensor, err := impl.GetSensor(&domain.RGetSensor{Address: req.SensorAddress})
	if nil != err {
		return nil, err
	}

	if sensor.Address.IsEmpty() {
		return nil, models.NewError(models.ECodeSensorHasNoAddress, "SensorAddress is empty")
	}

	var signed = &models.SMSignature{
		ID:        uuid.NewV4().String(),
		IsIotSign: true,
		IotID:     sensor.IotID,
		SensorID:  sensor.ID,
		Data:      req.Data,
		Signed:    req.Signed,
	}

	smx, err := signed.VerifySignature(*sensor.Address)
	if nil != err {
		return nil, err
	}

	sm, err := impl.insertMetric(smx, signed)
	if nil != err {
		return nil, nil
	}

	return sm, nil
}

func (impl *SensorRepo) CreateSMFromIot(req *domain.RCreateSMFromIOT,
) (*models.SM, error) {
	sensor, err := impl.GetSensor(&domain.RGetSensor{ID: req.SensorID})
	if nil != err {
		return nil, err
	}

	if !sensor.Address.IsEmpty() {
		return nil, models.NewError(models.ECodeSensorHasAddress, "SensorAddress is not empty")
	}

	var signed = &models.SMSignature{
		ID:        uuid.NewV4().String(),
		IsIotSign: true,
		IotID:     sensor.IotID,
		SensorID:  sensor.ID,
		Data:      req.Data,
		Signed:    req.Signed,
	}

	smx, err := signed.VerifySignature(req.IotAddress)
	if nil != err {
		return nil, err
	}

	sm, err := impl.insertMetric(smx, signed)
	if nil != err {
		return nil, nil
	}

	return sm, nil
}

func (impl *SensorRepo) insertMetric(smx *models.SMExtract, signed *models.SMSignature,
) (*models.SM, error) {
	sm := &models.SM{
		ID:        uuid.NewV4().String(),
		SignID:    signed.ID,
		Indicator: smx.Indicator,
		CreatedAt: time.Unix(smx.To, 0),
	}

	e1 := impl.tblMetrics().Transaction(func(dbTx *gorm.DB) error {
		err := dbTx.Table(models.TableNameSM).Create(sm).Error
		if nil != err {
			return models.ParsePostgresError("", err)
		}

		err = dbTx.Table(models.TableNameSMSignature).Create(signed).Error
		if nil != err {
			return models.ParsePostgresError("", err)
		}
		return nil
	})

	return sm, e1
}

func (impl *SensorRepo) GetMetrics(req *domain.RGetSM) ([]*models.SM, error) {
	var rs = make([]*models.SM, 0)
	var query = impl.tblMetrics()
	var err = query.Find(&rs).Error
	if nil != err {
		return nil, models.ParsePostgresError("Get metrics", err)
	}
	return rs, nil
}

func (impl *SensorRepo) tblSensor() *gorm.DB {
	return impl.db.Table(models.TableNameSensors)
}

func (impl *SensorRepo) tblMetrics() *gorm.DB {
	return impl.db.Table(models.TableNameSM)
}
