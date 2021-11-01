#!/bin/bash -e
shopt -s expand_aliases
kim_support=false
docker_support=true
exit_code=0

if command -v docker &> /dev/null
then
    if command -v kim &> /dev/null
    then
        kim_support=true
    fi
fi

[ -z "$USER" ] && echo "env variable USER must be set" && exit 1;

go get -u -d github.com/karlmutch/duat/cmd/semver
version=`$GOPATH/bin/semver`

for dir in cmd/*/ ; do
    base="${dir%%\/}"
    base="${base##*/}"
    if [ "$base" == "cli-experiment" ] ; then
        continue
    fi
    if [ "$base" == "cli-downstream" ] ; then
        continue
    fi
    echo "$dir"
    cd $dir
    if [ "$kim_support" = true ] ; then
#        nerdctl --namespace k8s.io build -t platform-services/$base:$version .
        kim build -t platform-services/$base:$version .
    else
        docker build -t platform-services/$base:$version .
    fi
    cd -
done

exit 0

if [ "$kube_cfg" == "false"]
then
./push.sh
fi
