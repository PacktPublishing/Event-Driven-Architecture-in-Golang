$global:coderoot = (Resolve-Path $PWD\..\..).ToString()

docker build -t deploytools:latest .

Function global:deploytools {
    docker run --rm -it `
    -v "//var/run/docker.sock://var/run/docker.sock" `
    -v $env:userprofile\.aws:/root/.aws `
    -v $env:userprofile\.kube:/root/.kube `
    -v ${coderoot}:/mallbots `
    -v ${PWD}:/mallbots/deployment/.current `
    -w /mallbots/deployment/.current deploytools `
    $args
}

echo "---"
echo ""
echo "Usage: deploytools <cmd [parameters]>"
