#!/usr/bin/env groovy
@Library('tools') _

kanikoBuild {
    registry '136813947591.dkr.ecr.us-east-1.amazonaws.com'
    image "${getEcrEnv()}/app/gitlab/gitlab-exporter"
    tags "v${env.BUILD_ID}", 'latest'
    ecrDeleteLatest true
}