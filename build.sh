#!/bin/bash -e
shopt -s expand_aliases
nerdctl_support=false
docker_support=true
exit_code=0

if ! command -v docker &> /dev/null
then
    docker_support=false
    if command -v nerdctl &> /dev/null
    then
        nerdctl_support=true
    fi
fi

[ -z "$USER" ] && echo "env variable USER must be set" && exit 1;

if [ "$docker_support" == "false" ]; then
    nerdctl build -t platform-services:latest --build-arg USER=$USER --build-arg USER_ID=`id -u $USER` --build-arg USER_GROUP_ID=`id -g $USER` .
else
    docker build -t platform-services:latest --build-arg USER=$USER --build-arg USER_ID=`id -u $USER` --build-arg USER_GROUP_ID=`id -g $USER` .
fi
docker_name=`petname`

if [ "$docker_support" == "false" ]; then
    nerdctl run --name $docker_name -e GITHUB_TOKEN=$GITHUB_TOKEN -v $GOPATH:/project platform-services 
else 
    docker run --name $docker_name -e GITHUB_TOKEN=$GITHUB_TOKEN -v $GOPATH:/project platform-services 
fi

exit_code=`docker inspect $docker_name --format='{{.State.ExitCode}}'`

if [ $exit_code -ne 0 ]; then
    echo "Error" $exit_code
    exit $exit_code
fi

echo "Build Done" ; docker container prune -f

./images.sh

exit 0

if [ "$kube_cfg" == "false"]
then
./push.sh
fi
