package customers

//go:generate buf generate

//go:generate mockery --quiet --dir ./customerspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c customersclient -m customersclient/models --with-flatten=remove-unused
