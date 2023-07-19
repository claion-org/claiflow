package fetcher

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/claion-org/claiflow/pkg/client/internal"
	"github.com/claion-org/claiflow/pkg/client/log"
	apiclient "github.com/claion-org/claiflow/pkg/server/api/client"
	sessionv1 "github.com/claion-org/claiflow/pkg/server/model"
)

func (f *Fetcher) HandShake() error {
	body := &apiclient.AuthRequestV1{
		ClusterUuid:      f.clusterId,
		Assertion:        f.bearerToken,
		ClientVersion:    f.clientVersion,
		ClientLibVersion: internal.ClientLibraryVersion,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Debugf("requested to handshake: cluster_uuid=%s, assertion=%s, client_version=%s\n", body.ClusterUuid, body.Assertion, body.ClientVersion)
	if err := f.serverAPI.Auth(ctx, body); err != nil {
		return err
	}
	sessionToken := f.serverAPI.GetToken()
	log.Debugf("received handshake request: token=%s\n", sessionToken)

	f.ChangeClientConfigFromToken()

	// save session_uuid from token
	claims := new(sessionv1.ClusterClientSessionClaim)
	jwt_token, _, err := jwt.NewParser().ParseUnverified(sessionToken, claims)
	if _, ok := jwt_token.Claims.(*sessionv1.ClusterClientSessionClaim); !ok || err != nil {
		if err == nil {
			err = fmt.Errorf("unable to convert token.claims to *sessionv1.ClientSessionPayload")
		}
		return err
	}
	if err := writeFile(".session", []byte(claims.UUID)); err != nil {
		return err
	}

	return nil
}

func (f *Fetcher) RetryHandshake() {
	maxRetryCnt := 5

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for retry := 0; ; <-ticker.C {
		log.Debugf("retry handshake: count=%d\n", retry+1)
		if err := f.HandShake(); err != nil {
			log.Warnf("failed to handshake retry: count=%d, error=%v\n", retry, err)
		} else {
			return
		}
		retry++

		if maxRetryCnt <= retry {
			f.Cancel()
			return
		}
	}
}

func writeFile(filename string, data []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}
