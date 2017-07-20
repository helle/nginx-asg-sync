package main

import "testing"

var validYaml = []byte(`region: us-west-2
upstream_conf_endpoint: http://127.0.0.1:8080/upstream_conf
status_endpoint: http://127.0.0.1:8080/status
sync_interval_in_seconds: 5
upstreams:
  - name: backend1
    autoscaling_group: backend-group
    port: 80
    kind: http
  - name: backend2
    autoscaling_group: backend-group
    port: 80
    kind: http
    max_conns: 13
`)

type testInput struct {
	cfg *config
	msg string
}

func getValidConfig() *config {
	upstreams := []upstream{
		upstream{
			Name:             "backend1",
			AutoscalingGroup: "backend-group",
			Port:             80,
			Kind:             "http",
		},
	}
	cfg := config{
		Region:                "us-west-2",
		UpstreamConfEndpont:   "http://127.0.0.1:8080/upstream_conf",
		StatusEndpoint:        "http://127.0.0.1:8080/status",
		SyncIntervalInSeconds: 1,
		Upstreams:             upstreams,
	}

	return &cfg
}

func getInvalidConfigInput() []*testInput {
	var input []*testInput

	invalidRegionCfg := getValidConfig()
	invalidRegionCfg.Region = ""
	input = append(input, &testInput{invalidRegionCfg, "invalid region"})

	invalidUpstreamConfEndponCfg := getValidConfig()
	invalidUpstreamConfEndponCfg.UpstreamConfEndpont = ""
	input = append(input, &testInput{invalidUpstreamConfEndponCfg, "invalid upstream_conf_endpoint"})

	invalidStatusEndpointCfg := getValidConfig()
	invalidStatusEndpointCfg.StatusEndpoint = ""
	input = append(input, &testInput{invalidStatusEndpointCfg, "invalid status_endpoint"})

	invalidSyncIntervalInSecondsCfg := getValidConfig()
	invalidSyncIntervalInSecondsCfg.SyncIntervalInSeconds = 0
	input = append(input, &testInput{invalidSyncIntervalInSecondsCfg, "invalid sync_interval_in_seconds"})

	invalidMissingUpstreamsCfg := getValidConfig()
	invalidMissingUpstreamsCfg.Upstreams = nil
	input = append(input, &testInput{invalidMissingUpstreamsCfg, "no upstreams"})

	invalidUpstreamNameCfg := getValidConfig()
	invalidUpstreamNameCfg.Upstreams[0].Name = ""
	input = append(input, &testInput{invalidUpstreamNameCfg, "invalid name of the upstream"})

	invalidUpstreamAutoscalingGroupCfg := getValidConfig()
	invalidUpstreamAutoscalingGroupCfg.Upstreams[0].AutoscalingGroup = ""
	input = append(input, &testInput{invalidUpstreamAutoscalingGroupCfg, "invalid autoscaling_group of the upstream"})

	invalidUpstreamPortCfg := getValidConfig()
	invalidUpstreamPortCfg.Upstreams[0].Port = 0
	input = append(input, &testInput{invalidUpstreamPortCfg, "invalid port of the upstream"})

	invalidUpstreamKindCfg := getValidConfig()
	invalidUpstreamKindCfg.Upstreams[0].Kind = ""
	input = append(input, &testInput{invalidUpstreamKindCfg, "invalid kind of the upstream"})

	return input
}

func TestExtraParams(t *testing.T) {
	ups := upstream{
		Name:             "backend1",
		AutoscalingGroup: "backend-group",
		Port:             80,
		Kind:             "http",
		MaxConns:         3,
	}

	extraParams := makeExtraParams(&ups)

	if extraParams != "&max_conns=3" {
		t.Errorf("makeExtraParams() generated wrong params: %v", extraParams)
	}
}

func TestUnmarshalConfig(t *testing.T) {
	cfg, err := unmarshalConfig(validYaml)
	if err != nil {
		t.Errorf("unmarshalConfig() failed for the valid yaml: %v", err)
	}

	if cfg.Upstreams[0].MaxConns != 0 {
		t.Errorf("unmarshalConfig() failed to read maxconns: %d", cfg.Upstreams[0].MaxConns)
	}

	if cfg.Upstreams[1].MaxConns != 13 {
		t.Errorf("unmarshalConfig() failed to read maxconns: %d", cfg.Upstreams[1].MaxConns)
	}
}

func TestValidateConfigNotValid(t *testing.T) {
	input := getInvalidConfigInput()

	for _, item := range input {
		err := validateConfig(item.cfg)
		if err == nil {
			t.Errorf("validateConfig() didn't fail for the invalid config file with %v", item.msg)
		}
	}
}

func TestValidateConfigValid(t *testing.T) {
	cfg := getValidConfig()

	err := validateConfig(cfg)
	if err != nil {
		t.Errorf("validateConfig() failed for the valid config: %v", err)
	}
}
