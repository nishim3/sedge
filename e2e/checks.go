package e2e

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/NethermindEth/sedge/internal/pkg/services"
	"github.com/cenkalti/backoff"
	"github.com/docker/docker/client"
	gapi "github.com/grafana/grafana-api-golang-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// checkMonitoringStackDir checks that the monitoring stack directory exists and contains the docker-compose file
func checkMonitoringStackDir(t *testing.T) {
	t.Logf("Checking monitoring stack directory")
	// Check monitoring folder exists
	dataDir, err := dataDirPath()
	if err != nil {
		t.Fatal(err)
	}
	monitoringDir := filepath.Join(dataDir, "monitoring")
	assert.DirExists(t, monitoringDir)

	// Check monitoring docker-compose file exists
	assert.FileExists(t, filepath.Join(monitoringDir, "docker-compose.yml"))
}

// checkMonitoringStackNotInstalled checks that the monitoring stack directory exists but is not installed
func checkMonitoringStackNotInstalled(t *testing.T) {
	t.Logf("Checking monitoring stack directory")
	// Check monitoring folder exists
	dataDir, err := dataDirPath()
	if err != nil {
		t.Fatal(err)
	}
	monitoringDir := filepath.Join(dataDir, "monitoring")
	assert.DirExists(t, monitoringDir)

	// Check monitoring docker-compose file does not exists
	assert.NoFileExists(t, filepath.Join(monitoringDir, "docker-compose.yml"))
}

// checkMonitoringStackContainers checks that the monitoring stack containers are running
func checkMonitoringStackContainers(t *testing.T) {
	t.Logf("Checking monitoring stack containers")
	checkContainerRunning(t, "sedge_grafana", "sedge_prometheus", "sedge_node_exporter")
}

// checkContainerRunning checks that the given containers are running
func checkContainerRunning(t *testing.T, containerNames ...string) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		t.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	dockerServiceManager := services.NewDockerServiceManager(cli)

	for _, containerName := range containerNames {
		t.Logf("Checking %s container is running", containerName)
		isRunning, err := dockerServiceManager.IsRunning(containerName)
		require.NoError(t, err)
		assert.True(t, isRunning, "%s container should be running", containerName)
	}
}

// checkPrometheusTargetsUp checks that the prometheus targets are up
func checkPrometheusTargetsUp(t *testing.T, targets ...string) {
	var (
		tries       int           = 0
		timeOut     time.Duration = 30 * time.Second
		promTargets *PrometheusTargetsResponse
		err         error
	)
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	b := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)
	err = backoff.Retry(func() error {
		tries++
		logPrefix := fmt.Sprintf("checkPrometheusTargetsUp (%d)", tries)
		promTargets, err = prometheusTargets(t)
		if err != nil {
			return logAndPipeError(t, logPrefix, err)
		}
		if promTargets.Status != "success" {
			return logAndPipeError(t, logPrefix, fmt.Errorf("expected status success, got %s", promTargets.Status))
		}
		if len(promTargets.Data.ActiveTargets) != len(targets) {
			return logAndPipeError(t, logPrefix, fmt.Errorf("expected %d targets, got %d", len(targets), len(promTargets.Data.ActiveTargets)))
		}
		for i, target := range promTargets.Data.ActiveTargets {
			var labels []string
			for label := range target.Labels {
				labels = append(labels, label)
			}
			if !slices.Contains(labels, "instance") {
				return logAndPipeError(t, logPrefix, fmt.Errorf("target %d does not have instance label", i))
			}
			instanceLabel := target.Labels["instance"]
			if !slices.Contains(targets, instanceLabel) {
				return logAndPipeError(t, logPrefix, fmt.Errorf("target %d instance label is not expected", i))
			}
			if target.Health == "unknown" {
				return logAndPipeError(t, logPrefix, fmt.Errorf("target %d health is unknown", i))
			}
		}
		return nil
	}, b)
	assert.NoError(t, err, `targets "%s" should be up, but after %d tries they are not`, targets, tries)
}

// checkPrometheusHealth checks that the prometheus health is ok
func checkGrafanaHealth(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tries := 0
	b := backoff.WithContext(backoff.NewConstantBackOff(time.Second), ctx)
	err := backoff.Retry(func() error {
		logPrefix := fmt.Sprintf("checkGrafanaHealth (%d)", tries+1)
		tries++
		// Check Grafana health
		gClient, err := gapi.New("http://localhost:3000", gapi.Config{
			BasicAuth: url.UserPassword("admin", "admin"),
		})
		if err != nil {
			return logAndPipeError(t, logPrefix, err)
		}
		healthResponse, err := gClient.Health()
		if err != nil {
			return logAndPipeError(t, logPrefix, err)
		}
		if healthResponse.Database != "ok" {
			return logAndPipeError(t, logPrefix, fmt.Errorf("expected database ok, got %s", healthResponse.Database))
		}
		return nil
	}, b)
	assert.NoError(t, err, "Grafana should be ok, but it is not")
}

// checkMonitoringStackContainersNotRunning checks that the monitoring stack containers are not running
func checkMonitoringStackContainersNotRunning(t *testing.T) {
	t.Logf("Checking monitoring stack containers are not running")
	checkContainerNotExisting(t, "sedge_grafana", "sedge_prometheus", "sedge_node_exporter")
}

// checkContainerNotExisting checks that the given containers are not existing
func checkContainerNotExisting(t *testing.T, containerNames ...string) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		t.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	dockerServiceManager := services.NewDockerServiceManager(cli)

	for _, containerName := range containerNames {
		t.Logf("Checking %s container is not existing", containerName)
		_, err := dockerServiceManager.ContainerStatus(containerName)
		assert.Error(t, err)
	}
}
