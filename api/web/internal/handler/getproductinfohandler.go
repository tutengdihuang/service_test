package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"service_test/api/web/internal/logic"
	"service_test/api/web/internal/svc"
)

func GetProductInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetProductInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetProductInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
