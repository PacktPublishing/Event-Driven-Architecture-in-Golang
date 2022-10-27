package am

//go:generate buf generate

//go:generate mockery --quiet --name ".*(Subscriber|Publisher|Handler)$"  --inpackage --case underscore
