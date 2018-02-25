BIN = linked-list2

$(BIN): main.go ../../graphviz/*.go
	go build

run: $(BIN)
	./$(BIN)

clean:
	go clean
	rm -f list-*.png
