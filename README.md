[![Build Status](https://jenkins.public.jamf.build/buildStatus/icon?job=Devops%2FBuild+Infrastructure%2FGitlab%2Fgitlab-exporter-github%2Fmaster)](https://jenkins.jamf.build/job/Devops/job/Build%20Infrastructure/job/Gitlab/job/gitlab-exporter-github/job/master/)

# Gitlab license exporter
Exposes License expiration date and active users from the Gitlab API, to a Prometheus compatible endpoint.

## Configuration
This exporter is setup to take input from environment variables:

### Required
* `TOKEN`: Admin token

### Optional
* `URL`: Gitlab url example: `https://gitlab.domain.com` (by default will use k8s service `gitlab-web`)

## Build and run
### Manually
```
go get
go build gitlabgoexporter.go
export TOKEN=token123token
export URL=https://gitlab.domain.com
./gitlabgoexporter.go
```
Visit http://localhost:2222


### Docker
Build a docker image:
`docker build -t <image-name> .`

Run:
* Custom URL:
	`docker image --env TOKEN=token123token --env URL=https://gitlab.domain.com <image-name>`

* Kubernetes Gitlab-Web service:
	`docker image --env TOKEN=token123token <image-name>`


### Kubernetes
```
apiVersion: v1
kind: Secret
metadata:
  name: gitlab-token
  namespace: {{ NAMESPACE }}
  labels:
    app: gitlab-exporter
type: Opaque
data:
  token: {{ TOKEN | b64encode }}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gitlab-exporter
  namespace: {{ NAMESPACE }}
  labels:
    app: gitlab-exporter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gitlab-exporter
  template:
    metadata:
      labels:
        app: gitlab-exporter
    spec:
      containers:
      - name: gitlab-exporter
        image: {{{ image-name }}}
        ports:
        - containerPort: 2222
        env:
        - name: TOKEN
          valueFrom:
            secretKeyRef:
              name: gitlab-token
              key: token
        resources:
          limits:
            cpu: 100m
            memory: 128Mi
          requests:
            cpu: 50m
            memory: 64Mi
---
apiVersion: v1
kind: Service
metadata:
  name: gitlab-exporter-svc
  namespace: {{ NAMESPACE }}
  labels:
    app: gitlab-exporter
spec:
  selector: 
    app: gitlab-exporter
  ports:
    - name: metrics
      port: 8080
      targetPort: 2222
---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: gitlab-export-metrics
  namespace: {{ NAMESPACE }}
spec:
  selector:
    matchLabels:
      app: gitlab-exporter
  endpoints:
  - port: metrics
    path: /
    interval: 30s
```

## Metrics
Metrics will be available on port 2222 by default

## Collectors
```
# HELP gitlab_active_users Gitlab active users
# HELP gitlab_license_expires_at Gitlab expiration day
# HELP gitlab_scrape_success Gitlab go exporter scrape status when try to read the API
# HELP gitlab_user_limit Users allowed by license
```