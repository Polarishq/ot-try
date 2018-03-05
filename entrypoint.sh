#!/usr/bin/env sh

term_handler() {
    echo "Caught SIGTERM signal!"
    kill -s SIGTERM "$child"
}

int_handler() {
    echo "Caught SIGINT signal!"
    kill -s SIGINT "$child"
}

trap term_handler SIGTERM
trap int_handler SIGINT

SERVER_ARGS="--host=0.0.0.0 --port=80
	--tls-host=0.0.0.0 --tls-port=443
	--tls-certificate=/tmp/ssl/server.crt --tls-key=/tmp/ssl/server.key
	--static=/var/www/static"


if [ "$SERVER_COVERAGE" == "true" ]; then
    echo Starting ot-try server with coverage
    COVERPROFILE=$ENVIRONMENT_NAME.coverprofile
    /go/bin/ot-try-coverage -test.run "^TestRunMain$" \
        -test.coverprofile=$COVERPROFILE \
        $SERVER_ARGS

    url="$ARTIFACTORY_URL/prod/ot-try/$CI_COMMIT_ID/$ENVIRONMENT_NAME.coverprofile"
    echo Uploading $COVERPROFILE to $url
    curl -u"$ARTIFACTORY_USERNAME:$ARTIFACTORY_PASSWORD" -T $COVERPROFILE $url
else
    echo Starting ot-try server in production mode
    /go/bin/ot-try $SERVER_ARGS &
    child=$!
    wait "$child"
fi
