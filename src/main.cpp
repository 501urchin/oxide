#include "oxide.hpp"

int main(){
    oxideContext ctx("hello.txt");


    try {
        ctx.ConnectToServer("127.0.0.1:8080", "testuser", SSH_PASSWORD_AUTH, "password123");
        ctx.Execute("lsd");
    }catch(const std::exception& e) {
        std::cout << e.what() << std::endl;
    }
}