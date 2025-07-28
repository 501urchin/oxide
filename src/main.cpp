#include "oxide.hpp"

int main(){
    oxideContext ctx("hello.txt");


    try {
        ctx.ConnectToServerViaPassword("127.0.0.1:8080", "testuser", "password123");
    }catch(const std::exception& e) {
        std::cout << e.what() << std::endl;
    }
}