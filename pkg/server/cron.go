package server

import (
	"database/sql"
	"time"

	"github.com/claion-org/claiflow/pkg/logger"
	"github.com/claion-org/claiflow/pkg/server/database/vanilla/excute"
	"github.com/claion-org/claiflow/pkg/server/status/globvar"
	"github.com/claion-org/claiflow/pkg/server/status/ticker"
	"github.com/pkg/errors"
)

func Cron_GlobalVariables(db *sql.DB, dialect excute.SqlExcutor) (func(), error) {
	const interval = 10 * time.Second

	//환경설정 updater 생성
	updator := globvar.NewGlobalVariablesUpdate(db, dialect)
	//환경변수 리스트 검사
	if err := updator.WhiteListCheck(); err != nil {
		//빠져있는 환경변수 추가
		if err := updator.Merge(); err != nil {
			return nil, errors.Wrapf(err, "global variables init merge")
		}
	}
	//환경변수 업데이트
	if err := updator.Update(); err != nil {
		return nil, errors.Wrapf(err, "global variables init update")
	}

	// //에러 핸들러 등록
	// errorHandlers := ticker.HashsetErrorHandlers{}
	// errorHandlers.Add(func(err error) {
	// 	logger.Logger().Error(err, "cron jobs")
	// })

	//new ticker
	tickerClose := ticker.NewTicker(interval,
		//global variables update
		func() {
			if err := updator.Update(); err != nil {
				// errorHandlers.OnError(errors.Wrapf(err, "global variables update"))
				logger.Logger().Error(err, "global variables update")
			}
		},
	)

	return tickerClose, nil
}
