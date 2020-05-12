initial:
	sqlite3 ./todos.db "create table todos (title varchar(100),done boolean)"
clean:
	go clean --cache
run: clean
	go run main.go
hot:
	re go run main.go
