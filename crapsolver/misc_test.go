package crapsolver

import (
	"log"
	"reflect"
	"testing"
)

func TestGetRestrictions(t *testing.T) {
	type args struct {
		sitekey string
	}
	tests := []struct {
		name    string
		args    args
		want    *Restrictions
		wantErr bool
	}{
		{
			name: "discord login",
			args: args{
				sitekey: "f5561ba9-8f1e-40ca-9b5b-a0b3f719ef34",
			},
			want: &Restrictions{
				MinSubmitTime: 3200,
				MaxSubmitTime: 13000,
				AlwaysText:    true,
				Domain:        "discord.com",
				Enabled:       true,
				OneclickOnly:  false,
				Rate:          0,
			},
		},
		{
			name: "discord guild join",
			args: args{
				sitekey: "b2b02ab5-7dae-4d6f-830e-7b55634c888b",
			},
			want: &Restrictions{
				MinSubmitTime: 3200,
				MaxSubmitTime: 13000,
				AlwaysText:    true,
				Domain:        "discord.com",
				Enabled:       true,
				OneclickOnly:  false,
				Rate:          0,
			},
		},
		{
			name: "discord register",
			args: args{
				sitekey: "4c672d35-0701-42b2-88c3-78380b0db560",
			},
			want: &Restrictions{
				MinSubmitTime: 3200,
				MaxSubmitTime: 13000,
				AlwaysText:    true,
				Domain:        "discord.com",
				Enabled:       false,
				OneclickOnly:  false,
				Rate:          0,
			},
		},
		{
			name: "discord friend",
			args: args{
				sitekey: "a9b5fb07-92ff-493f-86fe-352a2803b3df",
			},
			want: &Restrictions{
				MinSubmitTime: 3200,
				MaxSubmitTime: 13000,
				AlwaysText:    true,
				Domain:        "discord.com",
				Enabled:       true,
				OneclickOnly:  false,
				Rate:          0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRestrictions(tt.args.sitekey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRestrictions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRestrictions() = %v, want %v", got, tt.want)
			}

			log.Println(tt.args.sitekey, got.Domain, got.Enabled)
		})
	}
}
