package serversupport

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/danenmao/pterergate-dtf/dtf/errordef"
	"github.com/danenmao/pterergate-dtf/dtf/serversupport/serverhelper"
	"github.com/danenmao/pterergate-dtf/dtf/taskmodel"
)

type CollectorInvoker struct {
	ServerHost string
	ServerPort uint16
	URI        string
	url        string
}

// return an invoker function
// for executor to invoke collector
func (c *CollectorInvoker) GetInvoker() taskmodel.CollectorInvoker {
	// Host:Port
	// POST /collector
	c.url = fmt.Sprintf("%s:%d%s", c.ServerHost, c.ServerPort, c.URI)
	return func(results []taskmodel.SubtaskResult) error {
		return c.invoker(results)
	}
}

func (c *CollectorInvoker) invoker(results []taskmodel.SubtaskResult) error {
	body := CollectorRequestBody{
		Results: results,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return errordef.ErrOperationFailed
	}

	request := serverhelper.CommonRequest{}
	request.Body = string(data)
	return c.post(&request)
}

func (c *CollectorInvoker) post(request *serverhelper.CommonRequest) error {
	req, err := http.NewRequest(http.MethodPost, c.url, nil)
	if err != nil {
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	response := serverhelper.CommonResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Header.Code != errordef.Error_Msg_Success {
		return errors.New(response.Header.Message)
	}

	return nil
}
