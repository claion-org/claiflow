package globvar

import (
	"encoding/base64"
	"fmt"
	"time"
)

var (
	ClusterTokenExpirationTime        = clusterTokenExpirationTime{Duration: 365 * 24 * time.Hour}
	ClientSessionSignatureSecret      = clientSessionSignatureSecret{Bytes: []byte{}}
	ClientSessionExpirationTime       = clientSessionExpirationTime{Duration: 60 * time.Second}
	ClientConfigServiceValidityPeriod = clientConfigServiceValidityPeriod{Duration: 10 * time.Minute}
)

type clusterTokenExpirationTime struct{ Duration time.Duration }

func (v clusterTokenExpirationTime) UUID() string {
	return "0f5658f37f2b45d881f19c7f56ea2e23"
}

func (v clusterTokenExpirationTime) Name() string {
	return "ClusterToken/ExpirationTime"
}

func (v clusterTokenExpirationTime) Summary() string {
	return fmt.Sprintf("cluster token expiration time (default=%q, format=duration)", v.GetValue())
}

func (v clusterTokenExpirationTime) GetValue() string {
	return v.Duration.String()
}

func (v *clusterTokenExpirationTime) SetValue(s string) error {
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	v.Duration = d
	return nil
}

func (v clusterTokenExpirationTime) Clone() GlobVar {
	return &v
}

func (v clusterTokenExpirationTime) Add(t time.Time) time.Time {
	if v.Duration == 0 {
		return t.Truncate(24 * time.Hour).Add(365 * 24 * time.Hour)
	}

	return t.Truncate(24 * time.Hour).Add(v.Duration)
}

type clientSessionSignatureSecret struct{ Bytes []byte }

func (v clientSessionSignatureSecret) UUID() string {
	return "77f7b2aeb0aa4254ad073ae7743291ab"
}

func (v clientSessionSignatureSecret) Name() string {
	return "ClientSession/SignatureSecret"
}

func (v clientSessionSignatureSecret) Summary() string {
	return fmt.Sprintf("client session signature secret (default=%q, format=base64)", v.GetValue())
}

func (v clientSessionSignatureSecret) GetValue() string {
	return base64.StdEncoding.EncodeToString(v.Bytes)
}

func (v *clientSessionSignatureSecret) SetValue(s string) error {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	v.Bytes = b
	return nil
}

func (v clientSessionSignatureSecret) Clone() GlobVar {
	return &v
}

type clientSessionExpirationTime struct{ Duration time.Duration }

func (v clientSessionExpirationTime) UUID() string {
	return "af9a14a58b254d13ae69c065a27811b6"
}

func (v clientSessionExpirationTime) Name() string {
	return "ClientSession/ExpirationTime"
}

func (v clientSessionExpirationTime) Summary() string {
	return fmt.Sprintf("client session expiration time (default=%q, format=duration)", v.GetValue())
}

func (v clientSessionExpirationTime) GetValue() string {
	return v.Duration.String()
}

func (v *clientSessionExpirationTime) SetValue(s string) error {
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	v.Duration = d
	return nil
}

func (v clientSessionExpirationTime) Clone() GlobVar {
	return &v
}

func (v clientSessionExpirationTime) Add(t time.Time) time.Time {
	if v.Duration == 0 {
		return t.Add(60 * time.Second)
	}

	return t.Add(v.Duration)
}

type clientConfigServiceValidityPeriod struct{ Duration time.Duration }

func (v clientConfigServiceValidityPeriod) UUID() string {
	return "bc2cd0f95b6d4db68870d30862523a04"
}

func (v clientConfigServiceValidityPeriod) Name() string {
	return "ClientConfig/ServiceValidityPeriod"
}

func (v clientConfigServiceValidityPeriod) Summary() string {
	return fmt.Sprintf("service validity period (default=%q, format=duration)", v.GetValue())
}

func (v clientConfigServiceValidityPeriod) GetValue() string {
	return v.Duration.String()
}

func (v *clientConfigServiceValidityPeriod) SetValue(s string) error {
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	v.Duration = d
	return nil
}

func (v clientConfigServiceValidityPeriod) Clone() GlobVar {
	return &v
}

func (v clientConfigServiceValidityPeriod) Add(t time.Time) time.Time {
	if v.Duration == 0 {
		return t.Truncate(time.Second).Add(-1 * 10 * time.Minute)
	}

	return t.Truncate(time.Second).Add(-1 * v.Duration)
}
