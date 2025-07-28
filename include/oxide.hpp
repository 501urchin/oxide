#include <libssh/libssh.h>
#include <iostream>



struct oxideContext {
    ssh_session session;
    const std::string knownHostPath;

    oxideContext(const std::string& knownHostPath)
    : knownHostPath(knownHostPath){};

    ~oxideContext() {
        if (session != nullptr) {
            ssh_free(session);
        }
    }
    
    void ConnectToServerViaPassword(const std::string& host, const std::string& username, const std::string& password) ;
    void ConnectToServerViaKey(const std::string& host, const std::string& username, const std::string& key) ;
};