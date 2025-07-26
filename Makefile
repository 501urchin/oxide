CXX = clang++
CXXFLAGS = -std=c++17 -Wall -O3  
SRC = main.cpp
BIN = oxide


build: $(SRC)
	@$(CXX) $(CXXFLAGS) $(SRC) -o $(BIN)

run: build
	@./$(BIN)
	@rm -f $(BIN)

run-time: build
	@time ./$(BIN)
	@rm -f $(BIN)

clean:
	@rm -f $(BIN)
