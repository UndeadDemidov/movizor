package movizor

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

//func TestDestinationOptions_addValuesTo(t *testing.T) {
//	type fields struct {
//		Text         string
//		Lon          float32
//		Lat          float32
//		ExpectedTime time.Time
//	}
//	type args struct {
//		idx int
//		v   *url.Values
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			do := DestinationOptions{
//				Text:         tt.fields.Text,
//				Lon:          tt.fields.Lon,
//				Lat:          tt.fields.Lat,
//				ExpectedTime: tt.fields.ExpectedTime,
//			}
//			if err := do.addValuesTo(tt.args.idx, tt.args.v); (err != nil) != tt.wantErr {
//				t.Errorf("DestinationOptions.addValuesTo() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestSchedulingOptions_WeekdayOn(t *testing.T) {
	type fields struct {
		weekdays [7]bool
		FireAt   []time.Time
	}
	type args struct {
		day Weekday
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "monday_on",
			fields: fields{
				weekdays: [7]bool{},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Monday,
			},
		},
		{
			name: "tuesday_on",
			fields: fields{
				weekdays: [7]bool{},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Tuesday,
			},
		},
		{
			name: "wednesday_on",
			fields: fields{
				weekdays: [7]bool{},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Wednesday,
			},
		},
		{
			name: "thursday_on",
			fields: fields{
				weekdays: [7]bool{},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Thursday,
			},
		},
		{
			name: "friday_on",
			fields: fields{
				weekdays: [7]bool{},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Friday,
			},
		},
		{
			name: "saturday_on",
			fields: fields{
				weekdays: [7]bool{},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Saturday,
			},
		},
		{
			name: "sunday_on",
			fields: fields{
				weekdays: [7]bool{},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Sunday,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SchedulingOptions{
				weekdays: tt.fields.weekdays,
				FireAt:   tt.fields.FireAt,
			}
			s.WeekdayOn(tt.args.day)
			if !s.IsWeekdayOn(tt.args.day) {
				t.Errorf("SchedulingOptions.WeekdayOn() error, Weekday num %v is not properly set.", tt.args.day)
			}
		})
	}
}

func TestSchedulingOptions_WeekdayOff(t *testing.T) {
	type fields struct {
		weekdays [7]bool
		FireAt   []time.Time
	}
	type args struct {
		day Weekday
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "monday_off",
			fields: fields{
				weekdays: [7]bool{true, true, true, true, true, true, true},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Monday,
			},
		},
		{
			name: "tuesday_off",
			fields: fields{
				weekdays: [7]bool{true, true, true, true, true, true, true},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Tuesday,
			},
		},
		{
			name: "wednesday_off",
			fields: fields{
				weekdays: [7]bool{true, true, true, true, true, true, true},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Wednesday,
			},
		},
		{
			name: "thursday_off",
			fields: fields{
				weekdays: [7]bool{true, true, true, true, true, true, true},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Thursday,
			},
		},
		{
			name: "friday_off",
			fields: fields{
				weekdays: [7]bool{true, true, true, true, true, true, true},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Friday,
			},
		},
		{
			name: "saturday_off",
			fields: fields{
				weekdays: [7]bool{true, true, true, true, true, true, true},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Saturday,
			},
		},
		{
			name: "sunday_off",
			fields: fields{
				weekdays: [7]bool{true, true, true, true, true, true, true},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Sunday,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SchedulingOptions{
				weekdays: tt.fields.weekdays,
				FireAt:   tt.fields.FireAt,
			}
			s.WeekdayOff(tt.args.day)
			if s.IsWeekdayOn(tt.args.day) {
				t.Errorf("SchedulingOptions.WeekdayOff() error, Weekday num %v is not properly set.", tt.args.day)
			}
		})
	}
}

func TestSchedulingOptions_IsWeekdayOn(t *testing.T) {
	type fields struct {
		weekdays [7]bool
		FireAt   []time.Time
	}
	type args struct {
		day Weekday
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "monday_on",
			fields: fields{
				weekdays: [7]bool{true, false, false, false, false, false, false},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Monday,
			},
			want: true,
		},
		{
			name: "monday_off",
			fields: fields{
				weekdays: [7]bool{false, true, true, true, true, true, true},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Monday,
			},
			want: false,
		},
		{
			name: "sunday_on",
			fields: fields{
				weekdays: [7]bool{false, false, false, false, false, false, true},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Sunday,
			},
			want: true,
		},
		{
			name: "sunday_off",
			fields: fields{
				weekdays: [7]bool{true, true, true, true, true, true, false},
				FireAt:   []time.Time{},
			},
			args: args{
				day: Sunday,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SchedulingOptions{
				weekdays: tt.fields.weekdays,
				FireAt:   tt.fields.FireAt,
			}
			if got := s.IsWeekdayOn(tt.args.day); got != tt.want {
				t.Errorf("SchedulingOptions.IsWeekdayOn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchedulingOptions_addValuesTo(t *testing.T) {
	type fields struct {
		weekdays [7]bool
		FireAt   []time.Time
	}
	type args struct {
		v *url.Values
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "time is not set",
			fields: fields{
				weekdays: [7]bool{true, false, false, false, false, false, false},
				FireAt:   []time.Time{},
			},
			args: args{
				v: &url.Values{},
			},
			wantErr: true,
		},
		{
			name: "day is not set",
			fields: fields{
				weekdays: [7]bool{false, false, false, false, false, false, false},
				FireAt: []time.Time{
					time.Now(),
				},
			},
			args: args{
				v: &url.Values{},
			},
			wantErr: true,
		},
		{
			name: "values are nil",
			fields: fields{
				weekdays: [7]bool{false, false, false, false, false, false, false},
				FireAt: []time.Time{
					time.Now(),
				},
			},
			args:    args{nil},
			wantErr: true,
		},
		{
			name: "day & time are set",
			fields: fields{
				weekdays: [7]bool{true, false, true, false, true, false, false},
				FireAt: []time.Time{
					time.Now(),
				},
			},
			args: args{
				v: &url.Values{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SchedulingOptions{
				weekdays: tt.fields.weekdays,
				FireAt:   tt.fields.FireAt,
			}
			if err := s.addValuesTo(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SchedulingOptions.addValuesTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//func TestObjectOptions_addValuesTo(t *testing.T) {
//	type fields struct {
//		Title          string
//		Tags           []string
//		DateOff        time.Time
//		Tariff         TariffType
//		PackageProlong bool
//		Destinations   []DestinationOptions
//		Schedules      *SchedulingOptions
//		Metadata       map[string]string
//		CallToDriver   bool
//	}
//	type args struct {
//		v *url.Values
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			o := &ObjectOptions{
//				Title:          tt.fields.Title,
//				Tags:           tt.fields.Tags,
//				DateOff:        tt.fields.DateOff,
//				Tariff:         tt.fields.Tariff,
//				PackageProlong: tt.fields.PackageProlong,
//				Destinations:   tt.fields.Destinations,
//				Schedules:      tt.fields.Schedules,
//				Metadata:       tt.fields.Metadata,
//				CallToDriver:   tt.fields.CallToDriver,
//			}
//			if err := o.addValuesTo(tt.args.v); (err != nil) != tt.wantErr {
//				t.Errorf("ObjectOptions.addValuesTo() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestRequestPositionsOptions_addValuesTo(t *testing.T) {
	type fields struct {
		RequestLimit uint64
		Offset       uint64
		TimeFrom     time.Time
		TimeTo       time.Time
	}
	type args struct {
		v *url.Values
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test nice",
			fields: fields{
				10, 10, time.Now(), time.Now(),
			},
			args: args{
				v: &url.Values{},
			},
			wantErr: false,
		},
		{
			name: "test not nice",
			fields: fields{
				10, 10, time.Now(), time.Now(),
			},
			args:    args{nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rpo := &RequestPositionsOptions{
				RequestLimit: tt.fields.RequestLimit,
				Offset:       tt.fields.Offset,
				TimeFrom:     tt.fields.TimeFrom,
				TimeTo:       tt.fields.TimeTo,
			}
			if err := rpo.addValuesTo(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("RequestPositionsOptions.addValuesTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestObjectEventsOptions_values(t *testing.T) {
	type fields struct {
		RequestLimit uint64
		AfterEventID uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   url.Values
	}{
		{
			name:   "",
			fields: fields{},
			want:   url.Values{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eo := ObjectEventsOptions{
				RequestLimit: tt.fields.RequestLimit,
				AfterEventID: tt.fields.AfterEventID,
			}
			if got := eo.values(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectEventsOptions.values() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestNewSubscribeEventOptions(t *testing.T) {
//	type args struct {
//		o Object
//		e EventType
//	}
//	tests := []struct {
//		name string
//		args args
//		want SubscribeEventOptions
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewSubscribeEventOptions(tt.args.o, tt.args.e); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewSubscribeEventOptions() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestSubscribeEventOptions_SetSMSNotification(t *testing.T) {
//	type fields struct {
//		AllObjects bool
//		Objects    []Object
//		Event      EventType
//		notifyTo   notificationType
//		smsPhone   Object
//		email      string
//	}
//	type args struct {
//		phone Object
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			se := &SubscribeEventOptions{
//				AllObjects: tt.fields.AllObjects,
//				Objects:    tt.fields.Objects,
//				Event:      tt.fields.Event,
//				notifyTo:   tt.fields.notifyTo,
//				smsPhone:   tt.fields.smsPhone,
//				email:      tt.fields.email,
//			}
//			if err := se.SetSMSNotification(tt.args.phone); (err != nil) != tt.wantErr {
//				t.Errorf("SubscribeEventOptions.SetSMSNotification() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestSubscribeEventOptions_SetEMailNotification(t *testing.T) {
//	type fields struct {
//		AllObjects bool
//		Objects    []Object
//		Event      EventType
//		notifyTo   notificationType
//		smsPhone   Object
//		email      string
//	}
//	type args struct {
//		mail string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			se := &SubscribeEventOptions{
//				AllObjects: tt.fields.AllObjects,
//				Objects:    tt.fields.Objects,
//				Event:      tt.fields.Event,
//				notifyTo:   tt.fields.notifyTo,
//				smsPhone:   tt.fields.smsPhone,
//				email:      tt.fields.email,
//			}
//			if err := se.SetEMailNotification(tt.args.mail); (err != nil) != tt.wantErr {
//				t.Errorf("SubscribeEventOptions.SetEMailNotification() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestSubscribeEventOptions_SetTelegramNotification(t *testing.T) {
//	type fields struct {
//		AllObjects bool
//		Objects    []Object
//		Event      EventType
//		notifyTo   notificationType
//		smsPhone   Object
//		email      string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			se := &SubscribeEventOptions{
//				AllObjects: tt.fields.AllObjects,
//				Objects:    tt.fields.Objects,
//				Event:      tt.fields.Event,
//				notifyTo:   tt.fields.notifyTo,
//				smsPhone:   tt.fields.smsPhone,
//				email:      tt.fields.email,
//			}
//			se.SetTelegramNotification()
//		})
//	}
//}
//
//func TestSubscribeEventOptions_values(t *testing.T) {
//	type fields struct {
//		AllObjects bool
//		Objects    []Object
//		Event      EventType
//		notifyTo   notificationType
//		smsPhone   Object
//		email      string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		want    url.Values
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			se := SubscribeEventOptions{
//				AllObjects: tt.fields.AllObjects,
//				Objects:    tt.fields.Objects,
//				Event:      tt.fields.Event,
//				notifyTo:   tt.fields.notifyTo,
//				smsPhone:   tt.fields.smsPhone,
//				email:      tt.fields.email,
//			}
//			got, err := se.values()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("SubscribeEventOptions.values() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("SubscribeEventOptions.values() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
