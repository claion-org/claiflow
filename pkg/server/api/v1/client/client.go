package client

import (
	context "context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/claion-org/claiflow/pkg/echov4"
	"github.com/claion-org/claiflow/pkg/server/control"
	"github.com/claion-org/claiflow/pkg/server/model"
	"github.com/claion-org/claiflow/pkg/server/status/globvar"
	"github.com/golang-jwt/jwt/v4"
)

const (
	HTTP_HEADER_X_OLD_CLIENT_TOKEN = "x-sudory-client-token"
)

func GetClientSessionClaims(ctx context.Context, header http.Header) (string, *model.ClusterClientSessionClaim, error) {
	var token string
	// get old header
	if 0 < len(header.Get(HTTP_HEADER_X_OLD_CLIENT_TOKEN)) {
		token = header.Get(HTTP_HEADER_X_OLD_CLIENT_TOKEN)
	}

	// get new header
	if _, token_, ok := echov4.ParseAuthorizationHeader(header); ok {
		token = token_
	}

	// check token length
	if len(token) == 0 {
		err := fmt.Errorf("token is empty")
		return token, nil, err
	}

	var claims model.ClusterClientSessionClaim
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return globvar.ClientSessionSignatureSecret.Bytes, nil
	})
	if err != nil {
		return token, &claims, err
	}

	// check claims type
	if _, ok := jwtToken.Claims.(*model.ClusterClientSessionClaim); !ok {
		err := fmt.Errorf("invalid claim type")
		return token, &claims, err
	}
	// check valid
	if !jwtToken.Valid {
		err := fmt.Errorf("invalid claim")
		return token, &claims, err
	}

	// check cluster client token record in control
	var exp sql.NullTime
	columns := map[string]interface{}{
		model.ClusterClientTokenFieldsExpirationTime.String(): &exp,
	}
	if err := control.GetClusterClientTokenColumns(ctx, claims.ClusterClientTokenUUID, columns); err != nil {
		return token, &claims, err
	}

	// vaild expiration time
	if !exp.Valid {
		err := fmt.Errorf("the expiration_time field in the cluster client token is empty (uuid=%v)", claims.ClusterClientTokenUUID)
		return token, &claims, err
	}
	if time.Now().UTC().After(exp.Time.UTC()) {
		err := fmt.Errorf("cluster client token has expired (uuid=%v)", claims.ClusterClientTokenUUID)
		return token, &claims, err
	}

	return token, &claims, nil
}

func KeepAliveClientSessionStatus(ctx context.Context, claimsToken string, claims model.ClusterClientSessionClaim, keepAliveTimeout time.Duration, errorHandler ...func(err error)) {
	OnError := func(err error) {
		for i := range errorHandler {
			errorHandler[i](err)
		}
	}

	columnListForUpdate := []string{
		model.ClusterClientTokenFieldsIssuedAtTime.String(),
		model.ClusterClientTokenFieldsExpirationTime.String(),
		model.ClusterClientTokenFieldsUpdated.String(),
	}

	SaveSession := func(now time.Time) time.Time {
		exp := now.Truncate(time.Second).Add(time.Second).Add(keepAliveTimeout)

		var session = model.ClusterClientSession{
			ClusterUUID:            claims.ClusterUUID,
			ClusterClientTokenUUID: claims.ClusterClientTokenUUID,
			UUID:                   claims.UUID,
			Token:                  claimsToken,
			IssuedAtTime:           time.Unix(claims.IssuedAt, 0),
			ExpirationTime:         exp,
			Created:                now,
			Updated:                now,
		}

		if err := control.UpsertClusterClientSession(ctx, &session, columnListForUpdate); err != nil {
			OnError(err)
		}

		return exp
	}

	const (
		CRON_INTERVAL = time.Second * 1
	)

	for {
		// session keep alive
		exp := SaveSession(time.Now())

		select {
		case <-ctx.Done():
			// session was closed
			// to reset the session expiration time to keepalive until reconnection
			fmt.Println(time.Until(exp), keepAliveTimeout/2, time.Until(exp) < keepAliveTimeout/2)
			if time.Until(exp) < keepAliveTimeout/2 {
				SaveSession(time.Now())
			}
			return
		case <-time.After(time.Until(exp.Add(CRON_INTERVAL * -1))):
		}
	}
}

func PollServiceFilter(limit int, exp time.Time, now time.Time) func(*model.ClusterService) bool {
	return func(s *model.ClusterService) bool {
		// 사이즈 제한
		if !(0 < limit) {
			return false
		}

		// 유효 시간
		if !s.Created.After(exp) {
			return false
		}

		limit--
		return true
	}
}

type Map = map[string]any

func EmbedFields(m Map, mm ...Map) Map {
	if m == nil {
		m = Map{}
	}

	for i := range mm {
		for k, v := range mm[i] {
			m[k] = v
		}
	}

	return m
}
