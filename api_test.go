// Package movizor package provides access to http://MoVizor.ru API
// which provides access for GSM geo-position services of russian telecommunications operators.
// Beeline, MTS, Megafon, Tele2.
//
// As soon as MoVizor provides service only in Russia all documentation will be in russian.
package movizor

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

var (
	api *API
)

func init() {
	httpmock.Activate()
	api, _ = NewMovizorAPI("some", "test")
}

func TestAPI_GeoCode(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name     string
		filename string
		args     args
		want     GeoPoints
		wantErr  bool
	}{
		{
			name:     "good_addr",
			filename: "distance_search_resp1.json",
			args: args{
				addr: "Москва, Вятская 27/13",
			},
			want: GeoPoints{
				{Coordinates: Coordinates{
					Lat: 55.8070325,
					Lon: 37.5795807,
				},
					Description: "Вятская улица, Савёловский, Савёловский район, Северный административный округ, Москва, ЦФО, 127015, Россия"},
			},
			wantErr: false,
		},
		{
			name:     "bad_addr",
			filename: "error_response.json",
			args: args{
				addr: "XXXXX",
			},
			want:    GeoPoints{},
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
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/distance_search", responder)
			got, err := api.GeoCode(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.GeoCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API.GeoCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_GetBalance(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		filenameWant string
		wantErr      bool
	}{
		{
			name:         "balance_good",
			filename:     "balance_resp1.json",
			filenameWant: "balance.json",
			wantErr:      false,
		},
		{
			name:         "balance_bad",
			filename:     "error_response.json",
			filenameWant: "balance.json",
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
			want := Balance{}
			if err := want.UnmarshalJSON(dWant); err != nil {
				t.Errorf("Balance.UnmarshalJSON() error = %v", err)
			}

			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/balance", responder)

			got, err := api.GetBalance()
			if (err != nil) != tt.wantErr {
				t.Errorf("API.GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (!tt.wantErr && !reflect.DeepEqual(got, want)) ||
				(tt.wantErr && !reflect.DeepEqual(got, Balance{})) {
				t.Errorf("API.GetBalance() = %v, want %v", got, want)
			}
		})
	}
}

func TestAPI_AddEditObject(t *testing.T) {
	type args struct {
		o  Object
		oo *ObjectOptions
	}
	tests := []struct {
		name     string
		args     args
		filename string
		want     APIResponse
		wantErr  bool
	}{
		{
			name: "object_add&edit_good",
			args: args{
				o: "+7(900)129-4567",
				oo: &ObjectOptions{
					Title:          "gopher uses movizor",
					Tags:           []string{"gopher", "movizor"},
					DateOff:        time.Now(),
					Tariff:         TariffEvery15,
					PackageProlong: false,
					Destinations: []DestinationOptions{
						{
							Text:         "Москва",
							Lon:          37.622504,
							Lat:          55.753215,
							ExpectedTime: time.Now(),
						},
						{
							Text:         "СПб",
							Lon:          30.315868,
							Lat:          59.939095,
							ExpectedTime: time.Now(),
						},
					},
					Metadata: map[string]string{
						"gopher":  "here",
						"movizor": "there",
					},
					CallToDriver: false,
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
			name: "object_add&edit_bad_object",
			args: args{
				o: "+7(900)129-456",
				oo: &ObjectOptions{
					Title:          "gopher uses movizor",
					Tags:           []string{"gopher", "movizor"},
					DateOff:        time.Now(),
					Tariff:         TariffEvery15,
					PackageProlong: false,
					Destinations: []DestinationOptions{
						{
							Text:         "Москва",
							Lon:          37.622504,
							Lat:          55.753215,
							ExpectedTime: time.Now(),
						},
						{
							Text:         "СПб",
							Lon:          30.315868,
							Lat:          59.939095,
							ExpectedTime: time.Now(),
						},
					},
					Metadata: map[string]string{
						"gopher":  "here",
						"movizor": "there",
					},
					CallToDriver: false,
				},
			},
			filename: "success_response.json",
			want:     APIResponse{},
			wantErr:  true,
		},
		{
			name: "object_add&edit_bad_destination",
			args: args{
				o: "+7(900)129-4567",
				oo: &ObjectOptions{
					Title:          "gopher uses movizor",
					Tags:           []string{"gopher", "movizor"},
					DateOff:        time.Now(),
					Tariff:         TariffEvery15,
					PackageProlong: false,
					Destinations: []DestinationOptions{
						{
							Text:         "",
							Lon:          37.622504,
							Lat:          55.753215,
							ExpectedTime: time.Now(),
						},
						{
							Text:         "СПб",
							Lon:          30.315868,
							Lat:          59.939095,
							ExpectedTime: time.Now(),
						},
					},
					Metadata: map[string]string{
						"gopher":  "here",
						"movizor": "there",
					},
					CallToDriver: false,
				},
			},
			filename: "success_response.json",
			want:     APIResponse{},
			wantErr:  true,
		},
		{
			name: "object_add&edit_bad",
			args: args{
				o: "+7(900)129-4567",
				oo: &ObjectOptions{
					Title:          "gopher uses movizor",
					Tags:           []string{"gopher", "movizor"},
					DateOff:        time.Now(),
					Tariff:         TariffEvery15,
					PackageProlong: false,
					Destinations: []DestinationOptions{
						{
							Text:         "Москва",
							Lon:          37.622504,
							Lat:          55.753215,
							ExpectedTime: time.Now(),
						},
						{
							Text:         "СПб",
							Lon:          30.315868,
							Lat:          59.939095,
							ExpectedTime: time.Now(),
						},
					},
					Metadata: map[string]string{
						"gopher":  "here",
						"movizor": "there",
					},
					CallToDriver: false,
				},
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
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/object_add", responder)
			got, err := api.AddObject(tt.args.o, tt.args.oo)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.AddObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.Result, tt.want.Result) {
				t.Errorf("API.AddObject() = %v, want %v", got, tt.want)
			}

			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/object_edit", responder)
			got, err = api.EditObject(tt.args.o, tt.args.oo)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.EditObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.Result, tt.want.Result) {
				t.Errorf("API.EditObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_AddObjectToSlave(t *testing.T) {
	type args struct {
		o       Object
		oo      *ObjectOptions
		slaveID uint64
	}
	tests := []struct {
		name     string
		args     args
		filename string
		want     APIResponse
		wantErr  bool
	}{
		{
			name: "object_add_good",
			args: args{
				o: "+7(900)129-4567",
				oo: &ObjectOptions{
					Title:          "gopher uses movizor",
					Tags:           []string{"gopher", "movizor"},
					DateOff:        time.Now(),
					Tariff:         TariffEvery15,
					PackageProlong: false,
					Destinations: []DestinationOptions{
						{
							Text:         "Москва",
							Lon:          37.622504,
							Lat:          55.753215,
							ExpectedTime: time.Now(),
						},
						{
							Text:         "СПб",
							Lon:          30.315868,
							Lat:          59.939095,
							ExpectedTime: time.Now(),
						},
					},
					Metadata: map[string]string{
						"gopher":  "here",
						"movizor": "there",
					},
					CallToDriver: false,
				},
				slaveID: 11112222,
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
			name: "object_add_good_zero",
			args: args{
				o: "+7(900)129-4567",
				oo: &ObjectOptions{
					Title:          "gopher uses movizor",
					Tags:           []string{"gopher", "movizor"},
					DateOff:        time.Now(),
					Tariff:         TariffEvery15,
					PackageProlong: false,
					Destinations: []DestinationOptions{
						{
							Text:         "Москва",
							Lon:          37.622504,
							Lat:          55.753215,
							ExpectedTime: time.Now(),
						},
						{
							Text:         "СПб",
							Lon:          30.315868,
							Lat:          59.939095,
							ExpectedTime: time.Now(),
						},
					},
					Metadata: map[string]string{
						"gopher":  "here",
						"movizor": "there",
					},
					CallToDriver: false,
				},
				slaveID: 0,
			},
			filename: "success_response.json",
			want: APIResponse{
				Result:     "success",
				ResultCode: "OK",
				Message:    "Some message",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := ioutil.ReadFile(filepath.Join(dataPath, tt.filename))
			if err != nil {
				t.Errorf("err: %s", err)
			}

			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/object_add", responder)
			got, err := api.AddObjectToSlave(tt.args.o, tt.args.oo, tt.args.slaveID)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.AddObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.Result, tt.want.Result) {
				t.Errorf("API.AddObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_GetObjectInfo(t *testing.T) {
	type args struct {
		o Object
	}
	tests := []struct {
		name         string
		args         args
		filename     string
		filenameWant string
		wantErr      bool
	}{
		{
			name: "object_get_good1",
			args: args{
				o: "79001294567",
			},
			filename:     "object_get_resp1.json",
			filenameWant: "object_get1.json",
			wantErr:      false,
		},
		{
			name: "object_get_bad",
			args: args{
				o: "79001294567",
			},
			filename:     "error_response.json",
			filenameWant: "object_get1.json",
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
			want := ObjectInfo{}
			if err := want.UnmarshalJSON(dWant); err != nil {
				t.Errorf("ObjectInfo.UnmarshalJSON() error = %v", err)
			}
			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/object_get", responder)

			got, err := api.GetObjectInfo(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.GetObjectInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (!tt.wantErr && !reflect.DeepEqual(got, want)) ||
				(tt.wantErr && !reflect.DeepEqual(got, ObjectInfo{})) {
				t.Errorf("API.GetObjectInfo() = %v, want %v", got, want)
			}
		})
	}
}

func TestAPI_EditObjectWithActivate(t *testing.T) {
	type args struct {
		o        Object
		oo       *ObjectOptions
		activate bool
	}
	tests := []struct {
		name     string
		args     args
		filename string
		want     APIResponse
		wantErr  bool
	}{
		{
			name: "object_edit_good",
			args: args{
				o: "+7(900)129-4567",
				oo: &ObjectOptions{
					Title:          "gopher uses movizor",
					Tags:           []string{"gopher", "movizor"},
					DateOff:        time.Now(),
					Tariff:         TariffEvery15,
					PackageProlong: false,
					Destinations: []DestinationOptions{
						{
							Text:         "Москва",
							Lon:          37.622504,
							Lat:          55.753215,
							ExpectedTime: time.Now(),
						},
						{
							Text:         "СПб",
							Lon:          30.315868,
							Lat:          59.939095,
							ExpectedTime: time.Now(),
						},
					},
					Metadata: map[string]string{
						"gopher":  "here",
						"movizor": "there",
					},
					CallToDriver: false,
				},
				activate: true,
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
			name: "object_edit_good_false",
			args: args{
				o: "+7(900)129-4567",
				oo: &ObjectOptions{
					Title:          "gopher uses movizor",
					Tags:           []string{"gopher", "movizor"},
					DateOff:        time.Now(),
					Tariff:         TariffEvery15,
					PackageProlong: false,
					Destinations: []DestinationOptions{
						{
							Text:         "Москва",
							Lon:          37.622504,
							Lat:          55.753215,
							ExpectedTime: time.Now(),
						},
						{
							Text:         "СПб",
							Lon:          30.315868,
							Lat:          59.939095,
							ExpectedTime: time.Now(),
						},
					},
					Metadata: map[string]string{
						"gopher":  "here",
						"movizor": "there",
					},
					CallToDriver: false,
				},
				activate: false,
			},
			filename: "success_response.json",
			want: APIResponse{
				Result:     "success",
				ResultCode: "OK",
				Message:    "Some message",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := ioutil.ReadFile(filepath.Join(dataPath, tt.filename))
			if err != nil {
				t.Errorf("err: %s", err)
			}

			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/object_edit", responder)
			got, err := api.EditObjectWithActivate(tt.args.o, tt.args.oo, tt.args.activate)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.EditObjectWithActivate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.Result, tt.want.Result) {
				t.Errorf("API.EditObjectWithActivate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_GetObjects(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		filenameWant string
		wantErr      bool
	}{
		{
			name:         "object_list_good",
			filename:     "object_list_resp1.json",
			filenameWant: "object_list.json",
			wantErr:      false,
		},
		{
			name:         "object_list_bad",
			filename:     "error_response.json",
			filenameWant: "object_list.json",
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
			want := ObjectsWithStatus{}
			if err := json.Unmarshal(dWant, &want); err != nil {
				t.Errorf("ObjectsWithStatus.UnmarshalJSON() error = %v", err)
			}

			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/object_list", responder)

			got, err := api.GetObjects()
			if (err != nil) != tt.wantErr {
				t.Errorf("API.GetObjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (!tt.wantErr && !reflect.DeepEqual(got, want)) ||
				(tt.wantErr && !reflect.DeepEqual(got, ObjectsWithStatus{})) {
				t.Errorf("API.GetObjects() = %v, want %v", got, want)
			}
		})
	}
}

func TestAPI_DeleteObject(t *testing.T) {
	type args struct {
		o Object
	}
	tests := []struct {
		name     string
		args     args
		filename string
		want     APIResponse
		wantErr  bool
	}{
		{
			name: "object_delete_good",
			args: args{
				o: "+7(900)129-4567",
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
			name: "object_delete_bad",
			args: args{
				o: "+7(900)129-4567",
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
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/object_delete", responder)
			got, err := api.DeleteObject(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.DeleteObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Result, tt.want.Result) {
				t.Errorf("API.DeleteObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_ReactivateObject(t *testing.T) {
	type args struct {
		o Object
	}
	tests := []struct {
		name     string
		args     args
		filename string
		want     APIResponse
		wantErr  bool
	}{
		{
			name: "object_reactivate_good",
			args: args{
				o: "+7(900)129-4567",
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
			name: "object_reactivate_bad1",
			args: args{
				o: "+7(900)129-4567",
			},
			filename: "error_response.json",
			want: APIResponse{
				Result:     "error",
				ResultCode: "ACCESS_DENIED",
				ErrorText:  "Auth rate limit exceeded",
			},
			wantErr: true,
		},
		{
			name: "object_reactivate_bad2",
			args: args{
				o: "+7(900)129-456",
			},
			filename: "error_response.json",
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
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/object_reactivate", responder)
			got, err := api.ReactivateObject(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.ReactivateObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Result, tt.want.Result) {
				t.Errorf("API.ReactivateObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_CancelTariffChangeObject(t *testing.T) {
	type args struct {
		o Object
	}
	tests := []struct {
		name     string
		args     args
		filename string
		want     APIResponse
		wantErr  bool
	}{
		{
			name: "object_cancel_tariff_good",
			args: args{
				o: "+7(900)129-4567",
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
			name: "object_cancel_tariff_bad1",
			args: args{
				o: "+7(900)129-4567",
			},
			filename: "error_response.json",
			want: APIResponse{
				Result:     "error",
				ResultCode: "ACCESS_DENIED",
				ErrorText:  "Auth rate limit exceeded",
			},
			wantErr: true,
		},
		{
			name: "object_cancel_tariff_bad2",
			args: args{
				o: "+7(900)129-456",
			},
			filename: "error_response.json",
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
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/object_cancel_tariff", responder)
			got, err := api.CancelTariffChangeObject(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.CancelTariffChangeObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Result, tt.want.Result) {
				t.Errorf("API.CancelTariffChangeObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_GetObjectLastPosition(t *testing.T) {
	type args struct {
		o Object
	}

	var i Int = 790
	tests := []struct {
		name     string
		args     args
		filename string
		want     Position
		wantErr  bool
	}{
		{
			name: "pos_last_good",
			args: args{
				o: "+7(900)129-4567",
			},
			filename: "pos_last_resp1.json",
			want: Position{
				Coordinates: Coordinates{
					Lat: 55.798355,
					Lon: 37.579491,
				},
				Timestamp:        Time(time.Unix(1548165275, 0)),
				TimestampRequest: Time(time.Unix(1548165275, 0)),
				Deviation:        &i,
				CoordinatesAttributes: CoordinatesAttributes{
					Distance:  nil,
					ETA:       nil,
					ETAStatus: nil,
					Place:     "Москва",
				},
			},
			wantErr: false,
		},
		{
			name: "pos_last_bad1",
			args: args{
				o: "+7(900)129-4567",
			},
			filename: "error_response.json",
			want:     Position{},
			wantErr:  true,
		},
		{
			name: "pos_last_bad2",
			args: args{
				o: "+7(900)129-456",
			},
			filename: "error_response.json",
			want:     Position{},
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
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/pos_last", responder)
			got, err := api.GetObjectLastPosition(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.GetObjectLastPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API.GetObjectLastPosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_GetObjectPositions(t *testing.T) {
	type args struct {
		o   Object
		rpo *RequestPositionsOptions
	}
	tests := []struct {
		name         string
		args         args
		filename     string
		filenameWant string
		wantErr      bool
	}{
		{
			name: "pos_list_good",
			args: args{
				o: "+7(900)129-4567",
				rpo: &RequestPositionsOptions{
					RequestLimit: 0,
					Offset:       0,
					TimeFrom:     time.Unix(1556452070, 0),
					TimeTo:       time.Unix(1556453032, 0),
				},
			},
			filename:     "pos_list_resp1.json",
			filenameWant: "pos_list.json",
			wantErr:      false,
		},
		{
			name: "pos_list_bad1",
			args: args{
				o: "+7(900)129-4567",
				rpo: &RequestPositionsOptions{
					RequestLimit: 0,
					Offset:       0,
					TimeFrom:     time.Unix(1556452070, 0),
					TimeTo:       time.Unix(1556453032, 0),
				},
			},
			filename:     "error_response.json",
			filenameWant: "pos_list.json",
			wantErr:      true,
		},
		{
			name: "pos_list_bad2",
			args: args{
				o: "+7(900)129-456",
				rpo: &RequestPositionsOptions{
					RequestLimit: 0,
					Offset:       0,
					TimeFrom:     time.Unix(1556452070, 0),
					TimeTo:       time.Unix(1556453032, 0),
				},
			},
			filename:     "pos_list_resp1.json",
			filenameWant: "pos_list.json",
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
			want := Positions{}
			if err := json.Unmarshal(dWant, &want); err != nil {
				t.Errorf("Positions.UnmarshalJSON() error = %v", err)
			}

			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/pos_list", responder)

			got, err := api.GetObjectPositions(tt.args.o, tt.args.rpo)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.GetObjectPositions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (!tt.wantErr && !reflect.DeepEqual(got, want)) ||
				(tt.wantErr && !reflect.DeepEqual(got, Positions{})) {
				t.Errorf("API.GetObjectPositions() = %v, want %v", got, want)
			}
		})
	}
}

func TestAPI_RequestPosition(t *testing.T) {
	type args struct {
		o Object
	}
	tests := []struct {
		name     string
		args     args
		filename string
		want     PositionRequest
		wantErr  bool
	}{
		{
			name: "pos_request_good",
			args: args{
				o: "+7(900)129-4567",
			},
			filename: "pos_request_resp1.json",
			want: PositionRequest{
				RequestID: 22022757,
			},
			wantErr: false,
		},
		{
			name: "pos_request_bad1",
			args: args{
				o: "+7(900)129-4567",
			},
			filename: "error_response.json",
			want:     PositionRequest{},
			wantErr:  true,
		},
		{
			name: "pos_request_bad2",
			args: args{
				o: "+7(900)129-456",
			},
			filename: "pos_request_resp1.json",
			want:     PositionRequest{},
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
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/pos_request", responder)

			got, err := api.RequestPosition(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.RequestPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API.RequestPosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_GetRequestedPosition(t *testing.T) {
	type args struct {
		pr PositionRequest
	}
	tests := []struct {
		name         string
		args         args
		filename     string
		filenameWant string
		wantErr      bool
	}{
		{
			name: "pos_get_good",
			args: args{
				pr: PositionRequest{
					RequestID: 22022757,
				},
			},
			filename:     "pos_get_resp1.json",
			filenameWant: "pos_get1.json",
			wantErr:      false,
		},
		{
			name: "pos_get_bad1",
			args: args{
				pr: PositionRequest{
					RequestID: 22022757,
				},
			},
			filename:     "error_response.json",
			filenameWant: "pos_get1.json",
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
			want := Position{}
			if err := json.Unmarshal(dWant, &want); err != nil {
				t.Errorf("Positions.UnmarshalJSON() error = %v", err)
			}

			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/pos_get", responder)

			got, err := api.GetRequestedPosition(tt.args.pr)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.GetRequestedPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (!tt.wantErr && !reflect.DeepEqual(got, want)) ||
				(tt.wantErr && !reflect.DeepEqual(got, Position{})) {
				t.Errorf("API.GetRequestedPosition() = %v, want %v", got, want)
			}
		})
	}
}

func TestAPI_GetObjectsPositions(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		filenameWant string
		wantErr      bool
	}{
		{
			name:         "pos_objects_good",
			filename:     "pos_objects_resp1.json",
			filenameWant: "pos_objects.json",
			wantErr:      false,
		},
		{
			name:         "pos_objects_bad",
			filename:     "error_response.json",
			filenameWant: "pos_objects.json",
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
			want := ObjectPositions{}
			if err := json.Unmarshal(dWant, &want); err != nil {
				t.Errorf("Positions.UnmarshalJSON() error = %v", err)
			}

			responder := httpmock.NewBytesResponder(200, d)
			httpmock.RegisterResponder("GET", "https://movizor.ru/api/some/pos_objects", responder)

			got, err := api.GetObjectsPositions()
			if (err != nil) != tt.wantErr {
				t.Errorf("API.GetObjectsPositions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (!tt.wantErr && !reflect.DeepEqual(got, want)) ||
				(tt.wantErr && !reflect.DeepEqual(got, ObjectPositions{})) {
				t.Errorf("API.GetObjectsPositions() = %v, want %v", got, want)
			}
		})
	}
}
