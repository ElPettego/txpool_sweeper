TARGET = txpl_swpr
SRC    = src/txpl_swpr.go

$(TARGET) : $(SRC)
	go build $(SRC)

run:
	./$(TARGET) -network binance