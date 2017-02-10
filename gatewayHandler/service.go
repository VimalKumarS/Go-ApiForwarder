package Utility

import (
	"fmt"
	"io"
	"net/http"
)

type serviceClient interface {
	SendCommand(httpVerb, urlRoute string, body io.Reader, header http.Header) *http.Response
}

//ServiceWebClient allows communication to registry
type ServiceWebClient struct {
	URL string
}

//SendCommand sends the command to the designated service
func (client ServiceWebClient) SendCommand(httpVerb, urlRoute string, body io.Reader, header http.Header) *http.Response {
	httpclient := &http.Client{}

	url := fmt.Sprintf("%s%s", client.URL, urlRoute)
	req, _ := http.NewRequest(httpVerb, url, body)
	req.Header = header
	resp, err := httpclient.Do(req)

	if err != nil {
		Log.Printf("Errored when sending request to the server: %s\n", err.Error())
		return nil
	}
	return resp
}
