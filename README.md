# Go Monitoring Application

This application is developed for continuously checking website availability and exposes metrics for Prometheus. Logs failures are stored to local files for simplicity.

## Prerequisites

- Go 1.25 or later
- Kubernetes cluster (optional, for monitoring/alerting)
- kubectl and Helm 3.x (optional, for monitoring/alerting)

## Start the application

### 1. Install Dependencies

```bash
go mod download
```

### 2. Configure the Target

Edit `config.yaml` to set the website you want to monitor:

### 3. Run the Application

```bash
go run main.go
```

The app starts on port 8080.  visit localhost:8080/endpoint_name to access the endpoints.

#### API Endpoints

The application exposes three HTTP endpoints:

- **`GET /metrics`** - Metrics with timestamps and information:
Total requests made, Successful requests, Overall availability percentage

- **`GET /availability`** - Current availability (by percentage)

- **`GET /logs`** - Recent failure logs (last 10 entries by default), 
Older logs/ remaining logs before shutdown are written to `logs/monitoring.log`

### 4. Deploy Monitoring Stack (Optional)

Deploy Prometheus with Helm to visualize metrics:

```bash
cd helm

helm dependency update monitoring-stack

# Deploy to Kubernetes
helm upgrade --install monitoring-stack ./monitoring-stack --namespace monitoring --create-namespace

# Check deployment status
kubectl get pods -n monitoring
```

Prometheus gets the target metrics every 10 seconds (configurable)

**Access Prometheus UI:**
```bash
# Open in browser (NodePort default: 30900)
http://localhost:30900
```

**Configured Alerts:**
- **AvailabilityDrop** - Warns when availability drops below 95% for 1 minutes

## Cleanup

**Stop the application:** Press `Ctrl+C` 

**Remove Helm deployment:**
```bash
helm uninstall monitoring-stack -n monitoring
kubectl delete namespace monitoring
```

## Testing

Run the test suite:
```bash
go test ./test/...
```

or run the test separately, e.g.:
```bash
go test ./test/config_test.go
```

## Project Structure

```
api/
  handler.go          # HTTP handlers and API routes
config/
  config.go           # Configuration management
services/
  metrics.go          # Prometheus metrics definitions
  monitor.go          # Website monitoring and log management
helm/
  monitoring-stack/   # Prometheus Helm chart
    Chart.yaml        # helm dependencies
    values.yaml       # prometheus configuration(resources,scrap,rules)
logs/                 # failure logs (auto-created)
config.yaml           # monitored website settings
go.mod                # Go dependencies
main.go               # main function / entry point