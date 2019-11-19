// Package movizor package provides access to http://MoVizor.ru API
// which provides access for GSM geo-position services of russian telecommunications operators.
// Beeline, MTS, Megafon, Tele2.
//
// As soon as MoVizor provides service only in Russia all documentation will be in russian.
package movizor

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// API - это клиент к API Мовизора. Сервиса определения гео-координат на основе GSM сервиса.
type API struct {
	Endpoint string
	Project  string
	Token    string
	Client   *http.Client
	IsDebug  bool
	Retries  uint
}

// NewMovizorAPIWithEndpoint создает экземпляр Movizor API.
// Может быть указан нестандартный адрес сервиса МоВизора
// на случай, если такой появится.
func NewMovizorAPIWithEndpoint(endp string, prj string, token string) (*API, error) {
	api := &API{
		Endpoint: endp,
		Project:  prj,
		Token:    token,
		Client:   &http.Client{},
		IsDebug:  false,
		Retries:  1,
	}
	return api, nil
}

// NewMovizorAPI создает экземпляр Movizor API для стандартного адреса сервера.
func NewMovizorAPI(prj string, token string) (*API, error) {
	return NewMovizorAPIWithEndpoint(DefaultAPIMovizorEndpoint, prj, token)
}

// MakeRequest делает запрос в Movizor API с указанным действием и параметрами.
// Все методы API вызывают MakeRequest.
func (api *API) MakeRequest(action string, params url.Values) (APIResponse, error) {
	// MakeRequest itself
	endpAction := api.prepareRequest(action, params)

	req, err := http.NewRequest("GET", endpAction, nil)
	if err != nil {
		return APIResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var resp *http.Response

	// retry
	max := int(api.Retries)
	for i := 0; i < max; i++ {
		resp, err = api.Client.Do(req)
		if err == nil {
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		if api.IsDebug {
			log.Printf("ERROR: request: %s\nerror: %s", req.URL, err)
		}

		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return APIResponse{}, err
	}

	defer resp.Body.Close()

	// Response handling
	var apiResp APIResponse
	bytes, err := api.decodeAPIResponse(resp.Body, &apiResp)
	if err != nil {
		return APIResponse{}, err
	}

	if apiResp.Result == "success" {
		if api.IsDebug {
			log.Printf("INFO: request: %s\nresponse: %s", req.URL, bytes)
		}
		return apiResp, nil
	}

	err = fmt.Errorf("movizor API returns error on request: %s - %s",
		apiResp.ErrorCode, apiResp.ErrorText)
	if api.IsDebug {
		log.Printf("ERROR: request: %s\nresponse: %s", req.URL, bytes)
	}

	return apiResp, err
}

func (api *API) prepareRequest(action string, params url.Values) string {
	endpAction := fmt.Sprintf(fmt.Sprint(api.Endpoint, APIMovizorEndpointSuffix), api.Project, action)
	if params == nil {
		params = url.Values{}
	}
	params.Add("key", api.Token)
	uri, _ := url.Parse(endpAction)
	uri.RawQuery = params.Encode()
	endpAction = uri.String()
	return endpAction
}

func (api *API) decodeAPIResponse(responseBody io.Reader, resp *APIResponse) (_ []byte, err error) {
	if !api.IsDebug {
		dec := json.NewDecoder(responseBody)
		err = dec.Decode(resp)
		return
	}

	// if logging is on, read response body
	data, err := ioutil.ReadAll(responseBody)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return
	}

	return data, nil
}

// GetBalance возвращает текущее состояние баланса и установленные тарифы
// для всех видов мониторинга.
func (api *API) GetBalance() (Balance, error) {
	resp, err := api.MakeRequest("balance", nil)
	if err != nil {
		return Balance{}, err
	}

	var b Balance
	err = json.Unmarshal(resp.Data, &b)
	if err != nil {
		return Balance{}, err
	}
	return b, nil
}

// AddObject подключает абонента к мониторингу.
func (api *API) AddObject(o Object, oo *ObjectOptions) (APIResponse, error) {
	return api.AddObjectToSlave(o, oo, 0)
}

// AddObjectToSlave подключает абонента к мониторингу в подчиненный кабинет по ID этого кабинета.
// ID кабинета тоже самое, что и "Номер клиента" указанные в правом верхнем углу кабинета клиента.
func (api *API) AddObjectToSlave(o Object, oo *ObjectOptions, slaveID uint64) (APIResponse, error) {
	v, err := o.values()
	if err != nil {
		return APIResponse{}, err
	}

	if oo != nil {
		if err := oo.addValuesTo(&v); err != nil {
			return APIResponse{}, err
		}
	}

	if slaveID != 0 {
		v.Add("account", strconv.FormatUint(slaveID, 10))
	}

	resp, err := api.MakeRequest("object_add", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetObjectInfo возвращает информацию о ранее добавленном абоненте с наиболее полной
// информацией, включая все опции, с которыми добавлялся объект.
func (api *API) GetObjectInfo(o Object) (ObjectInfo, error) {
	v, err := o.values()
	if err != nil {
		return ObjectInfo{}, err
	}

	resp, err := api.MakeRequest("object_get", v)
	if err != nil {
		return ObjectInfo{}, err
	}

	var oi ObjectInfo
	err = json.Unmarshal(resp.Data, &oi)
	if err != nil {
		return ObjectInfo{}, err
	}

	return oi, nil
}

// EditObject производит изменение опций мониторинга ранее добавленного абонента.
func (api *API) EditObject(o Object, oo *ObjectOptions) (APIResponse, error) {
	return api.EditObjectWithActivate(o, oo, false)
}

// EditObjectWithActivate проивзодит изменение опций мониторинга ранее добавленного абонента
// с опцией немедленной активации или на следующие сутки.
//
//		api.EditObjectWithActivate(object, options, true) // немедленная активация
//		api.EditObjectWithActivate(object, options, false) // активация на следующие сутки
func (api *API) EditObjectWithActivate(o Object, oo *ObjectOptions, activate bool) (APIResponse, error) {
	v, err := o.values()
	if err != nil {
		return APIResponse{}, err
	}

	if oo != nil {
		err := oo.addValuesTo(&v)
		if err != nil {
			return APIResponse{}, err
		}
	}

	if activate {
		v.Add("activate", "1")
	}

	resp, err := api.MakeRequest("object_edit", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetObjects возвращает список абонентов, добавленных в мониторинг с их статусом
// и текущим местоположением.
func (api *API) GetObjects() (ObjectsWithStatus, error) {
	resp, err := api.MakeRequest("object_list", nil)
	if err != nil {
		return ObjectsWithStatus{}, err
	}

	var o ObjectsWithStatus
	err = json.Unmarshal(resp.Data, &o)
	if err != nil {
		return ObjectsWithStatus{}, err
	}

	return o, nil
}

// GetFilteredObjects возвращает список объектов отфильтрованных переданным параметром-функцией
func (api *API) GetFilteredObjects(filter func(ObjectStatus) bool) (ObjectsWithStatus, error) {
	ows, err := api.GetObjects()
	if err != nil {
		return ObjectsWithStatus{}, nil
	}

	result := ows[:0]
	for _, os := range ows {
		if !filter(os) {
			result = append(result, os)
		}
	}
	for i := len(result); i < len(ows); i++ {
		ows[i] = ObjectStatus{}
	}
	return result, nil
}

// DeleteObject отключает и удаляет абонента из системы мониторинга.
func (api *API) DeleteObject(o Object) (APIResponse, error) {
	v, err := o.values()
	if err != nil {
		return APIResponse{}, err
	}

	resp, err := api.MakeRequest("object_delete", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// ReactivateObject производит повторное подключение к системе абонента, если сработало автоматическое отключение.
// Невозможно повторно подключить ранее удаленный объект мониторинга.
func (api *API) ReactivateObject(o Object) (APIResponse, error) {
	v, err := o.values()
	if err != nil {
		return APIResponse{}, err
	}

	resp, err := api.MakeRequest("object_reactivate", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// CancelTariffChangeObject отменяет переход на новый тариф со следующего дня. Если с помощтю EditObject
// и без автоматической активации меняется тариф, то эту смену можно отменить.
func (api *API) CancelTariffChangeObject(o Object) (APIResponse, error) {
	v, err := o.values()
	if err != nil {
		return APIResponse{}, err
	}

	resp, err := api.MakeRequest("object_cancel_tariff", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetObjectLastPosition возвращает информацию о последнем зафиксированном в системе местоположении.
func (api *API) GetObjectLastPosition(o Object) (Position, error) {
	v, err := o.values()
	if err != nil {
		return Position{}, err
	}

	resp, err := api.MakeRequest("pos_last", v)
	if err != nil {
		return Position{}, err
	}

	var lp Position
	err = json.Unmarshal(resp.Data, &lp)
	if err != nil {
		return Position{}, err
	}

	return lp, nil
}

// GetObjectPositions возвращает информацию о всех координатах абонента.
// По умолчанию выдаются последние 1000 записей.
func (api *API) GetObjectPositions(o Object, rpo *RequestPositionsOptions) (Positions, error) {
	v, err := o.values()
	if err != nil {
		return Positions{}, err
	}

	if rpo != nil {
		err := rpo.addValuesTo(&v)
		if err != nil {
			return Positions{}, err
		}
	}
	resp, err := api.MakeRequest("pos_list", v)
	if err != nil {
		return Positions{}, err
	}

	var op Positions
	err = json.Unmarshal(resp.Data, &op)
	if err != nil {
		return Positions{}, err
	}

	return op, nil
}

// RequestPosition выполняет запрос на определение местоположения.
// Используется для определения координат вручную (TariffManual).
// Однако можно использовать и с другими тарифами.
// RequestPosition возвращает ID запроса, который передается в GetRequestedPosition
// для получения координат объекта.
func (api *API) RequestPosition(o Object) (PositionRequest, error) {
	v, err := o.values()
	if err != nil {
		return PositionRequest{}, err
	}

	resp, err := api.MakeRequest("pos_request", v)
	if err != nil {
		return PositionRequest{}, err
	}

	var pr PositionRequest
	err = json.Unmarshal(resp.Data, &pr)
	if err != nil {
		return PositionRequest{}, err
	}

	return pr, nil
}

// GetRequestedPosition получает информацию о сделанном запросе на определение
// местоположения по его идентификатору, который получен методом RequestPosition.
func (api *API) GetRequestedPosition(pr PositionRequest) (Position, error) {
	resp, err := api.MakeRequest("pos_get", pr.values())
	if err != nil {
		return Position{}, err
	}

	var p Position
	err = json.Unmarshal(resp.Data, &p)
	if err != nil {
		return Position{}, err
	}

	return p, nil
}

// GetObjectsPositions возвращает список объектов с их местоположением и ETA
// (estimated time of arrival)
func (api *API) GetObjectsPositions() (ObjectPositions, error) {
	resp, err := api.MakeRequest("pos_objects", nil)
	if err != nil {
		return ObjectPositions{}, err
	}

	var op ObjectPositions
	err = json.Unmarshal(resp.Data, &op)
	if err != nil {
		return ObjectPositions{}, err
	}

	return op, nil
}

// GeoCode выполняет геокодинг указанной строки. Метод возвращает набор подходящий точек
// с их точными найденными адресами.
func (api *API) GeoCode(addr string) (GeoPoints, error) {
	v := url.Values{}
	v.Add("query", addr)
	resp, err := api.MakeRequest("distance_search", v)
	if err != nil {
		return GeoPoints{}, err
	}

	var gp GeoPoints
	err = json.Unmarshal(resp.Data, &gp)
	if err != nil {
		return GeoPoints{}, err
	}

	return gp, nil
}

// GetOperatorInfo возвращает информацию по оператору объекта трекинга (номеру телефона)
func (api *API) GetOperatorInfo(o Object) (OperatorInfo, error) {
	v, err := o.values()
	if err != nil {
		return OperatorInfo{}, err
	}

	resp, err := api.MakeRequest("get_operator", v)
	if err != nil {
		return OperatorInfo{}, err
	}

	var oi OperatorInfo
	err = json.Unmarshal(resp.Data, &oi)
	if err != nil {
		return OperatorInfo{}, err
	}

	return oi, nil
}
