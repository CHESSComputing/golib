package MaterialCommons

import (
	"encoding/json"
	"fmt"
	"os"

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
func Publish(did, description string, record any) (string, string, error) {
	var err error
	var doi, doiLink string
	var projectID, datasetID int
	if srvConfig.Config == nil {
		srvConfig.Init()
	}

	// get MaterialCommons client
	getMcClient()

	// find out project ID to use
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

	// Create a temporary file with out record
	tempFile, err := os.CreateTemp("", "foxden-*.json")
	if err != nil {
		return doi, doiLink, err
	}
	defer os.Remove(tempFile.Name())
	content, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return doi, doiLink, err
	}
	_, err = tempFile.Write(content)
	if err != nil {
		return doi, doiLink, err
	}
	err = tempFile.Close()
	if err != nil {
		return doi, doiLink, err
	}

	// compose dataset file
	datasetFiles := []mcapi.DatasetFileUpload{
		mcapi.DatasetFileUpload{
			File:        tempFile.Name(),
			Description: description,
			Directory:   "foxden",
		},
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

	// publish deposit within project and dataset ids
	_, err = mcClient.PublishDataset(projectID, datasetID)
	if err != nil {
		return doi, doiLink, err
	}

	// Mint DOI using our project and dataset ids
	ds, err = mcClient.MintDOIForDataset(projectID, datasetID)
	if err == nil {
		doi = ds.DOI
		doiLink = fmt.Sprintf("https://doi.org/%s", doi)
	}
	return doi, doiLink, err
}
