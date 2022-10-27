package ddd

//go:generate mockery --quiet --name ".*(Aggregate|Entity|Subscriber|Publisher|Handler)$"  --inpackage --case underscore
