#!/usr/bin/env bash

# Ensure the database container is online and usable
# echo "Waiting for database..."
until docker exec -i dullahan.db mysql -u dullahan -pdullahan123 -D dullahan -e "SELECT 1" &> /dev/null
# EnablePostgreSQL: remove the line above, uncomment the following
# until docker exec -i dullahan.db psql -h localhost -U dullahan  -d dullahan -c "SELECT 1" &> /dev/null
do
  # printf "."
  sleep 1
done
