package repo

import (
	"time"

	"github.com/Dcarbon/go-shared/dmodels"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"github.com/Dcarbon/iott-cloud/internal/rss"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepo() (domain.IProject, error) {
	var db = rss.GetDB()
	err := db.AutoMigrate(
		&models.Project{},
		&models.ProjectImage{},
		&models.ProjectDescription{},
		&models.ProjectSpecs{},
	)
	if nil != err {
		return nil, err
	}

	var pp = &projectRepo{
		db: db,
	}
	return pp, nil
}

func (pRepo *projectRepo) Create(req *domain.RProjectCreate,
) (*models.Project, error) {

	var project = req.ToProject()
	var e1 = pRepo.tblProject().Transaction(func(dbTx *gorm.DB) error {
		err := dbTx.Table(models.TableNameProject).Create(project).Error
		if nil != err {
			return dmodels.ParsePostgresError("Create project", err)
		}

		return nil
	})

	if nil != e1 {
		return nil, e1
	}

	return project, nil
}

func (pRepo *projectRepo) UpdateDesc(req *domain.RProjectUpdateDesc,
) (*models.ProjectDescription, error) {
	var desc = req.ToProjectDesc()

	var err = pRepo.tblProjectDesc().
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "project_id"}, {Name: "language"}},
				UpdateAll: true,
			},
			clause.Insert{},
		).
		Create(desc).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Update project desc", err)
	}
	return desc, nil
}

func (pRepo *projectRepo) UpdateSpecs(req *domain.RProjectUpdateSpecs,
) (*models.ProjectSpecs, error) {
	var spec = req.ToProjectSpecs()

	var err = pRepo.tblProjectSpec().
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "project_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"specs", "updated_at"}),
			},
		).Create(spec).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Update project desc", err)
	}
	return spec, nil
}

func (pRepo *projectRepo) GetById(id int64, lang string) (*models.Project, error) {
	var project = &models.Project{}
	var query = pRepo.tblProject().Where("id = ?", id).
		Preload("Images", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("project_id, image")
		}).
		Preload("Specs")
	if lang != "" {
		query.Preload("Descs", func(tx *gorm.DB) *gorm.DB {
			return tx.Where("language = ?", lang)
		})
	}

	var err = query.First(project).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Project", err)
	}

	return project, nil
}

func (pRepo *projectRepo) GetList(filter *domain.RProjectFilter,
) ([]*models.Project, error) {
	var tbl = pRepo.tblProject().Offset(filter.Skip)
	if filter.Limit > 0 {
		tbl = tbl.Limit(filter.Limit)
	}

	if filter.Owner != "" {
		tbl = tbl.Where("owner = ?", filter.Owner)
	}

	var data = make([]*models.Project, 0)
	var err = tbl.Find(&data).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Project", err)
	}
	return data, nil
}

func (pRepo *projectRepo) GetByBB(min, max *models.Point4326, owner string,
) ([]*models.Project, error) {
	var data = make([]*models.Project, 0)
	var err = pRepo.tblProject().
		Where(
			"ST_WITHIN(pos, ST_MakeEnvelope(?, ?, ?, ?, 4326))",
			min.Lng, min.Lat, max.Lng, max.Lat).
		Find(&data).Error
	return data, dmodels.ParsePostgresError("Project", err)
}

func (pRepo *projectRepo) GetByID(id int64) (*models.Project, error) {
	var data = &models.Project{}
	var err = pRepo.tblProject().Where("id = ?", id).First(data).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("Project", err)
	}
	return data, nil
}

func (pRepo *projectRepo) ChangeStatus(id string, status models.ProjectStatus,
) error {
	var err = pRepo.tblProject().
		Where("id = ?", id).
		Update("status", status).
		Error
	return dmodels.ParsePostgresError("Project", err)
}

func (pRepo *projectRepo) GetOwner(projectId int64) (string, error) {
	var owner = ""
	var err = pRepo.tblProject().
		Where("id = ?", projectId).
		Pluck("owner", &owner).Error
	if nil != err {
		return "", dmodels.ParsePostgresError("Get owner ", err)
	}
	return owner, nil
}

func (pRepo *projectRepo) AddImage(req *domain.RProjectAddImage,
) (*models.ProjectImage, error) {
	var img = &models.ProjectImage{
		ID:        0,
		ProjectID: req.ProjectID,
		Image:     req.ImgPath,
		CreatedAt: time.Now(),
	}
	var err = pRepo.tblImage().Create(img).Error
	if nil != err {
		return nil, dmodels.ParsePostgresError("AddImage ", err)
	}
	return img, nil
}

func (pRepo *projectRepo) tblProject() *gorm.DB {
	return pRepo.db.Table(models.TableNameProject)
}

func (pRepo *projectRepo) tblProjectDesc() *gorm.DB {
	return pRepo.db.Table(models.TableNameProjectDesc)
}

func (pRepo *projectRepo) tblProjectSpec() *gorm.DB {
	return pRepo.db.Table(models.TableNameProjectSpecs)
}

func (pRepo *projectRepo) tblImage() *gorm.DB {
	return pRepo.db.Table(models.TableNameProjectImage)
}
