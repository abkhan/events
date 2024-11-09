package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"

	"events/server/response"
)

type ctxKey string

const (
	applicationKey ctxKey = "service_api_key"
	partnerKey     ctxKey = "service_partner_key"
)

var errStatusNotFound = &response.Error{
	Status:  response.StatusNotFound,
	Message: "status not found",
	Code:    http.StatusNotFound,
}

type initiator struct {
	keys         map[string]string
	totalCounter *prometheus.CounterVec
}

func NewInitiatorMiddleware(keys string) (gin.HandlerFunc, error) {
	apiKeys := strings.Split(keys, ",")

	k := make(map[string]string, len(apiKeys))
	for i := range apiKeys {
		pair := strings.SplitN(apiKeys[i], ":", 2)
		if len(pair) < 2 {
			return nil, fmt.Errorf("invalid api_keys: key \"%s\"", apiKeys[i])
		}

		k[pair[0]] = pair[1]
	}

	totalCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "fixme_ns_initiator",
			Name:      "total_request_to_image_service",
			Help:      "Count of requests to image service",
		},
		[]string{"api_key", "partner"},
	)

	return (&initiator{
		keys:         k,
		totalCounter: totalCounter,
	}).Handle, nil
}

func (u *initiator) isExists(key string) bool {

	logrus.Infof("Keys: %v", u.keys)

	_, ok := u.keys[key]

	return ok
}

var validate = validator.New()

func (u *initiator) Handle(c *gin.Context) {
	//uuid := c.Param("uuid")
	uuid := c.Request.Header["X-Api-Key"][0]
	logrus.Infof("uuid in request: %s", uuid)

	if err := validate.Var(uuid, "required,uuid4"); err != nil {
		u.checkAPIKey(c, uuid)
		return
	}

	// TODO: what projecte name we should validate against?
	//c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), partnerKey, p.Name))

	u.saveMetric(c.Request.Context())
}

func (u *initiator) checkAPIKey(c *gin.Context, uuid string) {
	if !u.isExists(uuid) {
		response.Abort(c, errStatusNotFound.WithError(fmt.Errorf("undefined api key: %s", uuid)))
		return
	}

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), applicationKey, uuid))

	u.saveMetric(c.Request.Context())
}

func GetAPIKey(ctx context.Context) string {
	if key := ctx.Value(applicationKey); key != nil {
		return key.(string)
	}

	return ""
}

func GetPartnerName(ctx context.Context) string {
	if partner := ctx.Value(partnerKey); partner != nil {
		return partner.(string)
	}

	return ""
}

func (u *initiator) saveMetric(ctx context.Context) {
	labels := []string{GetAPIKey(ctx), GetPartnerName(ctx)}
	u.totalCounter.WithLabelValues(labels...).Inc()
}
