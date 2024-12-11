if [ $# -ne 3 ]; then
  echo "Usage: $0 <username> <password> <database_name>"
  exit 1
fi

source tcdt_logo.sh
echo -e "${GREEN}Starting the database startup process...${NC}"

username=$1
password=$2
database_name=$3

DB_EXISTS=$(mysql -u "$USERNAME" -p"$PASSWORD" -e "SHOW DATABASES LIKE '$DB_NAME';" | grep -w "$DB_NAME")

if [ "$DB_EXISTS" ]; then
  echo "The database '$DB_NAME' already exists."
else
  echo "The database '$DB_NAME' does not exist."
  mysql -u $username -p$password -e "CREATE DATABASE $database_name;"
  mysql "GRANT SELECT, INSERT, UPDATE, CREATE , DELETE ON $database_name.* TO 'koosha'@'localhost';
"
  echo -e "${GREEN}Database created successfully!${NC}"
fi

echo -e "${GREEN}Starting the table creation process ...${NC}"
mysql -u $username -p$password $database_name < create_tables.sql
echo -e "${GREEN}Database startup process completed!${NC}"


