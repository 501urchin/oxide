#include <libssh/libssh.h>
#include <iostream>



struct oxideContext {
    ssh_session session;

 
    void ConnectToServerViaPassword(const std::string& host, const std::string& name, const std::string& password) const;
    void ConnectToServerViaKey(const std::string& host, const std::string& name, const std::string& key) const;
};