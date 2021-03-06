package facex

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
)

type FacexInput struct {
	Data []*FacexInputItem `json:"data"`
}

type FacexInputItem struct {
	URI       string            `json:"uri"`
	Attribute map[string]string `json:"attribute"`
}

func NewFaceBase64(dat []byte) string {
	return "data:application/octet-stream;base64," + base64.StdEncoding.EncodeToString(dat)
}

func NewFacexInput(uri, id string) FacexInput {
	return FacexInput{
		Data: []*FacexInputItem{
			&FacexInputItem{
				URI: uri,
				Attribute: map[string]string{
					"id": id,
				},
			},
		},
	}
}

type SearchInput struct {
	Data map[string]string `json:"data"`
}

type SearchResult struct {
	Message string            `json:"message"`
	Result  *ResultDetections `json:"result"`
}

type ResultDetections struct {
	Detections []*ResultValue `json:"detections"`
}

type ResultValue struct {
	Value *SearchResultValue `json:"value"`
}

type SearchResultValue struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type ListGroupResult struct {
	Code 	int 		`json:"code"`
	Message string 		`json:"message"`
	Result 	[]*ListGroupItem `json:"result"`
}

type ListGroupItem struct{
	Id		string				`json:"id"`
	Value	*ListGroupItemValue	`json:"value"`
}

type ListGroupItemValue struct{
	Name	string	`json:"name"`
}


func NewSearchResult(data []byte) (*SearchResult, error) {
	var ret SearchResult

	err := json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func NewListGroupResult(data []byte) (*ListGroupResult, error){
	var ret ListGroupResult

	err := json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (this *SearchResult) IsOK(threshold ...float64) bool {
	tv := 0.50 // a experience value
	if len(threshold) > 0 {
		tv = threshold[0]
	}

	return this.Result != nil &&
		this.Result.Detections != nil &&
		len(this.Result.Detections) > 0 &&
		this.Result.Detections[0].Value.Score > tv
}

func (this *SearchResult) Name() (ret string) {
	if len(this.Result.Detections) == 0 {
		return
	}

	return this.Result.Detections[0].Value.Name
}

func (this *SearchResult) Score() (ret float64) {
	if len(this.Result.Detections) == 0 {
		return
	}

	return this.Result.Detections[0].Value.Score
}

func NewSearchInput(uri string) *SearchInput {
	return &SearchInput{
		Data: map[string]string{
			"uri": uri,
		},
	}
}

func toPayload(in interface{}) io.Reader {
	data, _ := json.Marshal(in)

	return bytes.NewBuffer(data)
}
