package repo

import (
	"testing"

	"github.com/Dcarbon/go-shared/libs/utils"
	"github.com/Dcarbon/iott-cloud/internal/domain"
	"github.com/Dcarbon/iott-cloud/internal/models"
)

var pRepoTest domain.IProject

func init() {
	var err error

	// err := rss.InitResource(dbUrlTest, redisUrl)
	// utils.PanicError("", err)

	pRepoTest, err = NewProjectRepo()
	utils.PanicError("", err)
}

func TestProjectCreate(t *testing.T) {
	var p = &domain.RProjectCreate{
		Owner: adminAddr,
		Location: &models.Point4326{
			Lat: 21.015462,
			Lng: 105.804904,
		},
		// Status: models.ProjectStatusRegister,
	}
	_, err := pRepoTest.Create(p)
	utils.PanicError("", err)
	utils.Dump("Project created: ", p)
}

func TestProjectUpdateDesc(t *testing.T) {
	var desc = &models.ProjectDescription{
		ProjectID: 1,
		Language:  "en",
		Name:      "Test description",
		Desc:      "Test desc en",
	}
	var rs, err = pRepoTest.UpdateDesc(desc)
	utils.PanicError("", err)
	utils.Dump("", rs)
}

func TestProjectUpdateSpec(t *testing.T) {
	var spec = &models.ProjectSpec{
		ProjectID: 1,
		Specs: models.MapSFloat{
			"s": 51.0,
		},
	}
	var rs, err = pRepoTest.UpdateSpec(spec)
	utils.PanicError("TestProjectUpdateSpec", err)
	utils.Dump("TestProjectUpdateSpec", rs)
}

func TestProjectGetByID(t *testing.T) {
	var rs, err = pRepoTest.GetById(1, "")
	utils.PanicError("TestProjectGetByID", err)
	utils.Dump("TestProjectGetByID", rs)
}

func TestProjectGetList(t *testing.T) {

}

func TestProjectGetByBB(t *testing.T) {

}

func TestProjectChangeStatus(t *testing.T) {

}
