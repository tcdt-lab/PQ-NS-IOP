if [ $# -ne 2 ]; then
  echo "Usage: $0 <username> <database_name>"
  exit 1
fi

source tcdt_logo.sh
echo -e "${GREEN}Starting the database startup process...${NC}"

username=$1

database_name=$2

DB_EXISTS=$(mysql -u "$USERNAME" -p -e "SHOW DATABASES LIKE '$database_name';" | grep -w "$database_name")

if [ "$DB_EXISTS" ]; then
  echo "The database '$database_name' already exists."
else
  echo "The database '$database_name' does not exist."
  mysql -u $username -p  -e "CREATE DATABASE $database_name;"
  echo "The database '$database_name' is created. Granting permission process starts...."
  mysql -u $username -p -e "GRANT SELECT, INSERT, UPDATE, CREATE , DELETE ON $database_name.* TO 'koosha'@'localhost';"
  echo -e "${GREEN}Database created successfully!${NC}"
fi

echo -e "${GREEN}Starting the table creation process ...${NC}"
mysql -u $username -p $database_name < create_tables.sql
echo -e "${GREEN}Database startup process completed!${NC}"


