package imdg

import "encoding/json"

type JSONSerializer struct {

}

func (srz JSONSerializer) Encode(data interface {}) ([]byte, error){
	return json.Marshal(data)
}

func (srz JSONSerializer) Decode(data []byte) (map[string]interface {}, error) {
	result := make(map[string]interface {})
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
