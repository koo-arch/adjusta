env "local" {
  src = "ent://internal/infrastructure/ent/schema"
  dev = "docker://postgres/18/dev"

  migration {
    dir = "file://migrations"
  }
}
