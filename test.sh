#!/usr/bin/env bash

function createDatabase()
{
  echo "createDatabase"
  docker run \
    -d \
    -p 5432:5432 \
    -e POSTGRES_PASSWORD=tester \
    -e POSTGRES_USER=tester \
    -e POSTGRES_DB=postgres \
    --name tester_postgres \
    postgres:11.5
}

function injectStructure()
{
  echo "injectStructure"
  docker exec \
    -e PGPASSWORD=tester tester_postgres psql \
    -U tester \
    -d postgres \
    -c "CREATE TABLE "public"."agent" ("id" uuid, "name" varchar(200), "key" uuid, "secret" uuid, "company_id" uuid, PRIMARY KEY ("id"));"
}

function testCode()
{
  echo "Test Code"
  DB_DATABASE=postgres DB_TABLE=agent DB_HOSTNAME=0.0.0.0 DB_PORT=5432 DB_USERNAME=tester DB_PASSWORD=tester go test ./...
  echo "-------"
  DB_DATABASE=postgres DB_TABLE=agent DB_HOSTNAME=0.0.0.0 DB_PORT=5432 DB_USERNAME=tester DB_PASSWORD=tester go test ./... -bench=. -run=$$$
}

if [[ ! -z ${1} ]] || [[ "${1}" != "" ]]; then
  ${1}
else
  createDatabase
  sleep 5
  injectStructure
  testCode
fi
