package MaterialCommons

import (
	"fmt"

	srvConfig "github.com/CHESSComputing/golib/config"
	mcapi "github.com/materials-commons/gomcapi"
)

var mcClient *mcapi.Client

func getMcClient() {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	if mcClient != nil {
		return
	}
	args := &mcapi.ClientArgs{
		BaseURL: srvConfig.Config.DOI.URL,
		APIKey:  srvConfig.Config.DOI.AccessToken,
	}
	mcClient = mcapi.NewClient(args)
	return
}

func Publish(did, description string) (string, error) {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	getMcClient()
	var err error
	var doi string

	var projectID, datasetID int
	projectName := srvConfig.Config.DOI.ProjectName
	records, err := mcClient.ListProjects()
	if err != nil {
		return doi, err
	}
	for _, r := range records {
		if r.Name == projectName {
			projectID = r.ID
		}
	}
	// if project does not exist we'll create it
	if projectID == 0 {
		req := mcapi.CreateProjectRequest{
			Name:        "FOXDEN datasets",
			Description: "FOXDEN repository of public datasets",
			Summary:     "FOXDEN repository of public dataset",
		}
		proj, err := mcClient.CreateProject(req)
		if err != nil {
			return doi, err
		}
		projectID = proj.ID
	}

	// create new deposit
	name := fmt.Sprintf("FOXDEN dataset %s", did)
	summary := "FOXDEN dataset"
	deposit := mcapi.DepositDatasetRequest{
		Metadata: mcapi.DatasetMetadata{
			Name:        name,
			Description: description,
			Summary:     summary,
		},
	}
	ds, err := mcClient.DepositDataset(projectID, deposit)
	if err != nil {
		return doi, err
	}
	datasetID = ds.ID

	// publish deposit
	_, err = mcClient.PublishDataset(projectID, datasetID)
	if err != nil {
		return doi, err
	}
	ds, err = mcClient.MintDOIForDataset(projectID, datasetID)
	if err == nil {
		doi = ds.DOI
	}
	return doi, err
}
