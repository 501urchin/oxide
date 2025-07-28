#include <libssh/libssh.h>
#include <iostream>


enum oxideSshAuthMethod {
    SSH_KEY_AUTH,
    SSH_PASSWORD_AUTH,
};

struct oxideContext {
    private:
        ssh_session session;
        const std::string knownHostPath;
        

    public:
        oxideContext(const std::string& knownHostPath)
        : knownHostPath(knownHostPath){
            session = NULL;
        };

         ~oxideContext() {
            if (session != nullptr) {
                ssh_free(session);
            }
        }


        void ConnectToServer(const std::string& host, const std::string& username, oxideSshAuthMethod authMethod, const std::string& auth);
}; 