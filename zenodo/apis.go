package zenodo

// based on information collected from
// https://developers.zenodo.org/?shell#quickstart-upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
)

var Verbose int

// CreateRecord provides create record API
func CreateRecord(payload []byte) (CreateResponse, error) {
	/*
	 curl --request POST 'https://zenodo.org/api/deposit/depositions?access_token=<KEY>' \
	 --header 'Content-Type: application/json'  \
	 --data-raw '{}'
	*/
	var response CreateResponse
	// create new deposit
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/deposit/depositions?access_token=%s", zurl, token)
	req, err := http.NewRequest("POST", rurl, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: unable to post request to zenodo", err)
		return response, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if Verbose > 0 {
		log.Println("createDocument received response", string(data))
	}
	err = json.Unmarshal(data, &response)
	return response, err
}

// AddRecord represents add API to zenodo
/*
	curl --upload-file readme.md --request PUT
	'https://zenodo.org/api/files/50b47f75-c97d-47c6-af11-caa6e967c1d5/readme.md?access_token=<KEY>'
*/
func AddRecord(docId int64, fileName string, foxdenRecord any) error {
	data, err := json.Marshal(foxdenRecord)
	if err != nil {
		return err
	}
	records, err := DoiRecords(docId)
	if err != nil {
		return err
	}
	if len(records) != 1 {
		return errors.New("Too many DOI records")
	}
	rec := records[0]
	arr := strings.Split(rec.Links.Bucket, "/")
	bucket := arr[len(arr)-1]

	// create new deposit
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/files/%s/%s?access_token=%s", zurl, bucket, fileName, token)
	if Verbose > 0 {
		log.Println("request", rurl)
	}

	// place HTTP request to zenodo upstream server
	req, err := http.NewRequest("PUT", rurl, bytes.NewReader(data))
	if err != nil {
		return err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if Verbose > 0 {
		log.Println("updateDocument received response", string(data))
	}
	return err
}

// UpdateRecord updates Zenodo records with our metadata
/*
   # add mandatory metadata to our publication
   curl -v -X PUT "https://zenodo.org/api/deposit/depositions/<ID>?access_token=<TOKEN>" \
		   -H "Content-type: application/json" -d@meta1.json

	{
		"metadata": {
			"publication_type": "article",
			"upload_type":"publication",
			"description":"This is a test",
			"keywords": ["bla", "foo"],
			"title":"Test"
		}
	}
*/
func UpdateRecord(docId int64, mrec MetaDataRecord) error {
	rec := MetaRecord{Metadata: mrec}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/deposit/depositions/%d?access_token=%s", zurl, docId, token)

	// place HTTP request to zenodo upstream server
	req, err := http.NewRequest("PUT", rurl, bytes.NewReader(data))
	if err != nil {
		log.Println("ERROR: unable to PUT request to zenodo", err)
		return err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if Verbose > 0 {
		log.Println("updateDocument received response", string(data))
	}
	return err
}

// PublishRecord publishes docId record in Zenodo
// curl -v -X POST "https://zenodo.org/api/deposit/depositions/<ID>/actions/publish?access_token=<TOKEN>"
func PublishRecord(docId int64) (DoiRecord, error) {
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/deposit/depositions/%d/actions/publish?access_token=%s", zurl, docId, token)
	var record DoiRecord

	// place HTTP request to zenodo upstream server
	req, err := http.NewRequest("POST", rurl, nil)
	if err != nil {
		log.Println("ERROR: unable to POST request to zenodo", err)
		return record, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if Verbose > 0 {
		log.Println("publishDocument received response", string(data))
	}
	err = json.Unmarshal(data, &record)
	return record, err
}

// DoiRecords returns list of Zenodo DOI records
func DoiRecords(docId int64) ([]DoiRecord, error) {
	/*
	 curl 'https://zenodo.org/api/deposit/depositions?access_token=<KEY>'
	 curl 'https://zenodo.org/api/deposit/depositions/<123>?access_token=<KEY>'
	*/
	var records []DoiRecord
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/deposit/depositions?access_token=%s", zurl, token)
	if docId != 0 {
		rurl = fmt.Sprintf("%s/deposit/depositions/%d?access_token=%s", zurl, docId, token)
	}
	if Verbose > 0 {
		log.Println("request", rurl)
	}
	resp, err := http.Get(rurl)
	if err != nil {
		log.Println("ERROR: in GET request", err)
		return records, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if docId == 0 {
		if err := dec.Decode(&records); err != nil {
			log.Println("ERROR: unable to decode JSON response", err)
			return records, err
		}
		return records, nil
	} else {
		var record DoiRecord
		if err := dec.Decode(&record); err != nil {
			log.Println("ERROR: unable to decode JSON response", err)
			return records, err
		}
		records = append(records, record)
	}
	return records, nil
}

// GetRecordId extracts from our doi record id (last part of doi)
func GetRecordId(doi string) string {
	arr := strings.Split(doi, ".")
	recid := arr[len(arr)-1]
	return recid
}

// DepositRecords returns list of Zenodo Deposit records
/*
 curl 'https://zenodo.org/api/deposit/depositions?access_token=<KEY>'
 curl 'https://zenodo.org/api/deposit/depositions/<123>?access_token=<KEY>' where 123 is record id (last part of doi)
*/
func DepositRecords(doi string) ([]DepositRecord, error) {
	var records []DepositRecord
	// extract from our doi record id
	recid := GetRecordId(doi)
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/deposit/depositions?access_token=%s&q=recid:%s", zurl, token, recid)
	if Verbose > 0 {
		log.Println("request", rurl)
	}
	resp, err := http.Get(rurl)
	if err != nil {
		log.Println("ERROR: in GET request", err)
		return records, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&records); err != nil {
		log.Println("ERROR: unable to decode JSON response", err)
		return records, err
	}
	// if doi is provided filter out records with this doi
	if doi != "" {
		var out []DepositRecord
		for _, r := range records {
			if r.MetaDataResponseRecord.PrereserveDoi.Doi == doi {
				out = append(out, r)
			}
		}
		return out, nil
	}
	return records, nil
}

// MakePublic implements logic of publishing draft DOI
/*
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
     -X POST "https://zenodo.org/api/deposit/depositions/1234567/actions/publish"

*/
func MakePublic(rid string) error {
	zurl := srvConfig.Config.DOI.Zenodo.Url
	token := srvConfig.Config.DOI.Zenodo.AccessToken
	rurl := fmt.Sprintf("%s/deposit/depositions/%s/actions/publish", zurl, rid)
	if Verbose > 0 {
		log.Println("request", rurl)
	}
	log.Println("MakePublic POST to", rurl)
	req, err := http.NewRequest("POST", rurl, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: unable to post request to zenodo", err)
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if Verbose > 0 {
		log.Println("MakePublic received response", string(data))
	}
	if resp.StatusCode >= 400 && resp.StatusCode < 200 {
		log.Println("MakePublic received response", string(data))
		return errors.New(fmt.Sprintf("fail to make public DOI record, status code=%v", resp.StatusCode))
	}
	return err
}
