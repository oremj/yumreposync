package aws

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"os"

	"github.com/oremj/awsauth"
)

type STSResult struct {
	Credentials *Credentials `xml:"GetSessionTokenResult>Credentials"`
}

type Credentials struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}

func TemporaryCreds() (*Credentials, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	qs := url.Values{}
	qs.Add("Version", "2011-06-15")
	qs.Add("Action", "GetSessionToken")
	qs.Add("DurationSeconds", "900")

	req, err := http.NewRequest("GET", "https://sts.amazonaws.com/?"+qs.Encode(), nil)
	if err != nil {
		return nil, err
	}

	awsReq, err := awsauth.NewAWSRequest(req, accessKey, secretAccessKey)
	if err != nil {
		return nil, err
	}
	awsReq.Sign()

	c := new(http.Client)
	resp, err := c.Do(awsReq.Request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := new(STSResult)
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(res)
	if err != nil {
		return nil, err
	}

	return res.Credentials, nil
}
