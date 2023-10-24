package serversupport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/serversupport/serverhelper"
)

type SimpleHTTPClient struct {
	client *http.Client
}

func NewSimpleHTTPClient() *SimpleHTTPClient {
	s := &SimpleHTTPClient{}
	s.client = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxConnsPerHost: 1,
		},
	}

	return s
}

func (s *SimpleHTTPClient) Post(url string, userName string, requestBody string) error {
	// commonReq json
	commonReq := s.genCommonRequest(requestBody)
	commonReqData, err := json.Marshal(commonReq)
	if err != nil {
		return err
	}

	// send the request
	httpReq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(commonReqData)))
	if err != nil {
		return nil
	}

	// set http headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Accept-Encoding", "gzip")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", commonReq.Header.Sign))

	rsp, err := s.client.Do(httpReq)
	if err != nil {
		return err
	}

	// parse the response
	if rsp.StatusCode != 200 {
		return fmt.Errorf("error HTTP status %d", rsp.StatusCode)
	}

	defer rsp.Body.Close()
	respBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	// parse the common response
	commonResp := serverhelper.CommonResponse{}
	err = json.Unmarshal(respBody, &commonResp)
	if err != nil {
		return err
	}

	if commonResp.Header.Code != errordef.Error_Msg_Success {
		return errors.New(commonResp.Header.Message)
	}

	return nil
}

func (s *SimpleHTTPClient) genCommonRequest(requestBody string) *serverhelper.CommonRequest {
	req := &serverhelper.CommonRequest{
		Body: requestBody,
	}

	// fill header
	req.Header.Version = "1.0"
	req.Header.RequestId = uuid.NewString()
	req.Header.Timestamp = strconv.FormatInt(time.Now().Unix(), 16)
	req.Header.Module = ""
	req.Header.Action = ""

	// sign the request body
	req.Header.Sign, _ = s.sign(req)
	return req
}

func (s *SimpleHTTPClient) sign(req *serverhelper.CommonRequest) (string, error) {
	return "", nil
}
