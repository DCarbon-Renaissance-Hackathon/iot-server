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
		&models.SmFloat{},
		&models.SmGPS{},
		&models.SmSignature{},
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

func (impl *SensorRepo) GetSensor(req *domain.SensorID,
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
	var query = impl.tblSensor().Offset(req.Skip)
	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}

	if req.IotId != 0 {
		query = query.Where("iot_id = ?", req.IotId)
	}

	var err = query.Find(&sensors).Error
	if nil != err {
		return nil, models.ParsePostgresError("Get sensors", err)
	}
	return sensors, nil
}

func (impl *SensorRepo) CreateSM(req *domain.RCreateSM,
) (*models.SmSignature, error) {
	sensor, err := impl.GetSensor(&domain.SensorID{Address: req.SensorAddress})
	if nil != err {
		return nil, err
	}

	if sensor.Address.IsEmpty() {
		return nil, models.NewError(models.ECodeSensorHasNoAddress, "SensorAddress is empty")
	}

	var signed = &models.SmSignature{
		ID:        uuid.NewV4().String(),
		IsIotSign: true,
		IotID:     sensor.IotID,
		SensorID:  sensor.ID,
		Data:      req.Data,
		Signed:    req.Signed,
	}

	smx, err := signed.VerifySignature(*sensor.Address, sensor.Type)
	if nil != err {
		return nil, err
	}

	err = impl.insertMetric(smx, signed, sensor.Type)
	if nil != err {
		return nil, nil
	}

	return signed, nil
}

func (impl *SensorRepo) CreateSMFromIot(req *domain.RCreateSMFromIOT,
) (*models.SmSignature, error) {
	sensor, err := impl.GetSensor(&domain.SensorID{ID: req.SensorID})
	if nil != err {
		return nil, err
	}

	if !sensor.Address.IsEmpty() {
		return nil, models.NewError(models.ECodeSensorHasAddress, "SensorAddress is not empty")
	}

	var signed = &models.SmSignature{
		ID:        uuid.NewV4().String(),
		IsIotSign: true,
		IotID:     sensor.IotID,
		SensorID:  sensor.ID,
		Data:      req.Data,
		Signed:    req.Signed,
	}

	smx, err := signed.VerifySignature(req.IotAddress, sensor.Type)
	if nil != err {
		return nil, err
	}

	err = impl.insertMetric(smx, signed, sensor.Type)
	if nil != err {
		return nil, err
	}

	return signed, nil
}

func (impl *SensorRepo) insertMetric(smx *models.SMExtract, signed *models.SmSignature, stype models.SensorType,
) error {

	switch stype {
	case models.SensorTypePower:
		return impl.insertMetricFloat(smx, signed)
	case models.SensorTypeFlow:
		return impl.insertMetricFloat(smx, signed)
	case models.SensorTypeGPS:
		return impl.insertMetricGPS(smx, signed)

	}
	return models.NewError(models.ECodeSensorInvalidType, "Sensor type is not supported")
	// sm := &models.SM{
	// 	ID:     uuid.NewV4().String(),
	// 	SignID: signed.ID,
	// 	// Indicator: smx.Indicator,
	// 	CreatedAt: time.Unix(smx.To, 0),
	// }

	// e1 := impl.tblMetrics().Transaction(func(dbTx *gorm.DB) error {
	// 	err := dbTx.Table(models.TableNameSM).Create(sm).Error
	// 	if nil != err {
	// 		return models.ParsePostgresError("", err)
	// 	}

	// 	err = dbTx.Table(models.TableNameSMSignature).Create(signed).Error
	// 	if nil != err {
	// 		return models.ParsePostgresError("", err)
	// 	}
	// 	return nil
	// })

	// return sm, e1
}

func (impl *SensorRepo) insertMetricFloat(smx *models.SMExtract, signed *models.SmSignature,
) error {
	sm := &models.SmFloat{
		ID:        uuid.NewV4().String(),
		SignID:    signed.ID,
		Indicator: float64(smx.Indicator.Value),
		CreatedAt: time.Unix(smx.To, 0),
	}

	e1 := impl.tblFloatMetrics().Transaction(func(dbTx *gorm.DB) error {
		err := dbTx.Table(models.TableNameSmFloat).Create(sm).Error
		if nil != err {
			return models.ParsePostgresError("", err)
		}

		err = dbTx.Table(models.TableNameSmSignature).Create(signed).Error
		if nil != err {
			return models.ParsePostgresError("", err)
		}
		return nil
	})

	return e1
}

func (impl *SensorRepo) insertMetricGPS(smx *models.SMExtract, signed *models.SmSignature,
) error {

	sm := &models.SmGPS{
		ID:     uuid.NewV4().String(),
		SignID: signed.ID,
		Position: &models.Point4326{
			Lat: float64(smx.Indicator.Lat),
			Lng: float64(smx.Indicator.Lng),
		},
		CreatedAt: time.Unix(smx.To, 0),
	}

	e1 := impl.tblGPSMetrics().Transaction(func(dbTx *gorm.DB) error {
		err := dbTx.Table(models.TableNameSmGPS).Create(sm).Error
		if nil != err {
			return models.ParsePostgresError("", err)
		}

		err = dbTx.Table(models.TableNameSmSignature).Create(signed).Error
		if nil != err {
			return models.ParsePostgresError("", err)
		}
		return nil
	})

	return e1
}

func (impl *SensorRepo) GetMetrics(req *domain.RGetSM) ([]*models.SmFloat, error) {
	var rs = make([]*models.SmFloat, 0)
	var query = impl.tblFloatMetrics()
	var err = query.Find(&rs).Error
	if nil != err {
		return nil, models.ParsePostgresError("Get metrics", err)
	}
	return rs, nil
}

func (impl *SensorRepo) tblSensor() *gorm.DB {
	return impl.db.Table(models.TableNameSensors)
}

func (impl *SensorRepo) tblFloatMetrics() *gorm.DB {
	return impl.db.Table(models.TableNameSmFloat)
}

func (impl *SensorRepo) tblGPSMetrics() *gorm.DB {
	return impl.db.Table(models.TableNameSmFloat)
}
