# apigeehello

## Install

    git clone https://github.com/udhos/apigeehello
    cd apigeehello
    CGO_ENABLED=0 go install ./apiserver

## Run

    apiserver

## Test

    curl -X POST -d '{foobar}' http://evertonmarques-eval-test.apigee.net/apigee-hello-world/echo -H "accept: application/json"
