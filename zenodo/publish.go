package zenodo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

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
	if srvConfig.Config == nil {
		srvConfig.Init()
	}
	_httpReadRequest.GetToken()
	_httpWriteRequest.GetToken()
}

// helper function to check response status
func checkResponse(resp *http.Response) error {
	if resp.StatusCode != 200 {
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

func UpdateRecord(docId int64, mrec MetaDataRecord) error {
	// init reader/writer and srv config
	initSrv()
	data, err := json.Marshal(mrec)
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
	if err != nil || (publishResp.StatusCode < 200 || publishResp.StatusCode >= 400) {
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
