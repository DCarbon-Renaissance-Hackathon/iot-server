package repo

import (
	"github.com/Dcarbon/iott-cloud/internal/models"
	"gorm.io/gorm"
)

type SensorRepo struct {
	db *gorm.DB
}

func NewSensorRepo() (*SensorRepo, error) {
	var db = getSingletonDB()
	err := db.AutoMigrate(
		&models.Sensor{},
		&models.SensorMetrict{},
	)
	if nil != err {
		return nil, err
	}

	var impl = &SensorRepo{
		db: db,
	}

	return impl, nil
}

func (impl *SensorRepo) CreateSensor() (*models.Sensor, error) {
	return nil, nil
}

func (impl *SensorRepo) ChangeSensorStatus() (*models.Sensor, error) {
	return nil, nil
}

func (impl *SensorRepo) GetSensor() (*models.Sensor, error) {
	return nil, nil
}

func (impl *SensorRepo) GetSensors() ([]*models.Sensor, error) {
	return nil, nil
}

func (impl *SensorRepo) CreateMetric() (*models.SensorMetrict, error) {
	return nil, nil
}

func (impl *SensorRepo) CreateMetricFromRaw(metric *models.SensorMetricExtract,
) (*models.SensorMetrict, error) {
	return nil, nil
}

func (impl *SensorRepo) GetMetric() ([]*models.SensorMetrict, error) {
	return nil, nil
}
