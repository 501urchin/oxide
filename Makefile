CXX = clang++
CXXFLAGS = -std=c++20 -Wall -O3 -I./include -I/usr/local/include
LDFLAGS = -L/usr/local/lib -Wl,-rpath,/usr/local/lib -lssh

SRCS = ./src/main.cpp ./src/oxide.cpp ./src/ssh.cpp
BIN = oxide

build: $(SRCS)
	@$(CXX) $(CXXFLAGS) $(SRCS) $(LDFLAGS) -o $(BIN)

run: build
	@./$(BIN)
	@rm -f $(BIN)

run-time: build
	@time ./$(BIN)
	@rm -f $(BIN)

clean:
	@rm -f $(BIN)
