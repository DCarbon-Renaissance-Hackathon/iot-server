package repo

// type ProposalRepo struct {
// 	db *gorm.DB
// }

// func NewProposalRepo(dbUrl string) (domain.IProposal, error) {
// 	var db, err = dbutils.NewDB(dbUrl)
// 	if nil != err {
// 		return nil, err
// 	}

// 	err = db.AutoMigrate(
// 		&models.Proposal{},
// 	)
// 	if nil != err {
// 		return nil, err
// 	}

// 	var pRepo = &ProposalRepo{
// 		db: db,
// 	}
// 	return pRepo, nil
// }

// func (pRepo *ProposalRepo) Create(v *models.Proposal) error {
// 	err := pRepo.tblProposal().Create(v).Error
// 	return dmodels.ParsePostgresError("Proposal", err)
// }

// func (pRepo *ProposalRepo) GetList(skip, limit, iotId, projectId int64,
// ) ([]*models.Proposal, error) {
// 	var tbl = pRepo.tblProposal()
// 	if skip > 0 {
// 		tbl = tbl.Offset(int(skip))
// 	}

// 	if limit <= 0 || limit > 100 {
// 		limit = 50
// 	}
// 	tbl = tbl.Limit(int(limit))

// 	if iotId > 0 {
// 		tbl = tbl.Where("iot_id = ?", iotId)
// 	}
// 	if projectId > 0 {
// 		tbl = tbl.Where("project_id = ?", projectId)
// 	}
// 	var data = make([]*models.Proposal, 0, limit)
// 	var err = tbl.Find(&data).Error
// 	return data, dmodels.ParsePostgresError("Proposal", err)
// }

// func (pRepo *ProposalRepo) tblProposal() *gorm.DB {
// 	return pRepo.db.Table(models.TableNameProposal)
// }
