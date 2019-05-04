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
				t.Errorf("Positions.UnmarshalJSON() error = %v", err)
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

//
//func TestAPI_DeleteEventsSubscription(t *testing.T) {
//	type args struct {
//		id int64
//	}
//	tests := []struct {
//		name string
//
//		args    args
//		want    APIResponse
//		wantErr bool
//	}{
//		{ // TODO: Add test cases.
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			got, err := api.DeleteEventsSubscription(tt.args.id)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("API.DeleteEventsSubscription() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("API.DeleteEventsSubscription() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAPI_GetEventSubscriptions(t *testing.T) {
//	tests := []struct {
//		name    string
//		want    SubscribedEvents
//		wantErr bool
//	}{
//		{ // TODO: Add test cases.
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			got, err := api.GetEventSubscriptions()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("API.GetEventSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("API.GetEventSubscriptions() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestAPI_SubscribeEvent(t *testing.T) {
//	type args struct {
//		o SubscribeEventOptions
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    APIResponse
//		wantErr bool
//	}{
//		{ // TODO: Add test cases.
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			got, err := api.SubscribeEvent(tt.args.o)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("API.SubscribeEvent() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("API.SubscribeEvent() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
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
