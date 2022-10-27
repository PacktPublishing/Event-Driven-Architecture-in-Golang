resource null_resource init_db {
  provisioner "local-exec" {
    command = "psql --file sql/init_db.psql ${local.db_conn}/postgres"
  }
}
