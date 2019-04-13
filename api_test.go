// Package movizor package provides access to http://MoVizor.ru API
// which provides access for GSM geo-position services of russian telecommunications operators.
// Beeline, MTS, Megafon, Tele2.
//
// As soon as MoVizor provides service only in Russia all documentation will be in russian.
package movizor

import (
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
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
			name:     "Москва, Вятская 27/13",
			filename: "distance_search1.json",
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
			name:     "XXXXX",
			filename: "distance_search2.json",
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
