package swagger

import (
	"net/http"
	"wano-island/common/core"

	"github.com/go-chi/render"
	"go.uber.org/fx"
)

type swaggerUIHandler struct {
}

func newSwaggerUIHandler() *swaggerUIHandler {
	return &swaggerUIHandler{}
}

func (handler *swaggerUIHandler) Config() *core.HTTPRouteConfig {
	return &core.HTTPRouteConfig{
		Pattern: "GET /swagger",
	}
}

func (handler *swaggerUIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	render.HTML(w, r, `<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="/static/swagger/swagger-ui.css" />
    <link rel="stylesheet" type="text/css" href="/static/swagger/index.css" />
    <link rel="icon" type="image/png" href="/static/swagger/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="/static/swagger/favicon-16x16.png" sizes="16x16" />
  </head>

  <body>
    <div id="swagger-ui"></div>
    <script src="/static/swagger/swagger-ui-bundle.js" charset="UTF-8"> </script>
    <script src="/static/swagger/swagger-ui-standalone-preset.js" charset="UTF-8"> </script>
    <script src="/static/swagger/swagger-initializer.js" charset="UTF-8"> </script>
  </body>
</html>
`)
}

func NewSwaggerModule() fx.Option {
	return fx.Module(
		"Swagger Module",
		fx.Provide(core.AsRoute(newSwaggerUIHandler)),
	)
}
