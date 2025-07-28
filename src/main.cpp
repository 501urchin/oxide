#include "oxide.hpp"

int main(){
    oxideContext ctx("hello.txt");


    try {
        ctx.ConnectToServer("127.0.0.1:8080", "testuser", SSH_PASSWORD_AUTH, "password1d23");
    }catch(const std::exception& e) {
        std::cout << e.what() << std::endl;
    }
}