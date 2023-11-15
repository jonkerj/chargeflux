package submitter

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/jonkerj/chargeflux/pkg/smartevse"
)

type Submitter struct {
	context  context.Context
	url      string
	writeAPI api.WriteAPIBlocking
	tags     map[string]string
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

	p := influxdb2.NewPointWithMeasurement("smartevse")
	for tag, value := range s.tags {
		p.AddTag(tag, value)
	}
	settings.AddFieldsToPoint(p)

	slog.Debug("submitting point", "point", p)
	s.writeAPI.WritePoint(s.context, p)

	return nil
}
