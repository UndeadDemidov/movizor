// Package movizor package provides access to http://MoVizor.ru API
// which provides access for GSM geo-position services of russian telecommunications operators.
// Beeline, MTS, Megafon, Tele2.
//
// As soon as MoVizor provides service only in Russia all documentation will be in russian.
package movizor

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestAPI_GeoCode(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	api, err := NewMovizorAPI("some", "test")
	if err != nil {
		t.Errorf("API.GeoCode() can't create instance of API: %s", err)
	}

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
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	api, err := NewMovizorAPI("some", "test")
	if err != nil {
		t.Errorf("API.GeoCode() can't create instance of API: %s", err)
	}

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

func TestAPI_AddObject(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	api, err := NewMovizorAPI("some", "test")
	if err != nil {
		t.Errorf("API.GeoCode() can't create instance of API: %s", err)
	}

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
			name: "object_add_bad_object",
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
			name: "object_add_bad_destination",
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
			name: "object_add_bad",
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
		})
	}
}

func TestAPI_AddObjectToSlave(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	api, err := NewMovizorAPI("some", "test")
	if err != nil {
		t.Errorf("API.GeoCode() can't create instance of API: %s", err)
	}

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
