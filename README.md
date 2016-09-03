# spital

## Init DB

```bash
./spital-util data/boxes.csv data/tools.csv data/compisition.csv 
```

##Running the app

To run the whole app (server + db + client) you need a *docker-compose.yml* file :

```yaml
version: '2'
services:
  db:
    image: mysql:5.7
    volumes:
      - "~/.hospital_data/db:/var/lib/mysql"
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: xxx
      MYSQL_DATABASE: xxx
      MYSQL_USER: xxx
      MYSQL_PASSWORD: xxx
  spital-server:
    depends_on:
      - db
    image: fabienfoerster/spital-server:latest
    links:
      - db
    ports:
      - "5000:5000"
    restart: always
    environment:
      MYSQL_HOST: db:3306
      MYSQL_USER: xxx
      MYSQL_PASSWORD: xxx
  spital-client:
    image: fabienfoerster/spital-client:latest
    ports:
      - "3000:80"
    restart: always
    environment:
      SERVER_URL: spital-server:5000

```
