#!/usr/bin/env groovy

@Library('tools@customargs') _

def branch = (env.CHANGE_BRANCH ?: env.BRANCH_NAME).replaceAll('/', '-').toLowerCase()
dockerBuild image: "docker.jamf.build/devops/gitlab-exporter-${branch}",
            tag: "v${env.BUILD_ID}",
            dockerfile: "exporter",
            login: true