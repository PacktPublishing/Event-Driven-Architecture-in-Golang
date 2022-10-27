package depot

//go:generate buf generate

//go:generate mockery --quiet --dir ./depotpb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c depotclient -m depotclient/models --with-flatten=remove-unused
