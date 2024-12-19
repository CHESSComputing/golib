package MaterialCommons

import (
	"fmt"

	srvConfig "github.com/CHESSComputing/golib/config"
	mcapi "github.com/materials-commons/gomcapi"
)

var mcClient *mcapi.Client

// helper function to get MaterialCommons client
func getMcClient() {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	if mcClient != nil {
		return
	}
	args := &mcapi.ClientArgs{
		BaseURL: srvConfig.Config.MaterialCommons.Url,
		APIKey:  srvConfig.Config.MaterialCommons.AccessToken,
	}
	mcClient = mcapi.NewClient(args)
	return
}

// Publish function pulishes FOXDEN dataset with did and description in MaterialCommons
func Publish(did, description string) (string, string, error) {
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	getMcClient()
	var err error
	var doi, doiLink string

	var projectID, datasetID int
	projectName := srvConfig.Config.MaterialCommons.ProjectName
	if projectName == "" {
		projectName = "FOXDEN datasets"
	}
	records, err := mcClient.ListProjects()
	if err != nil {
		return doi, doiLink, err
	}
	for _, r := range records {
		if r.Name == projectName {
			projectID = r.ID
		}
	}
	// if project does not exist we'll create it
	if projectID == 0 {
		req := mcapi.CreateProjectRequest{
			Name:        projectName,
			Description: fmt.Sprintf("%s: repository of public datasets", projectName),
			Summary:     fmt.Sprintf("%s: repository of public datasets", projectName),
		}
		proj, err := mcClient.CreateProject(req)
		if err != nil {
			return doi, doiLink, err
		}
		projectID = proj.ID
	}
	// add dataset files to our deposit
	// TODO: obtain files from FOXDEN, e.g. get meta-data record
	var files []string
	var datasetFiles []mcapi.DatasetFileUpload
	for _, fname := range files {
		dir := "foxden"  // TODO: decide on directory placement
		desc := "foxden" // TODO: decide on description of the file
		f := mcapi.DatasetFileUpload{
			File:        fname,
			Description: desc,
			Directory:   dir,
		}
		datasetFiles = append(datasetFiles, f)
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
		Files: datasetFiles,
	}
	ds, err := mcClient.DepositDataset(projectID, deposit)
	if err != nil {
		return doi, doiLink, err
	}
	datasetID = ds.ID

	// publish deposit
	_, err = mcClient.PublishDataset(projectID, datasetID)
	if err != nil {
		return doi, doiLink, err
	}
	ds, err = mcClient.MintDOIForDataset(projectID, datasetID)
	if err == nil {
		doi = ds.DOI
		doiLink = fmt.Sprintf("https://doi.org/%s", doi)
	}
	return doi, doiLink, err
}
