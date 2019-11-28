package movizor

import (
	"encoding/json"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestAPI_GetEvents(t *testing.T) {
	type args struct {
		o ObjectEventsOptions
	}
	tests := []struct {
		name         string
		args         args
		filename     string
		filenameWant string
		wantErr      bool
	}{
		{
			name: "events_good",
			args: args{
				o: ObjectEventsOptions{
					RequestLimit: 0,
					AfterEventID: 0,
				},
			},
			filename:     "events_resp1.json",
			filenameWant: "events.json",
			wantErr:      false,
		},
		{
			name: "events_bad_json",
			args: args{
				o: ObjectEventsOptions{
					RequestLimit: 0,
					AfterEventID: 0,
				},
			},
			filename:     "bad.json",
			filenameWant: "events.json",
			wantErr:      true,
		},
		{
			name: "events_bad1",
			args: args{
				o: ObjectEventsOptions{
					RequestLimit: 0,
					AfterEventID: 0,
				},
			},
			filename:     "error_response.json",
			filenameWant: "events.json",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := ioutil.ReadFile(filepath.Join(dataPath, tt.filename))
			if err != nil {
				t.Errorf("err: %s", err)
			}

			dWant, err := ioutil.ReadFile(filepath.Join(dataPath, tt.filenameWant))
			if err != nil {
				t.Errorf("err: %s", err)
			}
			want := ObjectEvents{}
			if err := json.Unmarshal(dWant, &want); err != nil {
				t.Errorf("ObjectEvents.UnmarshalJSON() error = %v", err)
			}

			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/events", responder)

			got, err := api.GetEvents(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.GetEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (!tt.wantErr && !reflect.DeepEqual(got, want)) ||
				(tt.wantErr && !reflect.DeepEqual(got, ObjectEvents{})) {
				t.Errorf("API.GetEvents() = %v, want %v", got, want)
			}
		})
	}
}

func TestAPI_DeleteEventsSubscription(t *testing.T) {
	type args struct {
		id int64
	}
	tests := []struct {
		name     string
		args     args
		filename string
		want     APIResponse
		wantErr  bool
	}{
		{
			name: "events_subscribe_delete_good",
			args: args{
				id: 123,
			},
			filename: "success_response.json",
			want: APIResponse{
				Result:     "success",
				ResultCode: "OK",
				Message:    "Some message",
			},
			wantErr: false,
		},
		{
			name: "events_subscribe_delete_bad_json",
			args: args{
				id: 123,
			},
			filename: "bad.json",
			want:     APIResponse{},
			wantErr:  true,
		},
		{
			name: "events_subscribe_delete_bad1",
			args: args{
				id: 123,
			},
			filename: "error_response.json",
			want: APIResponse{
				Result:     "error",
				ResultCode: "ACCESS_DENIED",
				ErrorText:  "Auth rate limit exceeded",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := ioutil.ReadFile(filepath.Join(dataPath, tt.filename))
			if err != nil {
				t.Errorf("err: %s", err)
			}

			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/events_subscribe_delete", responder)

			got, err := api.DeleteEventsSubscription(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.DeleteEventsSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Result, tt.want.Result) {
				t.Errorf("API.DeleteEventsSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_GetEventSubscriptions(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		filenameWant string
		wantErr      bool
	}{
		{
			name:         "events_subscribe_list_good",
			filename:     "events_subscribe_list_resp1.json",
			filenameWant: "events_subscribe_list.json",
			wantErr:      false,
		},
		{
			name:         "events_subscribe_list_bad_json",
			filename:     "bad.json",
			filenameWant: "events_subscribe_list.json",
			wantErr:      true,
		},
		{
			name:         "events_subscribe_list_bad1",
			filename:     "error_response.json",
			filenameWant: "events_subscribe_list.json",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := ioutil.ReadFile(filepath.Join(dataPath, tt.filename))
			if err != nil {
				t.Errorf("err: %s", err)
			}

			dWant, err := ioutil.ReadFile(filepath.Join(dataPath, tt.filenameWant))
			if err != nil {
				t.Errorf("err: %s", err)
			}
			want := SubscribedEvents{}
			if err := json.Unmarshal(dWant, &want); err != nil {
				t.Errorf("SubscribedEvents.UnmarshalJSON() error = %v", err)
			}

			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/events_subscribe_list", responder)

			got, err := api.GetEventSubscriptions()
			if (err != nil) != tt.wantErr {
				t.Errorf("API.GetEventSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (!tt.wantErr && !reflect.DeepEqual(got, want)) ||
				(tt.wantErr && !reflect.DeepEqual(got, SubscribedEvents{})) {
				t.Errorf("API.GetEventSubscriptions() = %v, want %v", got, want)
			}
		})
	}
}

func TestAPI_SubscribeEvent(t *testing.T) {
	type args struct {
		o SubscribeEventOptions
	}
	tests := []struct {
		name     string
		args     args
		filename string
		want     APIResponse
		wantErr  bool
	}{
		{
			name: "events_subscribe_add_good",
			args: args{
				o: SubscribeEventOptions{
					AllObjects: true,
					Objects:    nil,
					Event:      NoConfirmationEvent,
					notifyTo:   emailNotification,
					smsPhone:   "",
					email:      "n.demidov@some.some",
				},
			},
			filename: "success_response.json",
			want: APIResponse{
				Result:     "success",
				ResultCode: "OK",
				Message:    "Some message",
			},
			wantErr: false,
		},
		{
			name: "events_subscribe_add_bad_no_phones",
			args: args{
				o: SubscribeEventOptions{
					AllObjects: false,
					Objects:    nil,
					Event:      NoConfirmationEvent,
					notifyTo:   emailNotification,
					smsPhone:   "",
					email:      "n.demidov@some.some",
				},
			},
			filename: "success_response.json",
			want:     APIResponse{},
			wantErr:  true,
		},
		{
			name: "events_subscribe_add_bad_no_event",
			args: args{
				o: SubscribeEventOptions{
					AllObjects: true,
					Objects:    nil,
					Event:      "",
					notifyTo:   emailNotification,
					smsPhone:   "",
					email:      "n.demidov@some.some",
				},
			},
			filename: "success_response.json",
			want:     APIResponse{},
			wantErr:  true,
		},
		{
			name: "events_subscribe_add_bad_no_notification_type",
			args: args{
				o: SubscribeEventOptions{
					AllObjects: true,
					Objects:    nil,
					Event:      NoConfirmationEvent,
					notifyTo:   "",
					smsPhone:   "",
					email:      "n.demidov@some.some",
				},
			},
			filename: "success_response.json",
			want:     APIResponse{},
			wantErr:  true,
		},
		{
			name: "events_subscribe_add_bad_json",
			args: args{
				o: SubscribeEventOptions{
					AllObjects: true,
					Objects:    nil,
					Event:      NoConfirmationEvent,
					notifyTo:   "emailNotification",
					smsPhone:   "",
					email:      "n.demidov@some.some",
				},
			},
			filename: "bad.json",
			want:     APIResponse{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := ioutil.ReadFile(filepath.Join(dataPath, tt.filename))
			if err != nil {
				t.Errorf("err: %s", err)
			}

			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/events_subscribe_add", responder)
			got, err := api.SubscribeEvent(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.SubscribeEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API.SubscribeEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_AddEventSubscription(t *testing.T) {
	type args struct {
		o SubscribeEventOptions
	}
	tests := []struct {
		name                 string
		subscriptionFileName string
		respondFileName      string
		args                 args
		wantErr              bool
	}{
		{
			name:                 "add_good",
			subscriptionFileName: "events_subscribe_list_resp1.json",
			respondFileName:      "success_response.json",
			args: args{
				o: SubscribeEventOptions{
					AllObjects: false,
					Objects: []Object{
						"79001294569",
					},
					Event:    OnParkingEvent,
					notifyTo: emailNotification,
					smsPhone: "",
					email:    "ndemidov@some.some",
				},
			},
			wantErr: false,
		},
		{
			name:                 "add_good",
			subscriptionFileName: "events_subscribe_list_resp1.json",
			respondFileName:      "success_response.json",
			args: args{
				o: SubscribeEventOptions{
					AllObjects: false,
					Objects: []Object{
						"79001294569",
					},
					Event:    ConfirmEvent,
					notifyTo: emailNotification,
					smsPhone: "",
					email:    "ndemidov@some.some",
				},
			},
			wantErr: false,
		},
		{
			name:                 "add_good",
			subscriptionFileName: "events_subscribe_list_resp1.json",
			respondFileName:      "success_response.json",
			args: args{
				o: SubscribeEventOptions{
					AllObjects: false,
					Objects: []Object{
						"79001294569",
					},
					Event:    ConfirmEvent,
					notifyTo: smsNotification,
					smsPhone: "79123456789",
					email:    "",
				},
			},
			wantErr: false,
		},
		{
			name:                 "add_good",
			subscriptionFileName: "events_subscribe_list_resp1.json",
			respondFileName:      "success_response.json",
			args: args{
				o: SubscribeEventOptions{
					AllObjects: false,
					Objects: []Object{
						"79001294569",
					},
					Event:    RejectEvent,
					notifyTo: emailNotification,
					smsPhone: "",
					email:    "ndemidov@some.some",
				},
			},
			wantErr: false,
		},
		{
			name:                 "add_bad",
			subscriptionFileName: "error_response.json",
			respondFileName:      "success_response.json",
			args: args{
				o: SubscribeEventOptions{
					AllObjects: false,
					Objects: []Object{
						"79001294569",
					},
					Event:    RejectEvent,
					notifyTo: emailNotification,
					smsPhone: "",
					email:    "ndemidov@some.some",
				},
			},
			wantErr: true,
		},
		{
			name:                 "add_bad",
			subscriptionFileName: "events_subscribe_list_resp1.json",
			respondFileName:      "error_response.json",
			args: args{
				o: SubscribeEventOptions{
					AllObjects: false,
					Objects: []Object{
						"79001294569",
					},
					Event:    RejectEvent,
					notifyTo: emailNotification,
					smsPhone: "",
					email:    "ndemidov@some.some",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subscriptionData, err := ioutil.ReadFile(filepath.Join(dataPath, tt.subscriptionFileName))
			if err != nil {
				t.Errorf("err: %s", err)
			}
			respondData, err := ioutil.ReadFile(filepath.Join(dataPath, tt.respondFileName))
			if err != nil {
				t.Errorf("err: %s", err)
			}

			subscriptionResp := httpmock.NewBytesResponder(200, subscriptionData)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/events_subscribe_list", subscriptionResp)

			subscribeResp := httpmock.NewBytesResponder(200, respondData)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/events_subscribe_add", subscribeResp)

			if err := api.AddEventSubscription(tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("API.AddEventSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//func TestAPI_ClearAllEventSubscriptions(t *testing.T) {
//	tests := []struct {
//		name    string
//		wantErr bool
//	}{
//		{ // TODO: Add test cases.
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := api.ClearAllEventSubscriptions(); (err != nil) != tt.wantErr {
//				t.Errorf("API.ClearAllEventSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestAPI_UnsubscribeObject(t *testing.T) {
//	type args struct {
//		o Object
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{ // TODO: Add test cases.
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			if err := api.UnsubscribeObject(tt.args.o); (err != nil) != tt.wantErr {
//				t.Errorf("API.UnsubscribeObject() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestAPI_ClearObjectEventSubscriptions(t *testing.T) {
//	type args struct {
//		o     Object
//		eType *EventType
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{ // TODO: Add test cases.
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			if err := api.ClearObjectEventSubscriptions(tt.args.o, tt.args.eType); (err != nil) != tt.wantErr {
//				t.Errorf("API.ClearObjectEventSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestAPI_ClearUnusedSubscriptions(t *testing.T) {
//	tests := []struct {
//		name    string
//		wantErr bool
//	}{
//		{ // TODO: Add test cases.
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			if err := api.ClearUnusedSubscriptions(); (err != nil) != tt.wantErr {
//				t.Errorf("API.ClearUnusedSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestAPI_removeObjectSubscriptions(t *testing.T) {
//	type args struct {
//		e SubscribedEvent
//		f shouldRemoveSubscription
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantErr bool
//	}{
//		{ // TODO: Add test cases.
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			if err := api.removeObjectSubscriptions(tt.args.e, tt.args.f); (err != nil) != tt.wantErr {
//				t.Errorf("API.removeObjectSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
