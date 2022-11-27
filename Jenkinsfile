#!/usr/bin/env groovy
@Library('tools') _

kanikoBuild {
    registry ${env.ECRREPO}
    image "${getEcrEnv()}/app/gitlab/gitlab-exporter"
    tags "v${env.BUILD_ID}", 'latest'
    ecrDeleteLatest true
}