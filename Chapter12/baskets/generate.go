package baskets

//go:generate buf generate

//go:generate mockery --quiet --dir ./basketspb -r --all --inpackage --case underscore
//go:generate mockery --quiet --dir ./internal -r --all --inpackage --case underscore

//go:generate swagger generate client -q -f ./internal/rest/api.swagger.json -c basketsclient -m basketsclient/models --with-flatten=remove-unused
