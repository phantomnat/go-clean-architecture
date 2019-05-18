
.PHONY: start-mysql stop-mysql


start-mysql:
	docker run --name mysql -e MYSQL_ROOT_PASSWORD=123456 -p 3306:3306 -d mysql:5.7

stop-mysql:
	docker stop mysql
	docker stop mysql