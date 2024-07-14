//go:build !windows

package main

import (
	"context"
	"encoding/json"
	"github.com/splicer3/grafana-gt7/pkg/gt7"
	"net"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/grafana-plugin-sdk-go/live"
)

var (
	_ backend.QueryDataHandler      = (*GT7TelemetryDatasource)(nil)
	_ backend.CheckHealthHandler    = (*GT7TelemetryDatasource)(nil)
	_ backend.StreamHandler         = (*GT7TelemetryDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*GT7TelemetryDatasource)(nil)
)

type Options struct {
	PlaystationIP string `json:"playstationIP"`
}

func getDatasourceSettings(s backend.DataSourceInstanceSettings) (*Options, error) {
	settings := &Options{}

	if err := json.Unmarshal(s.JSONData, settings); err != nil {
		return nil, err
	}

	return settings, nil
}

// NewGT7TelemetryDatasource creates a new datasource instance.
func NewGT7TelemetryDatasource(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	settings, err := getDatasourceSettings(s)
	if err != nil {
		return nil, err
	}

	return &GT7TelemetryDatasource{
		playstationIP: settings.PlaystationIP,
		streamConn:    nil,
		heartbeatConn: nil,
	}, nil
}

// GT7TelemetryDatasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type GT7TelemetryDatasource struct {
	playstationIP string
	streamConn    *net.UDPConn
	heartbeatConn *net.UDPConn
}

func (d *GT7TelemetryDatasource) Dispose() {
	// Clean up datasource instance resources.
	if d.streamConn != nil {
		d.streamConn.Close()
		d.streamConn = nil
	}

	if d.heartbeatConn != nil {
		d.heartbeatConn.Close()
		d.heartbeatConn = nil
	}
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *GT7TelemetryDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData called", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	WithStreaming bool   `json:"withStreaming"`
	Telemetry     string `json:"telemetry"`
}

func (d *GT7TelemetryDatasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}

	// create data frame response.
	frame := data.NewFrame("response")

	// add fields.
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
		data.NewField("values", nil, []float32{0, 0}),
	)

	// If query called with streaming on then return a channel
	// to subscribe on a client-side and consume updates from a plugin.
	// Feel free to remove this if you don't need streaming for your datasource.
	//streamPath := qm.Telemetry
	//if streamPath == "" {
	//	streamPath = "stream"
	//}
	streamPath := "dirt"
	if qm.WithStreaming {
		channel := live.Channel{
			Scope:     live.ScopeDatasource,
			Namespace: pCtx.DataSourceInstanceSettings.UID,
			Path:      streamPath,
		}
		frame.SetMeta(&data.FrameMeta{Channel: channel.String()})
	}

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *GT7TelemetryDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	log.DefaultLogger.Info("CheckHealth called", "request", req)

	var status = backend.HealthStatusOk
	var message = "Data source is working"

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

// SubscribeStream is called when a client wants to connect to a stream. This callback
// allows sending the first message.
func (d *GT7TelemetryDatasource) SubscribeStream(_ context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	log.DefaultLogger.Info("SubscribeStream called", "request", req)

	status := backend.SubscribeStreamStatusOK
	return &backend.SubscribeStreamResponse{
		Status: status,
	}, nil
}

// RunStream is called once for any open channel. Results are shared with everyone
// subscribed to the same channel.
func (d *GT7TelemetryDatasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	log.DefaultLogger.Info("RunStream called", "request", req)

	// Check if any existing stream exists and close it.
	if d.streamConn != nil {
		d.streamConn.Close()
		d.streamConn = nil
	}

	if d.heartbeatConn != nil {
		d.heartbeatConn.Close()
		d.heartbeatConn = nil
	}

	heartbeatConnChan := make(chan *net.UDPConn)
	streamConnChan := make(chan *net.UDPConn)
	gt7TelemetryChan := make(chan gt7.TelemetryFrame)
	gt7TelemetryErrorChan := make(chan error)

	if req.Path == "gt7" {
		go gt7.RunTelemetryServer(d.playstationIP, gt7TelemetryChan, gt7TelemetryErrorChan, heartbeatConnChan, streamConnChan)
	}

	lastTimeSent := time.Now()

	// Stream data frames periodically till stream closed by Grafana.
	for {
		select {
		case <-ctx.Done():
			log.DefaultLogger.Info("Context done, finish streaming", "path", req.Path)
			if d.streamConn != nil {
				d.streamConn.Close()
			}
			if d.heartbeatConn != nil {
				d.heartbeatConn.Close()
				d.heartbeatConn = nil
			}
			return nil

		case telemetryFrame := <-gt7TelemetryChan:
			if time.Now().Before(lastTimeSent.Add(time.Second / 60)) {
				// Drop frame
				continue
			}

			frame := gt7.TelemetryToDataFrame(telemetryFrame)
			lastTimeSent = time.Now()
			err := sender.SendFrame(frame, data.IncludeAll)
			if err != nil {
				log.DefaultLogger.Error("Error sending frame", "error", err)
				continue
			}

		case hbConn := <-heartbeatConnChan:
			d.heartbeatConn = hbConn

		case strConn := <-streamConnChan:
			d.streamConn = strConn

		case err := <-gt7TelemetryErrorChan:
			log.DefaultLogger.Error("Error from telemetry server", "error", err)
			return err
		}
	}
}

// PublishStream is called when a client sends a message to the stream.
func (d *GT7TelemetryDatasource) PublishStream(_ context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	log.DefaultLogger.Info("PublishStream called", "request", req)

	// Do not allow publishing at all.
	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied,
	}, nil
}
