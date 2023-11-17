package submitter

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/jonkerj/chargeflux/pkg/smartevse"
)

type (
	Submitter struct {
		context  context.Context
		url      string
		writeAPI api.WriteAPIBlocking
		tags     map[string]string
	}

	cfPoint struct {
		Point *write.Point
	}
)

func NewCfPoint(name string, tags map[string]string) *cfPoint {
	p := write.NewPointWithMeasurement(name)
	for tag, value := range tags {
		p.AddTag(tag, value)
	}

	return &cfPoint{
		Point: p,
	}
}

func (c *cfPoint) visit(field string, val int64) {
	slog.Debug("adding field", "field", field, "value", val)
	c.Point.AddField(field, val)
}

func NewSubmitter(smartEvseUrl string, influxdbUrl string, influxdbToken string, influxdbOrg string, influxdbBucket string, tags string) (*Submitter, error) {
	ctx := context.TODO()

	client := influxdb2.NewClient(influxdbUrl, influxdbToken)
	health, err := client.Health(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting status from influxdb: %v", err)
	}

	if health.Status != domain.HealthCheckStatusPass {
		return nil, fmt.Errorf("influxdb server was not healthy, status=%v", health.Status)
	}

	tagsAsStrings := strings.Split(tags, ",")
	tagMap := make(map[string]string)
	for _, tagAsString := range tagsAsStrings {
		kv := strings.Split(tagAsString, "=")
		tagMap[kv[0]] = kv[1]
	}

	return &Submitter{
		context:  ctx,
		url:      smartEvseUrl,
		writeAPI: client.WriteAPIBlocking(influxdbOrg, influxdbBucket),
		tags:     tagMap,
	}, nil
}

func (s *Submitter) Work() error {
	settings, err := smartevse.FromHTTP(s.url)
	if err != nil {
		return fmt.Errorf("error fetching http: %v", err)
	}

	p := NewCfPoint("smartevse", s.tags)
	smartevse.VisitNumericFields(*settings, p.visit)

	slog.Debug("submitting point", "point", p.Point)
	return s.writeAPI.WritePoint(s.context, p.Point)
}
