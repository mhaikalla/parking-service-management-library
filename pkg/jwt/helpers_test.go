package jwt

import (
	"reflect"
	"testing"
)

func TestTransCIAMClaims(t *testing.T) {

	ciamClaims := CIAMClaims{
		Subject:      "test trans",
		Issuer:       "middleware",
		SubscriberID: "88989887987",
		Audience:     "something",
		MSISDN:       "0899889887878",
		DeviceID:     "aksdjkjasdlkjasldkj",
		SubsType:     "SUBSTYPE",
	}

	targetClaims := JWTClaims{
		JWTID:    "test trans",
		Issuer:   "middleware",
		ClientID: "aksdjkjasdlkjasldkj;88989887987;0899889887878",
		Subject:  "something",
		Audience: "SUBSTYPE",
		MSISDN:   "0899889887878",
		SubsID:   "88989887987",
		DeviceID: "aksdjkjasdlkjasldkj",
	}

	type args struct {
		claims CIAMClaims
	}
	tests := []struct {
		name    string
		args    args
		want    JWTClaims
		wantErr bool
	}{
		{"1", args{claims: ciamClaims}, targetClaims, false},
		{
			"2",
			args{
				claims: CIAMClaims{
					Subject:      "test trans",
					Issuer:       "middleware",
					SubscriberID: "88989887987",
					Audience:     "something",
					MSISDN:       "0899889887878",
					DeviceID:     "aksdjkjasdlkjasldkj",
					// SubsType:     "SUBSTYPE",
				},
			},
			JWTClaims{},
			true,
		},
		{
			"3",
			args{
				claims: CIAMClaims{
					Subject:      "test trans",
					Issuer:       "middleware",
					SubscriberID: "88989887987",
					Audience:     "something",
					MSISDN:       "0899889887878",
					// DeviceID:     "aksdjkjasdlkjasldkj",
					SubsType: "SUBSTYPE",
				},
			},
			JWTClaims{},
			true,
		},
		{
			"4",
			args{
				claims: CIAMClaims{
					Subject:      "test trans",
					Issuer:       "middleware",
					SubscriberID: "88989887987",
					Audience:     "something",
					// MSISDN:       "0899889887878",
					DeviceID: "aksdjkjasdlkjasldkj",
					SubsType: "SUBSTYPE",
				},
			},
			JWTClaims{},
			true,
		},
		{
			"5",
			args{
				claims: CIAMClaims{
					Subject: "test trans",
					Issuer:  "middleware",
					// SubscriberID: "88989887987",
					Audience: "something",
					MSISDN:   "0899889887878",
					DeviceID: "aksdjkjasdlkjasldkj",
					SubsType: "SUBSTYPE",
				},
			},
			JWTClaims{},
			true,
		},
		{
			"5",
			args{
				claims: CIAMClaims{
					Subject:      "test trans",
					Issuer:       "middleware",
					SubscriberID: "88989887987",
					Audience:     "something",
					MSISDN:       "0899889887878",
					DeviceID:     "aksdjkjasdlkjasldkj",
					SubsType:     HomeFiberSubsType,
				},
			},
			JWTClaims{},
			true,
		},
		{
			"6",
			args{
				claims: CIAMClaims{
					Subject:      "test trans",
					Issuer:       "middleware",
					SubscriberID: "88989887987",
					Audience:     "something",
					MSISDN:       "0899889887878",
					DeviceID:     "aksdjkjasdlkjasldkj",
					AccountID:    "99889987",
					CustomerID:   "887766544",
					Email:        "foo@bar",
					SubsType:     HomeFiberSubsType,
				},
			},
			JWTClaims{
				JWTID:    "test trans",
				Issuer:   "middleware",
				ClientID: ";88989887987;99889987",
				Subject:  "something",
				Audience: "HOMEFIBER",
				MSISDN:   "99889987",
				SubsID:   "88989887987",
				DeviceID: "aksdjkjasdlkjasldkj",
			},

			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransCIAMClaims(tt.args.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransCIAMClaims() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransCIAMClaims() = %v, want %v", got, tt.want)
			}
		})
	}
}
