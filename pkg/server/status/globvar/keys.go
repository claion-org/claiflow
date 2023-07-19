package globvar

// import (
// 	"time"
// )

// //go:generate go run -mod=mod github.com/abice/go-enum@v0.5.3 --file=keys.go --names --nocase

// /*
// ENUM(
// BearerToken/ExpirationTime

// ClientSession/SignatureSecret
// ClientSession/ExpirationTime

// ClientConfig/ServiceValidityPeriod
// )
// */
// type Key int

// type StoreManager struct {
// 	Store map[Key]func(s string) error
// }

// func (gvset *StoreManager) Setter(gv Key, fn func(s string) error) {
// 	if gvset.Store == nil {
// 		gvset.Store = map[Key]func(s string) error{}
// 	}
// 	gvset.Store[gv] = fn
// }

// func (gvset *StoreManager) Call(gv Key, s string) error {
// 	fn, ok := gvset.Store[gv]
// 	if !ok {
// 		return nil
// 	}

// 	return fn(s)
// }

// var storeManager *StoreManager

// func init() {
// 	storeManager = &StoreManager{}
// 	for k, v := range defaultValueSet {
// 		storeManager.Setter(k, v.Setter)
// 	}
// }

// type bearerToken struct {
// 	// bearer token expiration time
// 	expirationTime time.Duration // 365 * 24 * 1h
// }

// func (value bearerToken) ExpirationTime(t time.Time) time.Time {
// 	if true {
// 		if value.expirationTime == 0 {
// 			return t.Truncate(24 * time.Hour).Add(365 * 24 * time.Hour)
// 		}

// 		return t.Truncate(24 * time.Hour).Add(value.expirationTime)
// 	}

// 	if value.expirationTime == 0 {
// 		return t.Add(365 * 24 * time.Hour)
// 	}

// 	return t.Add(value.expirationTime)
// }

// type clientSession struct {
// 	// client session signature secret
// 	signatureSecret string
// 	// client session expiration time
// 	expirationTime time.Duration // 60s
// }

// func (value clientSession) SignatureSecret() string {
// 	return value.signatureSecret
// }

// func (value clientSession) ExpirationTime(t time.Time) time.Time {
// 	if value.expirationTime == 0 {
// 		return t.Add(60 * time.Second)
// 	}

// 	return t.Add(value.expirationTime)
// }

// type clientConfig struct {
// 	// service validity period
// 	serviceValidityPeriod time.Duration // 10m
// }

// func (value clientConfig) ServiceValidityPeriod() time.Duration {
// 	if value.serviceValidityPeriod == 0 {
// 		return 10 * time.Minute
// 	}

// 	return value.serviceValidityPeriod
// }

// var (
// 	BearerToken = bearerToken{
// 		expirationTime: 365 * 24 * time.Hour,
// 	}
// 	ClientSession = clientSession{
// 		signatureSecret: "",
// 		expirationTime:  60 * time.Second,
// 	}
// 	ClientConfig = clientConfig{
// 		serviceValidityPeriod: 10 * time.Minute,
// 	}
// )
