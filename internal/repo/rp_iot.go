package repo

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/go-shared/ecodes"
	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/rss"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type iotRepo struct {
	db      *gorm.DB
	dMinter *esign.ERC712
}

func NewIOTRepo(dMinter *esign.ERC712,
) (domain.IIot, error) {
	var db = rss.GetDB()
	err := db.AutoMigrate(
		&models.IOTDevice{},
		&models.MintSign{},
		&models.Minted{},
		// &models.Metric{},
	)
	if nil != err {
		return nil, err
	}

	var ip = &iotRepo{
		db:      db,
		dMinter: dMinter,
	}
	return ip, nil
}

func (ip *iotRepo) Create(req *domain.RIotCreate,
) (*models.IOTDevice, error) {
	var iot = &models.IOTDevice{
		ID:       0,
		Project:  req.Project,
		Address:  req.Address,
		Type:     req.Type,
		Status:   dmodels.DeviceStatusRegister,
		Position: *req.Position,
	}

	var err = ip.tblIOT().Create(iot).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("IOT", err)
	}
	return iot, nil
}

func (ip *iotRepo) ChangeStatus(req *domain.RIotChangeStatus,
) (*models.IOTDevice, error) {
	var iot = &models.IOTDevice{}
	var err = ip.tblIOT().
		Model(iot).
		Clauses(clause.Returning{}).
		Where("id = ?", req.IotId).
		Update("status", req.Status).
		Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("IOT", err)
	}

	return iot, err
}

func (ip *iotRepo) Update(req *domain.RIotUpdate,
) (*models.IOTDevice, error) {
	var iot = &models.IOTDevice{}
	var err = ip.tblIOT().
		Model(iot).
		Clauses(clause.Returning{}).
		Where("id = ?", req.IotId).
		Updates(map[string]interface{}{
			"position": req.Position,
			// "updated_at": time.Now(),
		}).
		Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("IOT", err)
	}

	return iot, err
}

func (ip *iotRepo) GetIot(id int64) (*models.IOTDevice, error) {
	var iot = &models.IOTDevice{}
	var err = ip.tblIOT().Where("id = ?", id).First(iot).Error
	if nil != err {
		return iot, dmodels.ParsePostgresError("IOT", err)
	}
	return iot, nil
}

func (ip *iotRepo) GetIots(req *domain.RIotGetList,
) ([]*models.IOTDevice, error) {
	var data = make([]*models.IOTDevice, 0)
	var err = ip.queryGetIots(req).Find(&data).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("IOT", err)
	}
	return data, nil
}

func (ip *iotRepo) GetIotPositions(req *domain.RIotGetList,
) ([]*domain.PositionId, error) {
	var locs = make([]*domain.PositionId, 0)
	var err = ip.queryGetIots(req).Select("id, position").Find(&locs).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("IOT", err)
	}
	return locs, nil
}

func (ip *iotRepo) GetIotByAddress(addr dmodels.EthAddress,
) (*models.IOTDevice, error) {
	var iot = &models.IOTDevice{}
	var err = ip.tblIOT().Where("address = ?", &addr).Find(iot).Error
	if nil != err {
		return iot, dmodels.ParsePostgresError("IOT", err)
	}

	return iot, nil
}

func (ip *iotRepo) GetIOTStatus(iotAddr string,
) dmodels.DeviceStatus {
	var device = &models.IOTDevice{}
	var err = ip.tblIOT().
		Select("status").
		Where("address = ?", strings.ToLower(iotAddr)).
		First(&device).Error
	if nil != err {
		device.Status = dmodels.DeviceStatusReject
	}
	return device.Status
}

func (ip *iotRepo) CreateMint(req *domain.RIotMint,
) error {
	if req.Nonce <= 0 {
		return dmodels.ErrInvalidNonce()
	}

	newAmount, e1 := dmodels.NewBigNumberFromHex(req.Amount)
	if nil != e1 {
		return e1
	}

	req.Iot = strings.ToLower(req.Iot)
	iot, e1 := ip.GetIotByAddress(dmodels.EthAddress(req.Iot))
	if nil != e1 {
		return e1
	}

	if iot.Status < dmodels.DeviceStatusRegister {
		return dmodels.NewError(ecodes.IOTNotAllowed, "IOT is not allow")
	}

	var mint = &models.MintSign{
		ID:        0,
		Nonce:     req.Nonce,
		Amount:    req.Amount,
		IotId:     iot.ID,
		Iot:       req.Iot,
		R:         req.R,
		S:         req.S,
		V:         req.V,
		CreatedAt: time.Now(),
	}

	e1 = mint.Verify(ip.dMinter)
	if nil != e1 {
		return e1
	}

	var latest = make([]*models.MintSign, 0, 1)
	e1 = ip.tblSign().
		Where("iot = ?", mint.Iot).
		Order("created_at desc").
		Limit(1).
		Find(&latest).Error
	if nil != e1 {
		return dmodels.ParsePostgresError("", e1)
	}

	if len(latest) == 0 {
		latest = append(latest, &models.MintSign{})
	}

	if latest[0].Nonce == mint.Nonce || latest[0].Nonce+1 == mint.Nonce {
		oldAmount, e1 := dmodels.NewBigNumberFromHex(latest[0].Amount)
		if nil != e1 {
			oldAmount = dmodels.NewBigNumber(0)
		}

		var incAmount = big.NewInt(0).Sub(newAmount.Int, oldAmount.Int)
		var minted = &models.Minted{
			ID:     uuid.NewV4().String(),
			IotId:  iot.ID,
			Carbon: incAmount.Int64(),
		}

		return ip.db.Transaction(func(dbTx *gorm.DB) error {
			if latest[0].Nonce+1 == mint.Nonce {
				err := dbTx.Table(models.TableNameMintSign).Create(mint).Error
				if nil != err {
					return dmodels.ParsePostgresError("", err)
				}
			} else {
				err := dbTx.Table(models.TableNameMintSign).
					Where("id = ?", latest[0].ID).
					Updates(map[string]interface{}{
						"nonce":      mint.Nonce,
						"amount":     mint.Amount,
						"r":          mint.R,
						"s":          mint.S,
						"v":          mint.V,
						"updated_at": time.Now(),
					}).Error
				if nil != err {
					dmodels.ParsePostgresError("", err)
				}
			}

			err := dbTx.Table(models.TableNameMinted).Create(minted).Error
			if nil != err {
				return dmodels.ParsePostgresError("", err)
			}
			return nil
		})

	}
	return dmodels.ErrInvalidNonce()
}

func (ip *iotRepo) GetMintSigns(req *domain.RIotGetMintSignList,
) ([]*models.MintSign, error) {
	var iot, err = ip.GetIot(req.IotId)
	if nil != err {
		return nil, err
	}

	var signeds = make([]*models.MintSign, 0)
	var query = ip.tblSign().
		Where(
			"updated_at > ? AND updated_at < ? AND  iot = ?",
			time.Unix(req.From, 0), time.Unix(req.To, 0), iot.Address,
		)

	if req.Sort > 0 {
		query = query.Order("updated_at desc")
	} else {
		query = query.Order("updated_at asc")
	}

	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}
	err = query.Find(&signeds).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Get mint sign", err)
	}
	return signeds, nil
}

func (ip *iotRepo) GetMinted(req *domain.RIotGetMintedList,
) ([]*models.Minted, error) {
	var tz = "Asia/Ho_Chi_Minh"
	var query = ip.db.Table(models.TableNameMinted).
		Where(
			"created_at > ? AND created_at < ? AND iot_id = ? ",
			time.Unix(req.From, 0), time.Unix(req.To, 0), req.IotId,
		)
	if req.Interval > 0 {
		var group = "day"
		if req.Interval == 2 {
			group = "month"
		}

		query = query.Raw(
			fmt.Sprintf(`SELECT date_trunc('%s', created_at, ?) as ca, sum(carbon) as carbon
							FROM minted
							WHERE created_at > ? AND created_at < ? and iot_id = ?
							GROUP BY ca
							`, group),
			tz, time.Unix(req.From, 0).Format(time.RFC3339), time.Unix(req.To, 0).Format(time.RFC3339), req.IotId,
		)
	} else {
		query = query.Select("created_at as ca, carbon").Order("ca asc")
	}

	var data = make([]*aggMinted, 0)
	var err = query.Find(&data).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("", err)
	}

	var rs = make([]*models.Minted, len(data))
	for i, it := range data {
		rs[i] = &models.Minted{
			CreatedAt: it.Ca,
			Carbon:    it.Carbon,
		}
	}

	return rs, nil
}

func (ip *iotRepo) CountIot(req *domain.RIotCount) (int64, error) {
	var count = int64(0)
	var query = ip.tblIOT()
	var err = query.Count(&count).Error
	if nil != err {
		return 0, dmodels.ParsePostgresError("Count iot", err)
	}
	return count, nil
}

func (ip *iotRepo) IsIotActived(req *domain.RIsIotActiced,
) (bool, error) {
	var count = int64(0)
	var err = ip.db.Table(models.TableNameMinted).
		Where(
			"created_at >= ? AND created_at < ? AND iot_id = ?",
			time.Unix(req.From, 0), time.Unix(req.To, 0), req.IotId,
		).Count(&count).Error
	if nil != err {
		return false, dmodels.ParsePostgresError("Check iot is actived", err)
	}
	return count > 0, nil
}

func (ip *iotRepo) queryGetIots(req *domain.RIotGetList) *gorm.DB {
	var query = ip.tblIOT()
	if req.ProjectId != 0 {
		query = query.Where("project = ?", req.ProjectId)
	}

	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	return query
}

func (ip *iotRepo) tblIOT() *gorm.DB {
	return ip.db.Table(models.TableNameIOT)
}

func (ip *iotRepo) tblSign() *gorm.DB {
	return ip.db.Table(models.TableNameMintSign)
}

// func (ip *iotRepo) GetByBB(min, max *models.Point4326,
// ) ([]*models.IOTDevice, error) {
// 	var iots = make([]*models.IOTDevice, 0)
// 	var err = ip.tblIOT().
// 		Where(
// 			"ST_WITHIN(pos, ST_MakeEnvelope(?, ?, ?, ?, 4326))",
// 			min.Lng, min.Lat, max.Lng, max.Lat,
// 		).
// 		Find(&iots).Error
// 	return iots, dmodels.ParsePostgresError("IOT", err)
// }

type aggMinted struct {
	Ca     time.Time
	Carbon int64
}
