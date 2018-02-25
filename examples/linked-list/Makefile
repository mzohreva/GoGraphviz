BIN = linked-list

$(BIN): main.go ../../graphviz/*.go
	go build

run: $(BIN)
	./$(BIN) | dot -Tpng -o $(BIN).png

clean:
	go clean
	rm -f $(BIN).png
