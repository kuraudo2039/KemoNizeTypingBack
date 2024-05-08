package coyoteHttpApi

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"

	errorObj "gin_test/coyote/obj/error"
	roomObj "gin_test/coyote/obj/room"
	"gin_test/coyote/util"
)

func enterRoom(client *firestore.Client) func(*gin.Context) {
	return func(c *gin.Context) {
		// コンテキストを初期化
		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
		defer cancel()

		// パラメータ取得
		type ReqData struct {
			Password string `json:"password" binding:"required"`
			Name     string `json:"name" binding:"required"`
		}

		var reqData ReqData
		if err := c.ShouldBindJSON(&reqData); err != nil {
			util.Log(util.LogObj{"error", err.Error()})
			c.IndentedJSON(http.StatusBadRequest, errorObj.CreateErr(err))
			return
		}
		util.Log(util.LogObj{"requested enterRoom()", reqData})

		var resData *roomObj.Room
		// 1-1. パラメータを元に対象パスワードのルームがあるか&名前が重複していないか確認

		resData, err := roomObj.GetRoom(client, ctx, reqData.Password)
		if err != nil {
			util.Log(util.LogObj{"error", err.Error()})
			c.IndentedJSON(http.StatusInternalServerError, errorObj.CreateErr(err))
			return
		}

		if resData != nil {
			// 2-0. あったらメンバー重複チェック
			if resData.IsExistMember(reqData.Name) {
				util.Log(util.LogObj{"error(Conflict member name)", reqData.Name})
				c.IndentedJSON(http.StatusBadRequest, errorObj.CreateErrFromString("入室に失敗しました。\n既に使われている名前です。", 409))
				return
			}
			// 2-1. 問題なければ200応答
			c.IndentedJSON(http.StatusOK, resData)
		} else {
			// 2-2. 無かったらパスワードを元にルームを作成して201応答
			resData, err := roomObj.CreateRoom(client, ctx, roomObj.RoomData{reqData.Password, 0})
			if err != nil {
				util.Log(util.LogObj{"error", err.Error()})
				c.IndentedJSON(http.StatusInternalServerError, errorObj.CreateErr(err))
				return
			}
			c.IndentedJSON(http.StatusCreated, resData)
		}
	}
}
