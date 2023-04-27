package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/rss"
	"github.com/go-redis/redis/v8"
)

const (
	keyIotStatus = "iot_status"     // HashTable(HSET, HGET)
	keyIotMetric = "iot_metrics_%d" // HashTable(HSET, HGET)
)

type OperatorRepo struct {
	redis *redis.Client
}

func NewOperatorRepo() (*OperatorRepo, error) {
	var op = &OperatorRepo{
		redis: rss.GetRedis(),
	}
	return op, nil
}

func (op *OperatorRepo) SetStatus(req *domain.ROpSetStatus) error {
	var stt = &models.OpIotStatus{
		Id:     req.Id,
		Status: req.Status,
		Latest: time.Now().Unix(),
	}
	raw, err := json.Marshal(stt)
	if nil != err {
		return models.ErrInternal(err)
	}

	_, err = op.redis.HSet(
		context.TODO(),
		keyIotStatus,
		fmt.Sprintf("%d", req.Id), string(raw),
	).Result()
	if nil != err {
		return models.ErrInternal(err)
	}

	return nil
}

func (op *OperatorRepo) GetStatus(iotId int64) (*models.OpIotStatus, error) {
	str, err := op.redis.HGet(context.TODO(), keyIotStatus, fmt.Sprintf("%d", iotId)).Result()
	if err == redis.Nil || str == "" {
		return &models.OpIotStatus{
			Id:     iotId,
			Status: models.OpStatusInactived,
		}, nil
	}
	if nil != err {
		return nil, err
	}

	var stt = &models.OpIotStatus{}
	err = json.Unmarshal([]byte(str), stt)
	if nil != err {
		return nil, err
	}

	return stt, nil
}

func (op *OperatorRepo) ChangeMetrics(req *domain.RChangeMetric, sensorType models.SensorType,
) (*models.OpSensorMetric, error) {
	var metric = &models.OpSensorMetric{
		Id:     req.SensorId,
		Type:   sensorType,
		Metric: req.Metric,
		Latest: time.Now().Unix(),
	}
	var raw, err = json.Marshal(metric)
	if nil != err {
		return nil, models.ErrInternal(err)
	}

	err = op.redis.HSet(
		context.TODO(),
		fmt.Sprintf(keyIotMetric, req.IotId),
		fmt.Sprintf("%d", req.SensorId), string(raw)).Err()
	if nil != err {
		return nil, models.ErrInternal(err)
	}

	return metric, nil
}

func (op *OperatorRepo) GetMetrics(iotId int64) (*domain.RsGetMetrics, error) {
	data, err := op.redis.HGetAll(context.TODO(), fmt.Sprintf(keyIotMetric, iotId)).Result()
	if nil != err {
		if err == redis.Nil {
			return &domain.RsGetMetrics{}, nil
		}
		return nil, err
	}

	var metrics = make([]*models.OpSensorMetric, 0, len(data))
	for _, v := range data {
		m := &models.OpSensorMetric{}
		err = json.Unmarshal([]byte(v), m)
		if nil != err {
			log.Println("Marshall metric error: ", err)
		} else {
			metrics = append(metrics, m)
		}

	}

	return &domain.RsGetMetrics{
			Id:      iotId,
			Metrics: metrics,
		},
		nil
}
