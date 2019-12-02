package movizor

const (
	// DefaultAPIMovizorEndpoint is default API Movizor endpoint
	DefaultAPIMovizorEndpoint = "https://movizor.ru/api"
	// APIMovizorEndpointSuffix is API Movizor endpoint suffix to pass method and parameters
	// /%s/%s = /project/action
	APIMovizorEndpointSuffix = "/%s/%s"
)

// Operator представляет собой Операторов мобильной связи.
type Operator string

const (
	// МТС
	OperatorMTS Operator = "mts"
	// Мегафон
	OperatorMegafon Operator = "megafon"
	// Билайн
	OperatorBeeline Operator = "beeline"
	// Теле2
	OperatorTele2 Operator = "tele2"
)

// Service представляет собой тип доп сервиса МоВизора.
type Service string

const (
	EventSmsService   Service = "eventsms"   // Подписка на SMS по событиям.
	AutoInformService Service = "autoinform" // Голосовое информирование на телефон подключаемого абонента.
)

// TariffType представляет собой тип тарифа.
type TariffType string

const (
	TariffManual   TariffType = "0"   // Вручную
	TariffOnline   TariffType = "1"   // Онлайн
	TariffOneMonth TariffType = "3"   // Пакет на 1 месяц
	TariffEvery15  TariffType = "15"  // Каждые 15 мин
	TariffEvery30  TariffType = "30"  // Каждые 30 мин
	TariffEvery60  TariffType = "60"  // Каждые 60 мин
	TariffEvery180 TariffType = "180" // Каждые 3 часа
)

// Status представляет собой возможный статус состояния объекта в системе МоВизор.
type Status string

const (
	StatusNew                Status = "new"          // Новый
	StatusWaitOk             Status = "wait"         // Ожидание подтверждения
	StatusOk                 Status = "ok"           // Подтвержден и доступен для мониторинга
	StatusWaitOff            Status = "off_new"      // В процессе отключения от мониторинга
	StatusOff                Status = "off"          // Отключен от мониторинга
	StatusRejected           Status = "rejected"     // Абонент отказался от мониторинга
	StatusNotConfirmed       Status = "no_confirmed" // Абонент не подтвердил подключение
	StatusError              Status = "error"        // Оператор не поддерживается
	StatusOperatorNotAllowed Status = "wrong"        // Оператор не поддерживается
	StatusTrackCollision     Status = "error_exists" // Номер отслеживается другой компанией у оператора
	StatusTrackingIsOff      Status = "limited"      // У абонента ограничение на подключение к услуге
)

// ETAStatus представляет собой возможный статус ETA объекта.
type ETAStatus string

const (
	NewETAStatus      ETAStatus = "new"      // успевает
	OkETAStatus       ETAStatus = "ok"       // успевает
	LateETAStatus     ETAStatus = "late"     // опаздывает
	FinishedETAStatus ETAStatus = "finished" // прибыл
)

// EventType представляет собой возможный тип события, которые регистрируются в сервисе МоВизора.
// Не по всем типам событий можно производить подписку для нотификации.
type EventType string

const (
	AddEvent                    EventType = "add"             // - добавлен объект
	AutoOffEvent                EventType = "auto_off"        // - автоматическое отключение
	OffEvent                    EventType = "off"             // - отключение
	ConfirmEvent                EventType = "confirm"         // - объект подтвердил подключение
	RejectEvent                 EventType = "reject"          // - объект отказался от подключения
	RequestOkEvent              EventType = "request"         // - запрос: успешно
	RequestErrorEvent           EventType = "request_error"   // - запрос: ошибка
	RequestObjectOfflineEvent   EventType = "request_offline" // - запрос: телефон недоступен
	RequestObjectInRoamingEvent EventType = "request_roaming" // - запрос: телефон в роуминге
	ReactivateEvent             EventType = "reactivate"      // - повторное подключение
	ChangeTariffEvent           EventType = "tariff_auto"     // - смена тарифного плана
	InTimeEvent                 EventType = "pos_ok"          // - объект начал успевать
	LateEvent                   EventType = "pos_late"        // - объект начал опаздывать
	FinishedEvent               EventType = "pos_finished"    // - объект прибывает
	CallToDriverEvent           EventType = "autoinform"      // - автоинформатор
	NoConfirmationEvent         EventType = "no_confirm"      // - объект не подтвердил подключение
	ObjectLimitedEvent          EventType = "limit"           // - у объекта стоит ограничение
	OnRouteEvent                EventType = "onroute"         // - встал на маршрут
	ReturnRouteEvent            EventType = "returnroute"     // - вернулся на маршрут
	LeftRouteEvent              EventType = "leftroute"       // - отклонился от маршрута
	NotRouteEvent               EventType = "notroute"        // - не на маршруте
	OnParkingEvent              EventType = "onparking"       // - встал на парковку
	OffParkingEvent             EventType = "offparking"      // - начал движение
	MStopEvent                  EventType = "mstop"           // - приложение остановлено
	MStartEvent                 EventType = "mstart"          // - приложение запущено
)

type notificationType string

const (
	smsNotification      notificationType = "sms"
	emailNotification    notificationType = "email"
	telegramNotification notificationType = "telegram"
)

type Weekday int

const (
	Monday Weekday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)
