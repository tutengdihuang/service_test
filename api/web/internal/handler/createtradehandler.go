package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"service_test/api/web/internal/logic"
	"service_test/api/web/internal/svc"
	"service_test/api/web/internal/types"
)

func CreateTradeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateTradeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateTradeLogic(r.Context(), svcCtx)
		resp, err := l.CreateTrade(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
