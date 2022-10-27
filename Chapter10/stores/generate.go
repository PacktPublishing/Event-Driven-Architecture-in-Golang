package stores

//go:generate buf generate

//go:generate mockery --quiet --dir ./storespb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c storesclient -m storesclient/models --with-flatten=remove-unused
