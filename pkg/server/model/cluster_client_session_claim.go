package model

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ClusterClientSessionClaim struct {
	ClusterUUID            string `json:"cluster-uuid,omitempty"`
	ClusterClientTokenUUID string `json:"cluster-client-token-uuid,omitempty"`
	UUID                   string `json:"uuid,omitempty"`
	IssuedAt               int64  `json:"iat,omitempty"`
	ExpiresAt              int64  `json:"exp,omitempty"`
	ClientVersion          string `json:"client-version,omitempty"`
	ClientLibVersion       string `json:"client_lib_version"`
}

func NewClusterClientSessionClaim(clusterUUID, clientTokenUUID, clientSessionUUID string, iat, exp time.Time, clientVersion, clientLibVersion string) *ClusterClientSessionClaim {
	return &ClusterClientSessionClaim{
		ClusterUUID:            clusterUUID,
		ClusterClientTokenUUID: clientTokenUUID,
		UUID:                   clientSessionUUID,
		IssuedAt:               iat.UTC().Unix(),
		ExpiresAt:              exp.UTC().Unix(),
		ClientVersion:          clientVersion,
		ClientLibVersion:       clientLibVersion,
	}
}

func (claim ClusterClientSessionClaim) Valid() error {
	vErr := new(jwt.ValidationError)
	now := time.Now().UTC().Unix()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if !claim.VerifyExpiresAt(now, false) {
		delta := time.Unix(now, 0).Sub(time.Unix(claim.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("%s by %s", jwt.ErrTokenExpired, delta)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if !claim.VerifyIssuedAt(now, false) {
		vErr.Inner = jwt.ErrTokenUsedBeforeIssued
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

func (claim ClusterClientSessionClaim) VerifyExpiresAt(cmp int64, req bool) bool {
	if claim.ExpiresAt == 0 {
		return verifyExp(nil, time.Unix(cmp, 0), req)
	}

	t := time.Unix(claim.ExpiresAt, 0)
	return verifyExp(&t, time.Unix(cmp, 0), req)
}

func (claim ClusterClientSessionClaim) VerifyIssuedAt(cmp int64, req bool) bool {
	if claim.IssuedAt == 0 {
		return verifyIat(nil, time.Unix(cmp, 0), req)
	}

	t := time.Unix(claim.IssuedAt, 0)
	return verifyIat(&t, time.Unix(cmp, 0), req)
}

func verifyExp(exp *time.Time, now time.Time, required bool) bool {
	if exp == nil {
		return !required
	}
	return now.Before(*exp)
}

func verifyIat(iat *time.Time, now time.Time, required bool) bool {
	if iat == nil {
		return !required
	}
	return now.After(*iat) || now.Equal(*iat)
}
