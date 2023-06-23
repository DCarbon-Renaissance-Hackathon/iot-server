package repo

import (
	"fmt"
	"log"
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
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
		&models.Sm{},
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
		Status:    dmodels.DeviceStatusRegister,
		CreatedAt: time.Now(),
	}

	var err = impl.tblSensors().Create(sensor).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Create sensor", err)
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
		return nil, dmodels.ParsePostgresError("Change sensor status", err)
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
		return nil, dmodels.ParsePostgresError("Get sensor", err)
	}
	return sensor, nil
}

func (impl *SensorRepo) GetSensorType(req *domain.SensorID,
) (dmodels.SensorType, error) {
	var sensor = &models.Sensor{}
	var query = impl.tblSensors()
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	} else {
		query = query.Where("address = ?", req.Address.String())
	}

	var err = query.Select("id, type").First(sensor).Error
	if nil != err {
		return 0, dmodels.ParsePostgresError("Get sensor", err)
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
		return nil, dmodels.ParsePostgresError("Get sensors", err)
	}
	return sensors, nil
}

func (impl *SensorRepo) CreateSM(req *domain.RCreateSM,
) (*models.SmSignature, error) {
	// log.Println("Create sm payload: ", req)
	sensor, err := impl.GetSensor(&domain.SensorID{Address: req.SensorAddress})
	if nil != err {
		return nil, err
	}

	if sensor.Address.IsEmpty() {
		return nil, dmodels.NewError(dmodels.ECodeSensorHasNoAddress, "SensorAddress is empty")
	}

	var signed = &models.SmSignature{
		ID:        uuid.NewV4().String(),
		IsIotSign: false,
		IotID:     sensor.IotID,
		SensorID:  sensor.ID,
		Data:      req.Data,
		Signed:    req.Signed,
	}

	_, _, err = impl.insertMetric(sensor, signed, *sensor.Address)
	if nil != err {
		return nil, nil
	}

	return signed, nil
}

func (impl *SensorRepo) CreateSensorMetric(req *domain.RCreateSensorMetric,
) (*models.SmSignature, error) {
	sensor, err := impl.GetSensor(&domain.SensorID{ID: req.SensorID})
	if nil != err {
		return nil, err
	}

	var signAddr = sensor.Address
	if req.IsIotSign {
		if !sensor.Address.IsEmpty() {
			return nil, dmodels.NewError(dmodels.ECodeSensorHasAddress, "SensorAddress is not empty")
		}

		signAddr = &req.SignAddress
		if req.IotID != sensor.IotID {
			return nil, dmodels.ErrBadRequest("Iot id and sensor is not mathed")
		}
	}

	var signed = &models.SmSignature{
		ID:        uuid.NewV4().String(),
		IsIotSign: true,
		IotID:     sensor.IotID,
		SensorID:  sensor.ID,
		Data:      req.Data,
		Signed:    req.Signed,
	}

	_, smx, err := impl.insertMetric(sensor, signed, *signAddr)
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
	var query = impl.db.Table(models.TableNameSmSignature + " as tblSign").
		Offset(int(req.Skip)).
		Limit(int(req.Limit)).
		Select("tblSign.id, tblSign.data, tblSign.sensor_id, tblSign.iot_id, tblSign.created_at, sensors.type as sensor_type").
		Joins("JOIN sensors ON tblSign.sensor_id = sensors.id ")

	if req.From > 0 {
		query = query.Where("tblSign.created_at > ?", time.Unix(req.From, 0))
	}

	if req.To > 0 {
		query = query.Where("tblSign.created_at < ?", time.Unix(req.To, 0))
	}

	query = query.Where("tblSign.iot_id = ?", req.IotId)

	if req.SensorId != 0 {
		query = query.Where("tblSign.sensor_id = ?", req.SensorId)
	}

	var err = query.Find(&rs).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Get metrics", err)
	}

	for i, sign := range rs {
		rs[i] = &domain.Metric{
			ID:        sign.ID,
			IotId:     sign.IotId,
			Data:      sign.Data,
			SensorId:  sign.SensorId,
			CreatedAt: sign.CreatedAt,
		}
	}

	return rs, nil
}

func (impl *SensorRepo) GetAggregatedMetrics(req *domain.RSMAggregate,
) ([]*domain.TimeValue, error) {
	return impl.getMetricAggregate(req)
}

func (impl *SensorRepo) GetSignedMetric(req *domain.RGetSM) ([]*models.SmSignature, error) {
	var rs = make([]*models.SmSignature, 0)
	var query = impl.db.Table(models.TableNameSmSignature + " as tblSign").
		Offset(int(req.Skip)).
		Order("created_at asc")

	if req.Limit > 0 {
		query = query.Limit(int(req.Limit))
	}

	if req.From > 0 {
		query = query.Where("tblSign.created_at >= ?", time.Unix(req.From, 0))
	}

	if req.To > 0 {
		query = query.Where("tblSign.created_at < ?", time.Unix(req.To, 0))
	}

	if req.IotId != 0 {
		query = query.Where("tblSign.iot_id = ?", req.IotId)
	} else {
		query = query.Where("tblSign.iot_id != ?", 0)
	}

	if req.SensorId != 0 {
		query = query.Where("tblSign.sensor_id = ?", req.SensorId)
	}

	var err = query.Find(&rs).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Get metrics", err)
	}

	return rs, nil
}

func (impl *SensorRepo) GetMetricsLatest(req *domain.RGetSM) ([]*domain.Metric, error) {
	var rs = make([]*domain.Metric, 0)
	var query = impl.db.Table(models.TableNameSmSignature + " as tblSign").
		Limit(int(req.Limit)).
		Select("tblSign.id, tblSign.data, tblSign.sensor_id, tblSign.iot_id, tblSign.created_at, sensors.type as sensor_type").
		Joins("JOIN sensors ON tblSign.sensor_id = sensors.id ")
	if req.From > 0 {
		query = query.Where("tblSign.created_at > ?", time.Unix(req.From, 0))
	}
	if req.To > 0 {
		query = query.Where("tblSign.created_at < ?", time.Unix(req.To, 0))
	}
	query = query.Where("tblSign.iot_id = ?", req.IotId)
	var err = query.Find(&rs).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Get metrics", err)
	}
	for i, sign := range rs {
		rs[i] = &domain.Metric{
			ID:        sign.ID,
			IotId:     sign.IotId,
			Data:      sign.Data,
			SensorId:  sign.SensorId,
			CreatedAt: sign.CreatedAt,
		}
	}
	return rs, nil
}

func (impl *SensorRepo) insertMetric(sensor *models.Sensor, signed *models.SmSignature, addr dmodels.EthAddress,
) (*models.Sm, *models.SMExtract, error) {
	smx, err := signed.VerifySignature(addr, sensor.Type)
	if nil != err {
		return nil, nil, err
	}

	var data = &models.Sm{
		ID:        uuid.NewV4().String(),
		IotID:     signed.IotID,
		SensorID:  signed.SensorID,
		SignID:    signed.ID,
		Indicator: smx.Indicator,
		CreatedAt: time.Unix(smx.From, 0),
	}
	signed.CreatedAt = data.CreatedAt

	e1 := impl.tblMetrics().Transaction(func(dbTx *gorm.DB) error {
		err := dbTx.Table(models.TableNameSm).Create(data).Error
		if nil != err {
			return dmodels.ParsePostgresError("Save sensor metric data", err)
		}

		err = dbTx.Table(models.TableNameSmSignature).Create(signed).Error
		if nil != err {
			return dmodels.ParsePostgresError("Save sensor metric signature", err)
		}

		return nil
	})
	if nil != err {
		return nil, nil, e1
	}
	return data, smx, nil
}

func (impl *SensorRepo) migrateSM(signed *models.SmSignature,
) (*models.Sm, error) {
	smx, err := signed.ExtractData()
	if nil != err {
		return nil, err
	}

	var data = &models.Sm{
		ID:        uuid.NewV4().String(),
		IotID:     signed.IotID,
		SensorID:  signed.SensorID,
		SignID:    signed.ID,
		Indicator: smx.Indicator,
		CreatedAt: time.Unix(smx.From, 0),
	}

	err = impl.db.Table(models.TableNameSm).Create(data).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("", err)
	}

	return data, nil
}

func (impl *SensorRepo) getMetricAggregate(req *domain.RSMAggregate) ([]*domain.TimeValue, error) {
	var data = make([]*domain.TimeValue, 0)

	var query = impl.tblMetrics().Where(
		"created_at >= ? AND created_at < ? AND iot_id = ? AND sensor_id = ?",
		time.Unix(req.From, 0), time.Unix(req.To, 0), req.IotId, req.SensorId,
	).Group("time").Order("time desc")

	if req.Interval > 0 {
		var trunc = "day"
		if req.Interval == 2 {
			trunc = "month"
		}
		query = query.Select(
			fmt.Sprintf(
				"date_trunc('%s', created_at) as time, SUM (CAST (indicator ->> 'value' as float)) as val",
				trunc,
			),
		)
	} else {
		query = query.Select("created_at as time, CAST (indicator ->> 'value' as float) as val ")
	}

	var err = query.Find(&data).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("", err)
	}
	return data, nil
}

func (impl *SensorRepo) tblSensors() *gorm.DB {
	return impl.db.Table(models.TableNameSensors)
}

func (impl *SensorRepo) tblMetrics() *gorm.DB {
	return impl.db.Table(models.TableNameSm)
}

// func (impl *SensorRepo) insertMetric(smx *models.SMExtract, signed *models.SmSignature, stype dmodels.SensorType,
// ) error {
// 	switch stype {
// 	case dmodels.SensorTypePower:
// 		return impl.insertMetricFloat(smx, signed)
// 	case dmodels.SensorTypeFlow:
// 		return impl.insertMetricFloat(smx, signed)
// 	case dmodels.SensorTypeGPS:
// 		return impl.insertMetricGPS(smx, signed)
// 	}
// 	return dmodels.NewError(dmodels.ECodeSensorInvalidType, "Sensor type is not supported")
// }

// func (impl *SensorRepo) insertMetricFloat(smx *models.SMExtract, signed *models.SmSignature,
// ) error {
// 	sm := &models.SmFloat{
// 		ID:        uuid.NewV4().String(),
// 		SignID:    signed.ID,
// 		Indicator: float64(smx.Indicator.Val),
// 		CreatedAt: time.Unix(smx.To, 0),
// 	}
// 	e1 := impl.tblFloatMetrics().Transaction(func(dbTx *gorm.DB) error {
// 		err := dbTx.Table(models.TableNameSmFloat).Create(sm).Error
// 		if nil != err {
// 			return dmodels.ParsePostgresError("", err)
// 		}
// 		err = dbTx.Table(models.TableNameSmSignature).Create(signed).Error
// 		if nil != err {
// 			return dmodels.ParsePostgresError("", err)
// 		}
// 		return nil
// 	})
// 	return e1
// }

// func (impl *SensorRepo) insertMetricGPS(smx *models.SMExtract, signed *models.SmSignature,
// ) error {
// 	sm := &models.SmGPS{
// 		ID:     uuid.NewV4().String(),
// 		SignID: signed.ID,
// 		Position: &models.Point4326{
// 			Lat: float64(smx.Indicator.Lat),
// 			Lng: float64(smx.Indicator.Lng),
// 		},
// 		CreatedAt: time.Unix(smx.To, 0),
// 	}
// 	e1 := impl.tblGPSMetrics().Transaction(func(dbTx *gorm.DB) error {
// 		err := dbTx.Table(models.TableNameSmGPS).Create(sm).Error
// 		if nil != err {
// 			return dmodels.ParsePostgresError("", err)
// 		}
// 		err = dbTx.Table(models.TableNameSmSignature).Create(signed).Error
// 		if nil != err {
// 			return dmodels.ParsePostgresError("", err)
// 		}
// 		return nil
// 	})
// 	return e1
// }

// func (impl *SensorRepo) tblFloatMetrics() *gorm.DB {
// 	return impl.db.Table(models.TableNameSmFloat)
// }

// func (impl *SensorRepo) tblGPSMetrics() *gorm.DB {
// 	return impl.db.Table(models.TableNameSmFloat)
// }
