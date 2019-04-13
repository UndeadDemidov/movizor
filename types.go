package movizor

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// APIResponse представляет собой ответ от сервиса с описанием типа
// сообщения и сегмента с данными.
type APIResponse struct {
	Result      string          `json:"result"`                  // "success" or "error" expected
	ResultCode  string          `json:"code"`                    // "OK" expected
	MessageType string          `json:"message"`                 // "Balance info", ...
	Data        json.RawMessage `json:"data,omitempty"`          // optional Payload of response
	ErrorCode   string          `json:"error_code,omitempty"`    // Код ошибка
	ErrorText   string          `json:"error_text,omitempty"`    // Текст ошибки
	ErrorTextRU string          `json:"error_text_ru,omitempty"` // optional Текст ошибки на русском
}

// Tariff представляет собой структуру тариф для одного из сервисов.
type Tariff struct {
	AbonentPayment float64 `json:"abon"`    // Абоненская плата
	RequestCost    float64 `json:"request"` // Стоимость запроса
	TariffTitle    string  `json:"title"`   // Название тарифа
}

func (t *Tariff) UnmarshalJSON(data []byte) (err error) {
	type Alias Tariff
	aux := &struct {
		AbonentPayment json.Number `json:"abon"`    // Абоненская плата
		RequestCost    json.Number `json:"request"` // Стоимость запроса
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if t.AbonentPayment, err = aux.AbonentPayment.Float64(); err != nil {
		return err
	}

	if t.RequestCost, err = aux.RequestCost.Float64(); err != nil {
		return err
	}

	return nil
}

// Balance содержит текущий баланс по договору и список подключенных тарифов
// по мобильным операторам.
type Balance struct {
	Balance         float64 `json:"balance,string"` // Текущий остаток средств на балансе
	Credit          float64 `json:"credit,string"`  // Сумма кредитных средств на балансе
	ContractType    string  `json:"type"`           // Тип договора
	OperatorTariffs map[Operator]map[TariffType]Tariff
	ServiceTariffs  map[Service][]Tariff
}

func (b *Balance) UnmarshalJSON(data []byte) (err error) {
	type Alias Balance
	aux := &struct {
		TariffPlans map[string]json.RawMessage `json:"tariff"` // Список операторов с их тарифами и доп тарифы
		*Alias
	}{
		Alias: (*Alias)(b),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	ops := make(map[Operator]map[TariffType]Tariff)
	srv := make(map[Service][]Tariff)

	for key, data := range aux.TariffPlans {
		var sliceTrf []Tariff
		var trf map[TariffType]Tariff

		if err = json.Unmarshal(data, &sliceTrf); err == nil {
			srv[Service(key)] = sliceTrf
			continue
		}

		if err = json.Unmarshal(data, &trf); err == nil {
			ops[Operator(key)] = trf
			continue
		}

		return err
	}
	b.OperatorTariffs = ops
	b.ServiceTariffs = srv

	return nil
}

// CurrentCoordinates представляют собой текущие гео-координаты объекта.
// Допускается null значения.
type CurrentCoordinates struct {
	CurrentLon *Coordinate `json:"current_lon"` // Широта последнего местоположения
	CurrentLat *Coordinate `json:"current_lat"` // Долгота последнего местоположения
}

// ObjectInfo содержит почти полную информацию по объекту, включая опции,
// с которыми добавлялся объект.
type ObjectInfo struct {
	Phone         Object      `json:"phone"`                // Номер абонента
	Status        Status      `json:"status"`               // status type
	Confirmed     bool        `json:"confirmed"`            // Получено подтверждение от абонента
	Title         string      `json:"title"`                // Имя абонента (название объекта)
	Tariff        TariffType  `json:"tariff"`               // Текущий тарифный план
	TariffNew     *TariffType `json:"tariff_new,omitempty"` // Новый тарифный план со следующего дня
	LastTimestamp Time        `json:"last_timestamp"`       // Время последнего запроса на определение местоположения
	AtRequest     bool        `json:"at_request,omitempty"` // Производится определение местоположения в данный момент
	CurrentCoordinates
	CoordinatesAttributes
	OnParking    *bool             `json:"on_parking,omitempty"`   // Находится ли объект на парковке
	Destination  []Destination     `json:"destination,omitempty"`  // Список точек назначения, которые должен посетить Водитель.
	OfflineTime  Time              `json:"offline_time,omitempty"` // Время последнего известного местоположения
	PosError     bool              `json:"pos_error,omitempty"`    // Последнее местоположение не удалось определить
	TimestampOff Time              `json:"timestamp_off"`          // Время автоматического отключения от мониторинга
	TimestampAdd Time              `json:"timestamp_add"`          // Время добавления объекта в Мовизор
	Metadata     map[string]string `json:"metadata,omitempty"`     // Метаинформация объекта, массив
}

func (oi *ObjectInfo) UnmarshalJSON(data []byte) (err error) {
	type Alias ObjectInfo
	aux := &struct {
		Metadata json.RawMessage `json:"metadata,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(oi),
	}

	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var probe []interface{}
	if err = json.Unmarshal(aux.Metadata, &probe); err == nil {
		return nil
	}

	return json.Unmarshal(aux.Metadata, &oi.Metadata)
}

// Coordinates представляет собой структуру гео-координат.
type Coordinates struct {
	Lat Coordinate `json:"lat"` // Долгота
	Lon Coordinate `json:"lon"` // Широта
}

// Destination представляю собой структуру описания точки назначения,
// в которую следует объект.
type Destination struct {
	Text string `json:"text"`
	Coordinates
	Time   string    `json:"time"`
	Status ETAStatus `json:"status"`
}

// ObjectStatus представляет собой текущий статус объекта трекинга.
type ObjectStatus struct {
	Phone  Object `json:"phone"`  // Номер телефона абонента
	Status Status `json:"status"` // Статус добавления для отслеживания
}

// ObjectsWithStatus является списком объектов с их статусами.
type ObjectsWithStatus []ObjectStatus

func (os ObjectsWithStatus) Len() int           { return len(os) }
func (os ObjectsWithStatus) Swap(i, j int)      { os[i], os[j] = os[j], os[i] }
func (os ObjectsWithStatus) Less(i, j int) bool { return os[i].Phone < os[j].Phone }

// IsObjectIn проверяет, есть ли соответствующий объект в списке.
func (os ObjectsWithStatus) IsObjectIn(o Object) bool {
	for _, os := range os {
		if os.Phone == o {
			return true
		}
	}
	return false
}

// CoordinatesAttributes преставляет собой аттрибуты гео-координат.
// ETA, статус ETA, описание гео-координат (обратный гео-кодинг)
type CoordinatesAttributes struct {
	Distance *Int `json:"distance"`                         // Остаток в км до конечной точки
	ETA      *Int `json:"distance_forecast_time,omitempty"` // Прогноз оставшегося времени до конечной точки
	// Прогноз строится в зависимости от наличия информации о конечном пункте назначения и времени прибытия.
	// Если этой информации нет, значения элементов будут пустыми.
	ETAStatus *ETAStatus `json:"distance_forecast_status,omitempty"` // Прогноз успеваемости до конечной точки.
	// Прогноз строится в зависимости от наличия информации о конечном пункте назначения и времени прибытия.
	// Если этой информации нет, значения элементов будут пустыми.
	Place string `json:"place"` // Населенный пункт местоположения.
}

//func (ca *CoordinatesAttributes) UnmarshalJSON(data []byte) (err error) {
//	type Alias CoordinatesAttributes
//	aux := &struct {
//		Distance json.Number `json:"distance"`
//		*Alias
//	}{
//		Alias: (*Alias)(ca),
//	}
//
//	if err = json.Unmarshal(data, &aux); err != nil {
//		return err
//	}
//
//	if ca.Distance, err = aux.Distance.Int64(); err != nil {
//		return err
//	}
//
//	return nil
//}

// Positions является списком местоположений.
type Positions []Position

// Position содержит информацию о последнем зафиксированном в системе местоположении.
type Position struct {
	Coordinates
	Timestamp        Time `json:"timestamp"`                   // Время получения координат для этой точки
	TimestampRequest Time `json:"timestamp_request,omitempty"` // Время создания запроса на получение координат
	Deviation        *Int `json:"radius,omitempty"`            // Радиус погрешности (м)
	CoordinatesAttributes
}

// ObjectPositions является списком объектов с гео-координатами, последним
// временем обновления координат, текущим местонахождением и ETA.
type ObjectPositions []ObjectPosition

// ObjectPosition представляет собой гео-координаты, последнее временя
// обновления координат, текущее местонахождение и ETA объекта.
type ObjectPosition struct {
	Phone Object `json:"phone"` // Номер телефона абонента
	Position
}

// PositionRequest хранит ID запроса на опреления гео-координат.
type PositionRequest struct {
	RequestID int64 `json:"request_id"`
}

func (pr PositionRequest) values() url.Values {
	return url.Values{"id": {strconv.FormatInt(pr.RequestID, 10)}}
}

// GeoPoints содержит слайс точек геокодирования.
type GeoPoints []GeoPoint

// GeoPoint указывает на точку геокодирования.
type GeoPoint struct {
	Coordinates
	Description string `json:"title"`
}

// OperatorInfo содержит информацию о сотовом операторе.
type OperatorInfo struct {
	Operator Operator `json:"operator"`         // Внутренний идентификатор оператора
	Title    string   `json:"title"`            // Название оператора
	Region   string   `json:"region,omitempty"` // Домашний регион абонента
}

// ObjectEvents является списком событий по объектам.
type ObjectEvents []ObjectEvent

// ObjectEvent содержит информацию о событиях, которые происходили с объектом.
// Такие как: подтверждение трекинга, отклонение трекига, отклонения от маршрута следования и т.д.
type ObjectEvent struct {
	EventID   int64     `json:"id"`        // Идентификатор события (возрастающий номер события)
	Timestamp Time      `json:"timestamp"` // Время возникновения события
	Phone     Object    `json:"phone"`     // Номер телефона абонента, по которому произошло событие
	Event     EventType `json:"type"`      // Тип события
}

func (oe *ObjectEvent) UnmarshalJSON(data []byte) (err error) {
	type Alias ObjectEvent
	aux := &struct {
		EventID json.Number `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(oe),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if oe.EventID, err = aux.EventID.Int64(); err != nil {
		return err
	}

	return nil
}

// SubscribedEvents является списком подписок на события.
type SubscribedEvents []SubscribedEvent

// SubscribedEvent содержит информацию о подписке на одно из событий.
type SubscribedEvent struct {
	SubscriptionID         int64     `json:"id"`         // Идентификатор события (возрастающий номер события)
	IsAllObjectsSubscribed bool      `json:"phones_all"` // Уведомление о событии для всех объектов (в том числе добавляемых в будущем)
	ObjectsSubscribed      []Object  `json:"phones"`     // Список телефонов (объектов)
	Timestamp              Time      `json:"timestamp"`  // Время возникновения события
	Event                  EventType `json:"type"`       // Тип события, на которые зарегистрирована подписка
	Phone                  Object    `json:"phone"`      // Номер телефона абонента, по которому отправляются уведомления
	EMail                  string    `json:"email"`      // Email, по которому отправляются уведомления
	IsTelegram             bool      `json:"telegram"`   // Уведомления отправляются на аккаунт telegram указанный в настройках аккаунта
}

func (se *SubscribedEvent) UnmarshalJSON(data []byte) (err error) {
	type Alias SubscribedEvent
	aux := &struct {
		SubscriptionID         json.Number `json:"id"`
		IsAllObjectsSubscribed int         `json:"phones_all"`
		IsTelegram             int         `json:"telegram"`
		*Alias
	}{
		Alias: (*Alias)(se),
	}
	if err = json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if se.SubscriptionID, err = aux.SubscriptionID.Int64(); err != nil {
		return err
	}

	if aux.IsAllObjectsSubscribed == 1 {
		se.IsAllObjectsSubscribed = true
	}

	if aux.IsTelegram == 1 {
		se.IsTelegram = true
	}

	return nil
}

// MakeOptions создает опции на создание подписки на события на основе
// существуещей подписки для последующего добавления.
// Используется для редактирования подписки путем удаления старой и
// добавления новой.
func (se SubscribedEvent) MakeOptions() (seo SubscribeEventOptions, err error) {
	seo = SubscribeEventOptions{}
	seo.Event = se.Event
	if se.IsAllObjectsSubscribed {
		seo.AllObjects = true
	} else {
		seo.Objects = se.ObjectsSubscribed
	}
	switch {
	case se.Phone.String() != "":
		err = seo.SetSMSNotification(se.Phone)
	case se.EMail != "":
		err = seo.SetEMailNotification(se.EMail)
	default:
		seo.SetTelegramNotification()
	}
	return
}
