package MaterialsCommons

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
	mcapi "github.com/materials-commons/gomcapi"
)

var mcClient *mcapi.Client

// helper function to get MaterialsCommons client
func getMcClient() {
	if mcClient != nil {
		return
	}
	args := &mcapi.ClientArgs{
		BaseURL: srvConfig.Config.DOI.MaterialsCommons.Url,
		APIKey:  srvConfig.Config.DOI.MaterialsCommons.AccessToken,
	}
	mcClient = mcapi.NewClient(args)
	return
}

// Publish function pulishes FOXDEN dataset with did and description in MaterialsCommons
func Publish(did, description string, record map[string]any, publish bool) (string, string, error) {
	var err error
	var doi, doiLink string
	var projectID, datasetID int

	// get MaterialsCommons client
	getMcClient()

	// find out project ID to use
	projectName := srvConfig.Config.MaterialsCommons.ProjectName
	if projectName == "" {
		projectName = "FOXDEN datasets"
	}
	records, err := mcClient.ListProjects()
	if err != nil {
		log.Println("unable to list projects, error", err)
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
			log.Println("unable to create project, error", err)
			return doi, doiLink, err
		}
		projectID = proj.ID
	}

	// Create a temporary file with out record
	tempFile, err := os.CreateTemp("", "foxden-metadata.json")
	if err != nil {
		log.Println("unable to create temp foxden.json file, error", err)
		return doi, doiLink, err
	}
	defer os.Remove(tempFile.Name())
	content, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		log.Println("unable to marshal record, error", err)
		return doi, doiLink, err
	}
	_, err = tempFile.Write(content)
	if err != nil {
		log.Println("unable to write record, error", err)
		return doi, doiLink, err
	}
	err = tempFile.Close()
	if err != nil {
		log.Println("unable to close temp file, error", err)
		return doi, doiLink, err
	}

	// compose dataset file
	datasetFiles := []mcapi.DatasetFileUpload{
		mcapi.DatasetFileUpload{
			File:        tempFile.Name(),
			Description: "FOXDEN MetaData",
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
		log.Println("unable to deposit dataset, error", err)
		return doi, doiLink, err
	}
	datasetID = ds.ID

	// publish flag instructs that we should make public DOI
	// while MintDOIForDataset and PublishDataset APIs uses publishAsTestDataset flag
	draft := !publish

	// publish deposit within project and dataset ids
	_, err = mcClient.PublishDataset(projectID, datasetID, draft)
	if err != nil {
		log.Println("unable to publish dataset, error", err)
		return doi, doiLink, err
	}
	log.Printf("MaterialsCommons::PublishDataset API: projectID=%v datasetID=%v ds=%+v err=%v draft=%v", projectID, datasetID, ds, err, draft)

	// Mint DOI using our project and dataset ids
	ds, err = mcClient.MintDOIForDataset(projectID, datasetID, draft)
	if err == nil {
		if draft {
			doi = ds.TestDOI
		} else {
			doi = ds.DOI
		}
		mcUrl := strings.Replace(srvConfig.Config.MaterialsCommons.Url, "/api", "", -1)
		if strings.HasSuffix(mcUrl, "/") {
			mcUrl = strings.TrimSuffix(mcUrl, "/")
		}
		doiLink = fmt.Sprintf("%s/dois/%s", mcUrl, doi)
	}
	log.Printf("MaterialsCommons::MintDOIForDataset API: projectID=%v datasetID=%v ds=%+v doi=%v err=%v", projectID, datasetID, ds, doi, err)

	return doi, doiLink, err
}

// FindProjectDatastIDs finds both projectID and datasetID for a given doi
func FindProjectDatasetIDs(doi string) (int, int, error) {
	var projectID, datasetID int

	// get MaterialsCommons client
	getMcClient()

	// find out project ID to use
	projectName := srvConfig.Config.MaterialsCommons.ProjectName
	if projectName == "" {
		projectName = "FOXDEN datasets"
	}
	records, err := mcClient.ListProjects()
	if err != nil {
		log.Println("ERROR: unable to list projects, error", err)
		return projectID, datasetID, err
	}
	for _, r := range records {
		if r.Name == projectName {
			projectID = r.ID
			break
		}
	}
	// get list of datasets within our projectID
	datasets, err := mcClient.ListDatasets(projectID)
	if err != nil {
		log.Println("ERROR: unable to list datasets for projectID", projectID, "error:", err)
		return projectID, datasetID, err
	}
	for _, d := range datasets {
		if d.DOI == doi {
			datasetID = d.ID
		}
	}
	return projectID, datasetID, err
}

// MakePublic implements logic of publishing draft DOI
func MakePublic(doi string) error {
	// get MaterialsCommons client
	getMcClient()

	projectID, datasetID, err := FindProjectDatasetIDs(doi)
	if err != nil {
		log.Println("ERROR: unable to find project and datasetID for doi", doi, "error:", err)
		return err
	}

	// Mint DOI using our project and dataset ids
	draft := false
	_, err = mcClient.MintDOIForDataset(projectID, datasetID, draft)
	return err
}
