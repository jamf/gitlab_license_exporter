#!/usr/bin/env groovy

@Library(['tools', 'build-infrastructure']) _

def registryEnv = isProd() ? 'prod' : 'staging'

dockerBuild {
    agent 'aws-agent'
	registry '136813947591.dkr.ecr.us-east-1.amazonaws.com'
    image "${registryEnv}/app/gitlab/gitlab-exporter"
    tag "v${env.BUILD_ID}" 
    login true
}