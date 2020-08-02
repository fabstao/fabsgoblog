export DOCKER=/usr/local/bin/docker
export MYGO=/usr/local/bin/go
export GOOS=linux
export GOARCH=amd64
all:
	${MYGO} build -o build/server *.go

run:
	${MYGO} run server.go

docker:
	${DOCKER} build -t quay.io/fabstao/fabsgoblog .

docker-run:
	${DOCKER} run -d --name fabsgoblog -p 8019:8019 -e  SITEKEY="6LdR07EZAAAAAG594Nlkla3OhEXE-6DOzvip5avv" \
	-e FGOSECRET="Tequ1squiapan" -e PASSWD_ADMIN="fabsgobl0g" -e DBHOST="192.168.56.1" \
	-e DBPORT="5432" -e DBNAME="fabsgoblogdb" -e DBUSER="fabs" -e DBPASSWD="austin23" \
	quay.io/fabstao/fabsgoblog

clean:
	rm -f build/*




