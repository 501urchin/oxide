CXX = clang++
CXXFLAGS = -std=c++17 -Wall -O3 -I./include -I/usr/local/include -I/opt/homebrew/opt/openssl@3/include
LDFLAGS = -L/usr/local/lib -L/opt/homebrew/opt/openssl@3/lib -Wl,-rpath,/usr/local/lib -Wl,-rpath,/opt/homebrew/opt/openssl@3/lib -lssh -lcrypto -lz
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
