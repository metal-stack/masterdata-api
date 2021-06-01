package mapper

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/metal-stack/masterdata-api/api/rest/v1"
)

func TestTenantMapperRoundtrip(t *testing.T) {
	tests := []struct {
		name     string
		inAndOut *v1.Tenant
	}{
		{
			name: "full roundtrip",
			inAndOut: &v1.Tenant{
				Meta: &v1.Meta{
					Id:          "1",
					Kind:        "tenant",
					Apiversion:  "1",
					Version:     99,
					CreatedTime: mustParseTimeP("2021-01-21"),
					UpdatedTime: mustParseTimeP("2021-03-01"),
					Annotations: map[string]string{
						"a": "first",
						"b": "second",
					},
					Labels: []string{
						"first label",
					},
				},
				Name:        "tnt",
				Description: "tnt is a test tenant",
				DefaultQuotas: &v1.QuotaSet{
					Cluster: &v1.Quota{
						Quota: int32p(100),
					},
					Machine: &v1.Quota{
						Quota: int32p(10),
					},
					Ip: &v1.Quota{
						Quota: int32p(20),
					},
					Project: &v1.Quota{
						Quota: int32p(11),
					},
				},
				Quotas: &v1.QuotaSet{
					Cluster: &v1.Quota{
						Quota: int32p(100),
					},
					Machine: &v1.Quota{
						Quota: int32p(72),
					},
					Ip: &v1.Quota{
						Quota: int32p(30),
					},
					Project: &v1.Quota{
						Quota: int32p(7),
					},
				},
				IAMConfig: &v1.IAMConfig{
					IssuerConfig: &v1.IssuerConfig{
						URL:      "https://oidc.myissuer.com",
						ClientID: "47abcdef12",
					},
					IDMConfig: &v1.IDMConfig{
						IDMType: "ldap",
					},
				},
			},
		},
		{
			name: "minimal roundtrip",
			inAndOut: &v1.Tenant{
				Meta: &v1.Meta{
					Id:          "1",
					Kind:        "tenant",
					Apiversion:  "1",
					Version:     99,
					CreatedTime: mustParseTimeP("2021-01-21"),
					UpdatedTime: mustParseTimeP("2021-03-01"),
				},
				Name:        "tnt",
				Description: "tnt is a test tenant",
				IAMConfig: &v1.IAMConfig{
					IssuerConfig: &v1.IssuerConfig{
						URL:      "https://oidc.myissuer.com",
						ClientID: "47abcdef12",
					},
					IDMConfig: &v1.IDMConfig{
						IDMType: "ldap",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			gotMdMV1 := ToMdmV1Tenant(tt.inAndOut)
			gotV1 := ToV1Tenant(gotMdMV1)

			if !reflect.DeepEqual(tt.inAndOut, gotV1) {
				diff := cmp.Diff(tt.inAndOut, gotV1)
				t.Errorf("in = %v, want %v, diff %s", tt.inAndOut, gotV1, diff)
			}
		})
	}
}

func int32p(i int) *int32 {
	i32 := int32(i)
	return &i32
}

func mustParseTimeP(ts string) *time.Time {
	t, err := time.Parse("2006-01-02", ts)
	if err != nil {
		panic(err)
	}
	return timep(t)
}

func timep(t time.Time) *time.Time {
	return &t
}
