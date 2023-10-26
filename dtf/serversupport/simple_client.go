package serversupport

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/serversupport/serverhelper"
)

const TokenExpireDuration time.Duration = 5 * time.Minute
const Issuer = "SimpleInvoker"
const Subject = "pterergate-service"

type SimpleInvoker struct {
	client *http.Client
}

func NewSimpleInvoker() *SimpleInvoker {
	s := &SimpleInvoker{}
	s.client = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxConnsPerHost: 5,
		},
	}

	return s
}

func (s *SimpleInvoker) Post(url string, userName string, requestBody string) error {
	// generate the request json
	commonReq := s.genCommonRequest(requestBody)
	commonReqData, err := json.Marshal(commonReq)
	if err != nil {
		return err
	}

	sign, err := s.sign(userName, commonReq)
	if err != nil {
		return err
	}

	// new a request
	httpReq, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(commonReqData)))
	if err != nil {
		return nil
	}

	// set http headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Accept-Encoding", "gzip")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sign))

	// send the request
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

func (s *SimpleInvoker) genCommonRequest(requestBody string) *serverhelper.CommonRequest {
	req := &serverhelper.CommonRequest{
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

func CalcMsgHash(msg string) string {
	hash := sha256.New()
	hash.Write([]byte(msg))
	bytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(bytes)
	return hashCode
}

func (s *SimpleInvoker) sign(userName string, req *serverhelper.CommonRequest) (string, error) {
	now := time.Now()
	claims := serverhelper.CommonClaims{
		UserName: userName,
		BodyHash: req.Header.BodyHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenExpireDuration)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    Issuer,
			ID:        uuid.NewString(),
			Audience:  jwt.ClaimStrings{"executor", "collector"},
			Subject:   Subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString("")
}
