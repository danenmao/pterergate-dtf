package serverhelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/internal/msgsigner"
)

const TokenExpireDuration time.Duration = 5 * time.Minute
const Issuer = "SimpleInvoker"
const Subject = "pterergate-service"

type SimpleInvoker struct {
	client *http.Client
	signer *msgsigner.MsgSigner
}

func NewSimpleInvoker() *SimpleInvoker {
	s := &SimpleInvoker{}
	s.signer = msgsigner.NewMsgSigner()
	s.client = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxConnsPerHost: 5,
		},
	}

	return s
}

func (s *SimpleInvoker) Post(url string, userName string, requestBody string) (string, error) {
	// generate the request json
	commonReq := s.genCommonRequest(requestBody)
	commonReqData, err := json.Marshal(commonReq)
	if err != nil {
		return "", err
	}

	// sign the request body
	sign, err := s.sign(userName, commonReq)
	if err != nil {
		return "", err
	}

	// new a request
	httpReq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(commonReqData)))
	if err != nil {
		return "", nil
	}

	// set http headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Accept-Encoding", "gzip")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sign))

	// send the request
	rsp, err := s.client.Do(httpReq)
	if err != nil {
		glog.Warning("Failed to send a request to the server: ", err)
		return "", err
	}

	// parse the response
	if rsp.StatusCode != 200 {
		return "", fmt.Errorf("error HTTP status %d", rsp.StatusCode)
	}

	defer rsp.Body.Close()
	respBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	// parse the common response
	commonResp := CommonResponse{}
	err = json.Unmarshal(respBody, &commonResp)
	if err != nil {
		return "", err
	}

	if commonResp.Header.Code != errordef.Error_Msg_Success {
		return commonResp.Body, errors.New(commonResp.Header.Message)
	}

	return commonResp.Body, nil
}

func (s *SimpleInvoker) genCommonRequest(requestBody string) *CommonRequest {
	req := &CommonRequest{
		Body: requestBody,
	}

	// fill header
	req.Header.Version = "1.0"
	req.Header.RequestId = uuid.NewString()
	req.Header.Timestamp = strconv.FormatInt(time.Now().Unix(), 16)
	req.Header.Module = ""
	req.Header.Action = ""

	// calc the body hash
	req.Header.BodyHash = CalcMsgHash(requestBody)

	return req
}

func (s *SimpleInvoker) sign(userName string, req *CommonRequest) (string, error) {
	msg := CommonMessage{
		UserName: userName,
		BodyHash: req.Header.BodyHash,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	audience := []string{"executor", "collector"}
	return s.signer.Sign(Issuer, Subject, audience, string(data), TokenExpireDuration)
}
