package repo

import (
	"github.com/Dcarbon/iott-cloud/domain"
	"github.com/Dcarbon/iott-cloud/libs/dbutils"
	"github.com/Dcarbon/iott-cloud/models"
	"gorm.io/gorm"
)

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(dbUrl string) (domain.IProject, error) {
	var db, err = dbutils.NewDB(dbUrl)
	if nil != err {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Project{},
	)
	if nil != err {
		return nil, err
	}

	var pp = &projectRepo{
		db: db,
	}
	return pp, nil
}

func (pp *projectRepo) Create(project *models.Project) error {
	var err = pp.tblProject().Create(project).Error

	return models.ParsePostgresError("Project", err)
}

func (pp *projectRepo) GetById(id int64) (*models.Project, error) {
	var project = &models.Project{}
	var err = pp.tblProject().
		Where("id = ?", id).
		First(project).
		Error

	return project, models.ParsePostgresError("Project", err)
}

func (pp *projectRepo) GetList(skip, limit int64, owner string,
) ([]*models.Project, error) {
	var tbl = pp.tblProject()
	if skip > 0 {
		tbl = tbl.Offset(int(skip))
	}
	if limit > 0 {
		tbl = tbl.Limit(int(limit))
	}
	if owner != "" {
		tbl = tbl.Where("owner = ?", owner)
	}

	var data = make([]*models.Project, 0)
	var err = tbl.Find(&data).Error
	if nil != err {
		return nil, models.ParsePostgresError("Project", err)
	}
	return data, nil
}

func (pp *projectRepo) GetByBB(min, max *models.Point4326, owner string,
) ([]*models.Project, error) {
	var data = make([]*models.Project, 0)
	var err = pp.tblProject().
		Where(
			"ST_WITHIN(pos, ST_MakeEnvelope(?, ?, ?, ?, 4326))",
			min.Lng, min.Lat, max.Lng, max.Lat).
		Find(&data).Error
	return data, models.ParsePostgresError("Project", err)
}

func (pp *projectRepo) GetByID(id int64) (*models.Project, error) {
	var data = &models.Project{}
	var err = pp.tblProject().Where("id = ?", id).First(data).Error
	if nil != err {
		return nil, models.ParsePostgresError("Project", err)
	}
	return data, nil
}

func (pp *projectRepo) ChangeStatus(id string, status models.ProjectStatus,
) error {
	var err = pp.tblProject().
		Where("id = ?", id).
		Update("status", status).
		Error
	return models.ParsePostgresError("Project", err)
}

func (pp *projectRepo) tblProject() *gorm.DB {
	return pp.db.Table(models.TableNameProject)
}
