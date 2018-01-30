# Introduction

The experiment server is used to persist experiment details and to record changes to the state of experiments.  Items included within an experiment include layer definitions and meta-data items.

The experiment server offers a gRPC API that can be accessed using a machine-to-machine or human-to-machine (HCI) interface.  The HCI interface can be interacted with using the grpc_cli tool provided with the gRPC toolkita  More information about grpc_cli can be found at, https://github.com/grpc/grpc/blob/master/doc/command_line_tool.md.

# Experiment Database

The experiment server makes use of a Postgres DB.  The installation process is specific to AWS Aurora.  To begin installation you will need to create a Postgres Aurora RDS instance.  Use values of your choosing for the DB name and user/password combinations.

Parameters that impact deployment of your Aurora instance include, RDS Endpoint, DB Name, user name, and the password.

The command to install the postgres schema into your DB instance will appear similar to the following:

```
PGUSER=pl PGHOST=dev-platform.cluster-cff2uhtd2jzh.us-west-2.rds.amazonaws.com PGDATABASE=platform psql -f sql/platform.sql -d postgres
```

## Installation

Now go to the experimentsrv.yaml file and change the Egress rule to point at your endpoint.  The deployment spec, PGHOST and PGDATABASE should also be modified to the endpoint and the DB Name respectively.

You should now create or edit the secrets.yaml file ready for deployment with the user name and the password.

# Service Authentication

Service authetication is explained within the top level README.md file for the github.com/SentientTechnologies/platform-services repository.  All calls into the experiment service must contain metadata for the autorization brearer token and have the all:experiments claim in order to be accepted.

# Manually exercising the server

```
export AUTH0_DOMAIN=sentientai.auth0.com
export AUTH0_TOKEN=$(curl -s --request POST --url 'https://sentientai.auth0.com/oauth/token' --header 'content-type: application/json' --data '{ "client_id":"71eLNu9Bw1rgfYz9PA2gZ4Ji7ujm3Uwj", "client_secret": "AifXD19Y1EKhAKoSqI5r9NWCdJJfyN0x-OywIumSd9hqq_QJr-XlbC7b65rwMjms", "audience": "http://api.sentient.ai/experimentsrv", "grant_type": "http://auth0.com/oauth/grant-type/password-realm", "username": "karlmutch@gmail.com", "password": "Passw0rd!", "scope": "all:experiments", "realm": "Username-Password-Authentication" }' | jq -r '"\(.access_token)"')
/tmp/grpc_cli call localhost:30001 ai.sentient.experiment.Service.Get "id: ''" --metadata authorization:"Bearer $AUTH0_TOKEN"
connecting to localhost:30001
Sending client initial metadata:
authorization : ...
Rpc failed with status code 2, error message: selecting an experiment requires either the DB id or the experiment unique id to be specified stack="[db.go:533 server.go:42 experimentsrv.pb.go:375 auth.go:88 experimentsrv.pb.go:377 server.go:900 server.go:1122 server.go:617]"

/tmp/grpc_cli call localhost:30001 ai.sentient.experiment.Service.Get "id: 't'" --metadata authorization:"$AUTH"
connecting to localhost:30001
Sending client initial metadata:
authorization : ...
Rpc failed with status code 2, error message: no matching experiments found for caller specified input parameters
```

If you wish to exercise the server while it is deployed into an Istio orchestrated cluster then you should kubectl exec into the istio-proxy to get localized access to the service.  The application server container runs only Alpine so for a full set of tools such as curl and jq using the istio-proxy is a better option, with some minor additions. The Istio container will need the following commands run in order to activate the needed features for local testing:

```
# Copy the grpc CLI tool into the running container
kubectl cp `which grpc_cli` experiments-v1-bc46b5d68-bcdkv:/tmp/grpc_cli -c istio-proxy
kubectl exec -it experiments-v1-bc46b5d68-bcdkv -c istio-proxy /bin/bash
sudo apt-get update
sudo apt-get install -y libgflags2v5 ca-certificates jq
export AUTH0_TOKEN=$(curl -s --request POST --url 'https://sentientai.auth0.com/oauth/token' --header 'content-type: application/json' --data '{ "client
_id":"71eLNu9Bw1rgfYz9PA2gZ4Ji7ujm3Uwj", "client_secret": "AifXD19Y1EKhAKoSqI5r9NWCdJJfyN0x-OywIumSd9hqq_QJr-XlbC7b65rwMjms", "audience": "http://api.sentient.
ai/experimentsrv", "grant_type": "http://auth0.com/oauth/grant-type/password-realm", "username": "karlmutch@gmail.com", "password": "Passw0rd!", "scope": "all:
experiments", "realm": "Username-Password-Authentication" }' | jq -r '"\(.access_token)"')
/tmp/grpc_cli call 100.96.1.14:30001 ai.sentient.experiment.Service.Get "id: 't'"  --metadata authorization:"Bearer $AUTH0_TOKEN"
```

When Istio is used without a Load balancer the IP to be used can be determined by using the following command:

```
kubectl -n istio-system get po -l istio=ingress -o jsonpath='{.items[0].status.hostIP}'
```

# Using IP Port addresses for serving grpc

The grpc listener is configured so that both an IPv4 and IPv6 adapter is attemped individually.  Tgis is because the behavior on some OS's such as Alpine is to not use the 'tcp' protocol with an empty address for both IPv4 and Ipv6, this is the case on Alpine.  Some OS's such as Ubuntu do open both IPv6 and IPv4 in these cases.  As such when using Ubuntu you will always need to specify a single IP, ":30001" for example, to override the default.  The default is designed for use with Alpine containerized deployments such as those in production.

# Functional testing

The server provides some testing for the DB and core functionality of the server.  In order to run this you can use the go test command and point at the relevant go package directories from which you wish to run the tests, for example to run the experiment DB and server tests you could use commands like the following:

```
cd cmd/experimentsrv
export AUTH0_DOMAIN=sentientai.auth0.com
export AUTH0_TOKEN=$(curl -s --request POST --url 'https://sentientai.auth0.com/oauth/token' --header 'content-type: application/json' --data '{ "client_id":"71eLNu9Bw1rgfYz9PA2gZ4Ji7ujm3Uwj", "client_secret": "AifXD19Y1EKhAKoSqI5r9NWCdJJfyN0x-OywIumSd9hqq_QJr-XlbC7b65rwMjms", "audience": "http://api.sentient.ai/experimentsrv", "grant_type": "http://auth0.com/oauth/grant-type/password-realm", "username": "karlmutch@gmail.com", "password": "Passw0rd!", "scope": "all:experiments", "realm": "Username-Password-Authentication" }' | jq -r '"\(.access_token)"')
LOGXI=*=TRC PGUSER=pl PGHOST=dev-platform.cluster-cff2uhtd2jzh.us-west-2.rds.amazonaws.com PGDATABASE=platform go test -v . -ip-port ":30001"
```