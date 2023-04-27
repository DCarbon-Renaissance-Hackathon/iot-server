package repo

import (
	"fmt"
	"log"
	"time"

	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/rss"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SensorRepo struct {
	db      *gorm.DB
	opCache domain.IOperator
}

func NewSensorRepo() (*SensorRepo, error) {
	var db = rss.GetDB()
	err := db.AutoMigrate(
		&models.Sensor{},
		&models.SmSignature{},
		&models.SmFloat{},
		&models.SmGPS{},
		// &models.Sm{},
	)
	if nil != err {
		return nil, err
	}

	var impl = &SensorRepo{
		db: db,
	}

	return impl, nil
}

func (impl *SensorRepo) SetOperatorCache(op domain.IOperator) {
	impl.opCache = op
}

func (impl *SensorRepo) CreateSensor(req *domain.RCreateSensor,
) (*models.Sensor, error) {
	var sensor = &models.Sensor{
		ID:        0,
		IotID:     req.IotID,
		Address:   &req.Address,
		Type:      req.Type,
		Status:    models.DeviceStatusRegister,
		CreatedAt: time.Now(),
	}

	var err = impl.tblSensors().Create(sensor).Error
	if nil != err {
		return nil, models.ParsePostgresError("Create sensor", err)
	}

	return sensor, nil
}

func (impl *SensorRepo) ChangeSensorStatus(req *domain.RChangeSensorStatus,
) (*models.Sensor, error) {
	var sensor = &models.Sensor{}
	var err = impl.tblSensors().
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
	var query = impl.tblSensors()
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

func (impl *SensorRepo) GetSensorType(req *domain.SensorID,
) (models.SensorType, error) {
	var sensor = &models.Sensor{}
	var query = impl.tblSensors()
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	} else {
		query = query.Where("address = ?", req.Address.String())
	}
	var err = query.Select("id, type").First(sensor).Error
	if nil != err {
		return 0, models.ParsePostgresError("Get sensor", err)
	}
	return sensor.Type, nil
}

func (impl *SensorRepo) GetSensors(req *domain.RGetSensors,
) ([]*models.Sensor, error) {
	var sensors = make([]*models.Sensor, 0, req.Limit)
	var query = impl.tblSensors().Offset(req.Skip)
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
	log.Println("Create sm payload: ", req)
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

	if impl.opCache != nil {
		err = impl.opCache.SetStatus(&domain.ROpSetStatus{
			Id:     sensor.IotID,
			Status: models.OpStatusActived,
		})
		if nil != err {
			log.Println("Save iot status error: ", err)
		}

		_, err = impl.opCache.ChangeMetrics(&domain.RChangeMetric{
			IotId:    sensor.IotID,
			SensorId: sensor.ID,
			Metric:   smx.Indicator,
		}, sensor.Type)
		if nil != err {
			log.Println("Save sensor metric cache error: ", err)
		}
	}

	return signed, nil
}

func (impl *SensorRepo) GetMetrics(req *domain.RGetSM) ([]*domain.Metric, error) {
	var rs = make([]*domain.Metric, 0)
	var query = impl.db.Table(models.TableNameSm+" as tblSM").
		Select("tblSM.id, tblSM.indicator, tblSM.created_at, tblSign.iot_id, tblSign.sensor_id ").
		Joins(fmt.Sprintf("JOIN %s as tblSign ON tblSign.id = tblSM.sign_id", models.TableNameSmSignature)).
		Where("iot_id = ?", req.IotId).
		Order("tblSM.created_at asc")
	if req.From > 0 {
		query = query.Where("tblSM.created_at > ?", time.Unix(req.From, 0))
	}

	if req.To > 0 {
		query = query.Where("tblSM.created_at < ?", time.Unix(req.To, 0))
	}

	var err = query.Find(&rs).Error
	if nil != err {
		return nil, models.ParsePostgresError("Get metrics", err)
	}

	return rs, nil
}

// func (impl *SensorRepo) insertMetric(smx *models.SMExtract, signed *models.SmSignature, stype models.SensorType,
// ) error {
// 	var data = &models.Sm{
// 		ID:        uuid.NewV4().String(),
// 		SignID:    signed.ID,
// 		Indicator: &smx.Indicator,
// 		CreatedAt: time.Unix(smx.To, 0),
// 	}
// 	e1 := impl.tblMetrics().Transaction(func(dbTx *gorm.DB) error {
// 		err := dbTx.Table(models.TableNameSm).Create(data).Error
// 		if nil != err {
// 			return models.ParsePostgresError("Save sensor metric data", err)
// 		}
// 		err = dbTx.Table(models.TableNameSmSignature).Create(signed).Error
// 		if nil != err {
// 			return models.ParsePostgresError("Save sensor metric signature", err)
// 		}
// 		return nil
// 	})
// 	return e1
// }

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
}

func (impl *SensorRepo) insertMetricFloat(smx *models.SMExtract, signed *models.SmSignature,
) error {
	sm := &models.SmFloat{
		ID:        uuid.NewV4().String(),
		SignID:    signed.ID,
		Indicator: float64(smx.Indicator.Val),
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

func (impl *SensorRepo) tblSensors() *gorm.DB {
	return impl.db.Table(models.TableNameSensors)
}

// func (impl *SensorRepo) tblMetrics() *gorm.DB {
// 	return impl.db.Table(models.TableNameSm)
// }

func (impl *SensorRepo) tblFloatMetrics() *gorm.DB {
	return impl.db.Table(models.TableNameSmFloat)
}

func (impl *SensorRepo) tblGPSMetrics() *gorm.DB {
	return impl.db.Table(models.TableNameSmFloat)
}
