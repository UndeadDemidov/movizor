package movizor

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DestinationOptions указывает точку следования объекта мониторинга.
// ToDo: заменить на Coordinates
type DestinationOptions struct {
	Text         string
	Lon          float32
	Lat          float32
	ExpectedTime time.Time
}

func (do DestinationOptions) addValuesTo(idx int, v *url.Values) error {
	if v == nil {
		return errors.New("trying to add to nothing")
	}
	if do.Text == "" {
		return errors.New("option Text is not provided")
	}

	v.Add(fmt.Sprintf("destination[%d][text]", idx), do.Text)
	v.Add(fmt.Sprintf("destination[%d][coord]", idx), fmt.Sprintf("%.8f,%.8f", do.Lat, do.Lon))
	if !do.ExpectedTime.IsZero() {
		v.Add(fmt.Sprintf("destination[%d][time]", idx), do.ExpectedTime.Format("02.01.2006 15:04"))
	}

	return nil
}

// SchedulingOptions является расписанием запросов на определение координат объекта.
type SchedulingOptions struct {
	weekdays [7]bool
	FireAt   []time.Time // st Массив времени в расписании. Передается в многомерном массиве,
	// каждый вложенный элемент является временем для срабатывания расписания в формате hh:mm
}

// WeekdayOn добавляет день недели в расписание запросов на определение координат объекта.
func (s *SchedulingOptions) WeekdayOn(day Weekday) {
	s.weekdays[int(day)] = true
}

// WeekdayOff исключает день недели из расписания запросов на определение координат объекта.
func (s *SchedulingOptions) WeekdayOff(day Weekday) {
	s.weekdays[int(day)] = false
}

// IsWeekdayOn возвращает текущее состояние расписания запросов для указанного дня недели.
func (s *SchedulingOptions) IsWeekdayOn(day Weekday) bool {
	return s.weekdays[int(day)]
}

func (s *SchedulingOptions) addValuesTo(v *url.Values) error {
	if len(s.FireAt) <= 0 {
		return errors.New("time to fire scheduling is not set (set FireAt property)")
	}

	if v == nil {
		return errors.New("trying to add to nothing")
	}

	// sw1 string Включить расписание на понедельник
	// sw2 string Включить расписание на вторник
	// sw3 string Включить расписание на среду
	// sw4 string Включить расписание на четверг
	// sw5 string Включить расписание на пятницу
	// sw6 string Включить расписание на субботу
	// sw7 string Включить расписание на воскресенье

	chk := false
	for idx, val := range s.weekdays {
		if val {
			chk = true
			v.Add(fmt.Sprintf("sw%d", idx+1), "1")
		}
	}
	if !chk {
		return errors.New("no single weekday to fire scheduling is set")
	}

	for _, val := range s.FireAt {
		v.Add("st[]", val.Format("15:04"))
	}
	return nil
}

// ObjectOptions предоставляет опции для AddObject (add_object) и EditObject (edit_object)
type ObjectOptions struct {
	Title          string               //title - Название объекта
	Tags           []string             //tags - Список меток через запятую
	DateOff        time.Time            //dateoff - Дата и время автоматического отключения абонента (dd.mm.yyyy hh:mm:ss)
	Tariff         TariffType           //tariff - Id-тарифного плана
	PackageProlong bool                 //package_prolong - Автоматически продлевать пакет (при использовании пакетного тарифа)
	Destinations   []DestinationOptions // destination[] - массив конечных точек маршрута.
	Schedules      *SchedulingOptions   // рассписание для ручного обновления координат
	Metadata       map[string]string    // metadata Массив с дополнительной информацией по объекту для отображения в событиях и отчетах.
	CallToDriver   bool                 // autoinform integer Включить услугу автоинформатора.
}

func (o *ObjectOptions) addValuesTo(v *url.Values) error {
	if v == nil {
		return errors.New("trying to add to nothing")
	}
	if o.Title != "" {
		v.Add("title", o.Title)
	}
	if !o.DateOff.IsZero() {
		v.Add("dateoff", o.DateOff.Format("02.01.2006 15:04:05"))
	}
	if o.Tariff != "" {
		v.Add("tariff", string(o.Tariff))
	}
	if o.PackageProlong {
		v.Add("package_prolong", "1")
	}
	if len(o.Tags) > 0 {
		v.Add("tags", strings.Join(o.Tags, ","))
	}
	// Каждый элемент обозначает одну запись метаинформации для объекта. Имена элементов и значения произвольные:
	// metadata[Исполнитель] = Петров;
	// metadata[Склад] = Восточный;
	for key, val := range o.Metadata {
		v.Add(fmt.Sprintf("metadata[%s]", key), val)
	}
	if o.CallToDriver {
		v.Add("autoinform", "1")
	}

	if o.Schedules != nil {
		err := o.Schedules.addValuesTo(v)
		if err != nil {
			return err
		}
	}

	for key, val := range o.Destinations {
		err := val.addValuesTo(key, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// RequestPositionsOptions - представляет собой указание выборки списка координат методом
// GetObjectPositions.
type RequestPositionsOptions struct {
	RequestLimit uint64    // req_limit - Разрешить делать не более X запросов в сутки на определение координат всех объектов
	Offset       uint64    // offset - Смещение количества получаемых координат
	TimeFrom     time.Time // date_start - Unix Timestamp. Фильтрация вывода, начиная с этой даты
	TimeTo       time.Time // date_end - Unix Timestamp. Фильтрация вывода, до этой даты
}

func (rpo *RequestPositionsOptions) addValuesTo(v *url.Values) error {
	if v == nil {
		return errors.New("trying to add to nothing")
	}
	if rpo.RequestLimit != 0 {
		v.Add("req_limit", strconv.FormatUint(rpo.RequestLimit, 10))
	}
	if rpo.Offset != 0 {
		v.Add("offset", strconv.FormatUint(rpo.Offset, 10))
	}
	if !rpo.TimeFrom.IsZero() {
		v.Add("date_start", strconv.FormatInt(rpo.TimeFrom.Unix(), 10))
	}
	if !rpo.TimeTo.IsZero() {
		v.Add("date_end", strconv.FormatInt(rpo.TimeTo.Unix(), 10))
	}

	return nil
}

// ObjectEventsOptions предоставляет опции для получения списка событий через
// метод GetEvents.
type ObjectEventsOptions struct {
	RequestLimit uint64
	AfterEventID uint64
}

func (eo ObjectEventsOptions) values() url.Values {
	v := url.Values{}
	if eo.RequestLimit != 0 {
		v.Add("req_limit", strconv.FormatUint(eo.RequestLimit, 10))
	}
	if eo.AfterEventID != 0 {
		v.Add("afterid", strconv.FormatUint(eo.AfterEventID, 10))
	}
	return v
}

// SubscribeEventOptions предоставляет опции подписки на нотификацию по событиям.
// Если установлен признак AllObjects, то список Objects игнорируется.
type SubscribeEventOptions struct {
	AllObjects bool
	Objects    []Object
	Event      EventType
	notifyTo   notificationType
	smsPhone   Object
	email      string
}

// NewSubscribeEventOptions возвращает экземпляр SubscribeEventOptions для указанного
// объекта и типа события.
func NewSubscribeEventOptions(o Object, e EventType) SubscribeEventOptions {
	return SubscribeEventOptions{
		Objects: []Object{o},
		Event:   e,
	}
}

// IsSubscriberEqualTo проверяет, что у опций подписки на конкретное событие один и тот же подписчик
func (se *SubscribeEventOptions) IsSubscriberEqualTo(event SubscribedEvent) bool {
	return se.Event == event.Event && ((se.notifyTo == telegramNotification && event.IsTelegram) ||
		(se.notifyTo == emailNotification && se.email == event.EMail) ||
		(se.notifyTo == smsNotification && se.smsPhone == event.Phone))
}

// SetSMSNotification устанавливает нотификацию на указанный телефон по СМС. Работает только та нотификация,
// которая была установлена последней в данной подписке. Это особенности API Movizor.
func (se *SubscribeEventOptions) SetSMSNotification(phone Object) error {
	if phone.String() == "" {
		return fmt.Errorf("invalid phone number %s", string(phone))
	}

	se.notifyTo = smsNotification
	se.smsPhone = phone
	return nil
}

// SetEMailNotification устанавливает нотификацию на указанный почтовый адрес. Работает только та нотификация,
// которая была установлена последней в данной подписке. Это особенности API Movizor.
func (se *SubscribeEventOptions) SetEMailNotification(mail string) error {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(mail) {
		return fmt.Errorf("%s is not valid email address", mail)
	}
	se.notifyTo = emailNotification
	se.email = mail
	return nil
}

// SetTelegramNotification устанавливает нотификацию на Телеграм указанные в профиле держателя аккаута.
// Работает только та нотификация, которая была установлена последней в данной подписке. Это особенности API Movizor.
func (se *SubscribeEventOptions) SetTelegramNotification() {
	se.notifyTo = telegramNotification
}

func (se SubscribeEventOptions) values() (url.Values, error) {
	if !se.AllObjects && len(se.Objects) == 0 {
		return url.Values{}, errors.New("no single phone is set to subscribe for event")
	}
	if string(se.Event) == "" {
		return url.Values{}, errors.New("event to subscribe is not set")
	}
	if string(se.notifyTo) == "" {
		return url.Values{}, errors.New("notification type is not set")
	}

	v := url.Values{}

	if se.AllObjects {
		v.Add("phones_all", "1")
	} else {
		for _, val := range se.Objects {
			v.Add("phones[]", val.String())

		}
	}

	v.Add("events", string(se.Event))
	v.Add("notify_type", string(se.notifyTo))
	switch se.notifyTo {
	case smsNotification:
		v.Add("notify_value", se.smsPhone.String())
	case emailNotification:
		v.Add("notify_value", se.email)
	}

	return v, nil
}
