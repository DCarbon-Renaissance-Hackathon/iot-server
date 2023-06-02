package repo

import (
	"strings"
	"time"

	"github.com/Dcarbon/go-shared/libs/esign"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/rss"
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

func (ip *iotRepo) Create(req *domain.RIotCreate) (*models.IOTDevice, error) {
	var iot = &models.IOTDevice{
		ID:       0,
		Project:  req.Project,
		Address:  req.Address,
		Type:     req.Type,
		Status:   models.DeviceStatusRegister,
		Position: *req.Position,
	}

	var err = ip.tblIOT().Create(iot).Error
	if nil != err {
		return nil, models.ParsePostgresError("IOT", err)
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
		return nil, models.ParsePostgresError("IOT", err)
	}

	return iot, err
}

func (ip *iotRepo) GetIOT(id int64) (*models.IOTDevice, error) {
	var iot = &models.IOTDevice{}
	var err = ip.tblIOT().Where("id = ?", id).First(iot).Error
	if nil != err {
		return iot, models.ParsePostgresError("IOT", err)
	}
	return iot, nil
}

func (ip *iotRepo) GetIOTByAddress(addr models.EthAddress) (*models.IOTDevice, error) {
	var iot = &models.IOTDevice{}
	var err = ip.tblIOT().Where("address = ?", &addr).First(iot).Error
	if nil != err {
		return iot, models.ParsePostgresError("IOT", err)
	}

	return iot, nil
}

func (ip *iotRepo) GetByBB(min, max *models.Point4326,
) ([]*models.IOTDevice, error) {
	var iots = make([]*models.IOTDevice, 0)
	var err = ip.tblIOT().
		Where(
			"ST_WITHIN(pos, ST_MakeEnvelope(?, ?, ?, ?, 4326))",
			min.Lng, min.Lat, max.Lng, max.Lat,
		).
		Find(&iots).Error
	return iots, models.ParsePostgresError("IOT", err)
}

func (ip *iotRepo) GetIOTStatus(iotAddr string) models.DeviceStatus {
	var device = &models.IOTDevice{}
	var err = ip.tblIOT().
		Select("status").
		Where("address = ?", strings.ToLower(iotAddr)).
		First(&device).Error
	if nil != err {
		device.Status = models.DeviceStatusReject
	}
	return device.Status
}

func (ip *iotRepo) CreateMint(mint *models.MintSign,
) error {
	mint.IOT = strings.ToLower(mint.IOT)

	var iot, err = ip.GetIOTByAddress(models.EthAddress(mint.IOT))
	if nil != err {
		return err
	}

	if iot.Status < models.DeviceStatusRegister {
		return models.NewError(models.ECodeIOTNotAllowed, "IOT is not allow")
	}

	err = mint.Verify(ip.dMinter)
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
		err = models.ParsePostgresError("", ip.tblSign().Create(mint).Error)
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
		err = models.ParsePostgresError("", ip.tblSign().Create(mint).Error)
	} else {
		err = models.NewError(models.ECodeIOTInvalidNonce, "Invalid nonce")
	}

	return err
}

func (ip *iotRepo) GetMintSigns(req *domain.RIotGetMintSignList,
) ([]*models.MintSign, error) {
	var iot, err = ip.GetIOT(req.IotId)
	if nil != err {
		return nil, err
	}

	var signeds = make([]*models.MintSign, 0)
	err = ip.tblSign().
		Where(
			"created_at > ? AND created_at < ? AND  iot = ?",
			time.Unix(req.From, 0), time.Unix(req.To, 0), iot.Address,
		).
		Order("id asc").
		Find(&signeds).
		Error
	if nil != err {
		return nil, models.ParsePostgresError("Get mint sign", err)
	}
	return signeds, nil
}

func (ip *iotRepo) tblIOT() *gorm.DB {
	return ip.db.Table(models.TableNameIOT)
}

func (ip *iotRepo) tblSign() *gorm.DB {
	return ip.db.Table(models.TableNameMintSign)
}
