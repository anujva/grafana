package clientmiddleware

import (
	"context"
	"net/http"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/plugins/manager/client/clienttest"
	"github.com/grafana/grafana/pkg/services/contexthandler/ctxkey"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/web"
)

func TestForwardIDMiddleware(t *testing.T) {
	pluginContext := backend.PluginContext{
		DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{},
	}

	t.Run("Should set forwarded id header if present", func(t *testing.T) {
		cdt := clienttest.NewClientDecoratorTest(t, clienttest.WithMiddlewares(NewForwardIDMiddleware()))

		ctx := context.WithValue(context.Background(), ctxkey.Key{}, &contextmodel.ReqContext{
			Context:      &web.Context{Req: &http.Request{}},
			SignedInUser: &user.SignedInUser{IDToken: "some-token"},
		})

		err := cdt.Decorator.CallResource(ctx, &backend.CallResourceRequest{
			PluginContext: pluginContext,
		}, nopCallResourceSender)
		require.NoError(t, err)

		require.Equal(t, "some-token", cdt.CallResourceReq.Headers[forwardIDHeaderName][0])
	})

	t.Run("Should not set forwarded id header if not present", func(t *testing.T) {
		cdt := clienttest.NewClientDecoratorTest(t, clienttest.WithMiddlewares(NewForwardIDMiddleware()))

		ctx := context.WithValue(context.Background(), ctxkey.Key{}, &contextmodel.ReqContext{
			Context:      &web.Context{Req: &http.Request{}},
			SignedInUser: &user.SignedInUser{},
		})

		err := cdt.Decorator.CallResource(ctx, &backend.CallResourceRequest{
			PluginContext: pluginContext,
		}, nopCallResourceSender)
		require.NoError(t, err)

		require.Len(t, cdt.CallResourceReq.Headers[forwardIDHeaderName], 0)
	})
}
