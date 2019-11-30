package movizor

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// GetEvents получает список событий, с возможностью определить с какого id события выводить данные.
func (api *API) GetEvents(o ObjectEventsOptions) (ObjectEvents, error) {
	resp, err := api.MakeRequest("events", o.values())
	if err != nil {
		return ObjectEvents{}, err
	}

	var oe ObjectEvents
	err = json.Unmarshal(resp.Data, &oe)
	if err != nil {
		return ObjectEvents{}, err
	}

	return oe, nil
}

// DeleteEventsSubscription удаляет подписку по ее id. Для получения id используйте GetEventSubscriptions.
func (api *API) DeleteEventsSubscription(id int64) (APIResponse, error) {
	v := url.Values{}
	v.Add("id", strconv.FormatInt(id, 10))
	resp, err := api.MakeRequest("events_subscribe_delete", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetEventSubscriptions получает список подписок активных на текущий момент.
func (api *API) GetEventSubscriptions() (SubscribedEvents, error) {
	resp, err := api.MakeRequest("events_subscribe_list", nil)
	if err != nil {
		return SubscribedEvents{}, err
	}

	var se SubscribedEvents
	err = json.Unmarshal(resp.Data, &se)
	if err != nil {
		return SubscribedEvents{}, err
	}

	return se, nil
}

// SubscribeEvent выполняет подписку на указанное тип события для всех объектов (телефонов) или по списку.
func (api *API) SubscribeEvent(o SubscribeEventOptions) (APIResponse, error) {
	v, err := o.values()
	if err != nil {
		return APIResponse{}, err
	}

	resp, err := api.MakeRequest("events_subscribe_add", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// AddSubscription добавляет подписку на событие по телефону в существующую подписку подписчика,
// если таковой будет найден. Если такой подписчик не будет найден, то создастся новая подписка.
func (api *API) AddEventSubscription(o SubscribeEventOptions) error {
	events, err := api.GetEventSubscriptions()
	if err != nil {
		return err
	}

	opt := o
	for _, e := range events {
		if opt.IsSubscriberEqualTo(e) {
			if e.IsAllObjectsSubscribed {
				return nil
			}

			opt, err = e.MakeOptions()
			if err != nil {
				return err
			}
			opt.Objects = append(opt.Objects, o.Objects...)
			break
		}
	}
	_, err = api.SubscribeEvent(opt)
	if err != nil {
		return err
	}

	return nil
}

// ClearAllEventSubscriptions удаляет все подписки в аккаунте.
func (api *API) ClearAllEventSubscriptions() error {
	events, err := api.GetEventSubscriptions()
	if err != nil {
		return err
	}
	for _, e := range events {
		_, err = api.DeleteEventsSubscription(e.SubscriptionID)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnsubscribeObject производит отписку от всех событий для определенного телефона.
// Если существует подписка на все телефоны на какое-то событие, то она не будет затронута.
// Т.е. удаляются только подписки с явным указанием номера телефона.
func (api *API) UnsubscribeObject(o Object) error {
	return api.ClearObjectEventSubscriptions(o, nil)
}

// ClearObjectEventSubscriptions производит отписку от конкретного события для определенного телефона.
// Если существует подписка на все телефоны на какое-то событие, то она не будет затронута.
// Т.е. удаляются только подписки с явным указанием номера телефона.
func (api *API) ClearObjectEventSubscriptions(o Object, eType *EventType) error {
	events, err := api.GetEventSubscriptions()
	if err != nil {
		return err
	}

	isUnused := func(obj Object, ev *EventType) bool {
		return obj == o && (eType == nil || *ev == *eType)
	}

	for _, e := range events {
		if e.IsAllObjectsSubscribed {
			continue
		}

		err := api.removeObjectSubscriptions(e, isUnused)
		if err != nil {
			return err
		}
	}
	return nil
}

// ClearUnusedSubscriptions удаляет все не используемые подписки. То есть,
// если конретного объекта нет в списке трекинга (в любом статусе), то подписка
// на события по этому объекуту удаляется.
// Удаление касается всех подписок для конкретных телефонов.
// Общие подписки для всех телефонов не затраниваются.
func (api *API) ClearUnusedSubscriptions() error {
	events, err := api.GetEventSubscriptions()
	if err != nil {
		return err
	}

	trackList, err := api.GetObjects()
	if err != nil {
		return err
	}

	isUnused := func(o Object, e *EventType) bool {
		return !trackList.IsObjectIn(o)
	}

	for _, e := range events {
		if e.IsAllObjectsSubscribed {
			continue
		}

		err := api.removeObjectSubscriptions(e, isUnused)
		if err != nil {
			return err
		}
	}
	return nil
}

type shouldRemoveSubscription func(Object, *EventType) bool

func (api *API) removeObjectSubscriptions(e SubscribedEvent, f shouldRemoveSubscription) error {
	for i, phone := range e.ObjectsSubscribed {
		if f(phone, &e.Event) {
			_, err := api.DeleteEventsSubscription(e.SubscriptionID)
			if err != nil {
				return err
			}

			if len(e.ObjectsSubscribed) > 1 {
				seo, err := e.MakeOptions()
				if err != nil {
					return err
				}

				seo.Objects = append(seo.Objects[:i], seo.Objects[i+1:]...)
				seo.Objects = deDupObjectSlice(seo.Objects)
				_, err = api.SubscribeEvent(seo)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func deDupObjectSlice(src []Object) []Object {
	m := make(map[Object]bool)
	var result []Object
	for _, v := range src {
		if _, f := m[v]; !f {
			m[v] = true
			result = append(result, v)
		}
	}
	return result
}
