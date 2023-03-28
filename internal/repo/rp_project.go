package repo

import (
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepo() (domain.IProject, error) {
	var db, err = getSingletonDB()
	if nil != err {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Project{},
		&models.ProjectDescription{},
		&models.ProjectSpec{},
	)
	if nil != err {
		return nil, err
	}

	var pp = &projectRepo{
		db: db,
	}
	return pp, nil
}

func (pRepo *projectRepo) Create(project *models.Project) error {
	var err = pRepo.tblProject().Create(project).Error

	return models.ParsePostgresError("Project", err)
}

func (pRepo *projectRepo) UpdateDesc(req *models.ProjectDescription,
) (*models.ProjectDescription, error) {
	var err = pRepo.tblProjectDesc().
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "project_id"}, {Name: "language"}},
				UpdateAll: true,
			},
			clause.Insert{},
		).
		Create(req).Error
	if nil != err {
		return nil, models.ParsePostgresError("Update project desc", err)
	}
	return req, nil
}

func (pRepo *projectRepo) UpdateSpec(req *models.ProjectSpec,
) (*models.ProjectSpec, error) {
	var err = pRepo.tblProjectSpec().
		Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "project_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"specs", "updated_at"}),
			},
		).Create(req).Error
	if nil != err {
		return nil, models.ParsePostgresError("Update project desc", err)
	}
	return req, nil
}

func (pRepo *projectRepo) GetById(id int64, lang string, withSpec bool) (*models.Project, error) {
	var project = &models.Project{}
	var query = pRepo.tblProject().Where("id = ?", id)
	if lang != "" {
		query.Preload("Descs", func(tx *gorm.DB) *gorm.DB {
			return tx.Where("language = ?", lang)
		})
	}
	if withSpec {
		query.Preload("Specs")
	}
	var err = query.First(project).Error

	return project, models.ParsePostgresError("Project", err)
}

func (pRepo *projectRepo) GetList(filter *domain.RFilterProject,
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
		return nil, models.ParsePostgresError("Project", err)
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
	return data, models.ParsePostgresError("Project", err)
}

func (pRepo *projectRepo) GetByID(id int64) (*models.Project, error) {
	var data = &models.Project{}
	var err = pRepo.tblProject().Where("id = ?", id).First(data).Error
	if nil != err {
		return nil, models.ParsePostgresError("Project", err)
	}
	return data, nil
}

func (pRepo *projectRepo) ChangeStatus(id string, status models.ProjectStatus,
) error {
	var err = pRepo.tblProject().
		Where("id = ?", id).
		Update("status", status).
		Error
	return models.ParsePostgresError("Project", err)
}

func (pRepo *projectRepo) GetOwner(projectId int64) (string, error) {
	var owner = ""
	var err = pRepo.tblProject().
		Where("id = ?", projectId).
		Pluck("owner", &owner).Error
	if nil != err {
		return "", models.ParsePostgresError("Get owner ", err)
	}
	return owner, nil
}

func (pRepo *projectRepo) tblProject() *gorm.DB {
	return pRepo.db.Table(models.TableNameProject)
}

func (pRepo *projectRepo) tblProjectDesc() *gorm.DB {
	return pRepo.db.Table(models.TableNameProjectDesc)
}

func (pRepo *projectRepo) tblProjectSpec() *gorm.DB {
	return pRepo.db.Table(models.TableNameProjectSpec)
}
