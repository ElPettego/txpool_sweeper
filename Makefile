TARGET = txpl_swpr
SRC    = src/txpl_swpr.go

$(TARGET) : $(SRC)
	go build -o build/$(TARGET) $(SRC)

run:
	swg_tmr ./build/$(TARGET) -network binance

contract:
	solc --bin --abi --optimize --overwrite --evm-version paris src/corsa.sol -o build/solidity/

crun:
	swg_tmr go build -o build/$(TARGET) $(SRC) && swg_tmr ./build/$(TARGET) -network binance 

init_db:
	swg_tmr go build -o build/$(TARGET) $(SRC) && swg_tmr ./build/$(TARGET) -init_db yes 
