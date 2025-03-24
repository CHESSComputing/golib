package zenodo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
	"github.com/CHESSComputing/golib/services"
)

var _httpReadRequest, _httpWriteRequest *services.HttpRequest
var Verbose int

func initSrv() {
	if _httpReadRequest == nil {
		_httpReadRequest = services.NewHttpRequest("read", Verbose)
	}
	if _httpWriteRequest == nil {
		_httpWriteRequest = services.NewHttpRequest("write", Verbose)
	}
	_httpReadRequest.GetToken()
	_httpWriteRequest.GetToken()
}

// helper function to check response status
func checkResponse(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		// read response body to get error
		if data, err := io.ReadAll(resp.Body); err == nil {
			msg := fmt.Sprintf("call to upstream service not successfull, %s", string(data))
			return errors.New(msg)
		}
		msg := fmt.Sprintf("call to upstream service fails with code %d", resp.StatusCode)
		return errors.New(msg)
	}
	return nil
}

// CreateRecord provides create record API
func CreateRecord() (int64, error) {
	// init reader/writer and srv config
	initSrv()

	var docId int64
	// create new DOI resource
	rurl := fmt.Sprintf("%s/create", srvConfig.Config.Services.PublicationURL)
	resp, err := _httpWriteRequest.Post(rurl, "application/json", bytes.NewBuffer([]byte{}))
	defer resp.Body.Close()
	if err != nil {
		return docId, err
	}
	if err := checkResponse(resp); err != nil {
		return docId, err
	}

	// capture response and extract document id (did)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return docId, err
	}
	var doc CreateResponse
	err = json.Unmarshal(data, &doc)
	if err != nil {
		return docId, err
	}
	return doc.Id, nil
}

func getBid(did int64) (string, error) {
	rurl := fmt.Sprintf("%s/docs/%d", srvConfig.Config.Services.PublicationURL, did)
	resp, err := _httpReadRequest.Get(rurl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var rec DoiRecord
	data, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(data, &rec)
	if err != nil {
		return "", err
	}
	arr := strings.Split(rec.Links.Bucket, "/")
	bid := arr[len(arr)-1]
	return bid, nil
}

// AddRecord represents add API to zenodo
func AddRecord(docId int64, name string, foxdenRecord any) error {
	initSrv()
	data, err := json.Marshal(foxdenRecord)
	if err != nil {
		return err
	}
	bid, err := getBid(docId)
	if err != nil {
		return err
	}
	rurl := fmt.Sprintf("%s/add/%s/%s", srvConfig.Config.Services.PublicationURL, bid, name)
	resp, err := _httpWriteRequest.Put(rurl, "application/json", bytes.NewBuffer(data))
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if err := checkResponse(resp); err != nil {
		return err
	}
	return nil
}

func UpdateRecord(docId int64, mrec MetaDataRecord) error {
	// init reader/writer and srv config
	initSrv()
	rec := MetaRecord{Metadata: mrec}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	rurl := fmt.Sprintf("%s/update/%d", srvConfig.Config.Services.PublicationURL, docId)
	metaResp, err := _httpWriteRequest.Put(rurl, "application/json", bytes.NewBuffer(data))
	defer metaResp.Body.Close()
	if err != nil {
		return err
	}
	if err := checkResponse(metaResp); err != nil {
		return err
	}
	return nil
}

func PublishRecord(docId int64) (DoiRecord, error) {
	// init reader/writer and srv config
	initSrv()

	var doiRecord DoiRecord

	// publish the record
	rurl := fmt.Sprintf("%s/publish/%d", srvConfig.Config.Services.PublicationURL, docId)
	publishResp, err := _httpWriteRequest.Post(rurl, "application/json", bytes.NewBuffer([]byte{}))
	defer publishResp.Body.Close()
	if err != nil {
		return doiRecord, err
	}
	if err := checkResponse(publishResp); err != nil {
		return doiRecord, err
	}

	// fetch our document
	rurl = fmt.Sprintf("%s/docs/%d", srvConfig.Config.Services.PublicationURL, docId)
	docsResp, err := _httpReadRequest.Get(rurl)
	defer docsResp.Body.Close()
	if err != nil {
		return doiRecord, err
	}
	if err := checkResponse(docsResp); err != nil {
		return doiRecord, err
	}
	data, err := io.ReadAll(docsResp.Body)
	if err != nil {
		return doiRecord, err
	}

	// parse doi record
	err = json.Unmarshal(data, &doiRecord)
	return doiRecord, err
}
