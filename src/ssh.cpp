#include "oxide.hpp"
#include <stdexcept>
#include <libssh/libssh.h>


void verifyKnownHost(ssh_session session) {
    unsigned char *hash = NULL;
    size_t hlen;
    ssh_key srv_pubkey = NULL;
    int rc;

    try {
        rc = ssh_get_server_publickey(session, &srv_pubkey);
        if (rc < 0) {
            throw std::runtime_error("failed to get server public key");
        }

        rc = ssh_get_publickey_hash(srv_pubkey, SSH_PUBLICKEY_HASH_SHA256, &hash, &hlen);
        if (rc < 0) {
            throw std::runtime_error("failed to get server public key hash");
        }

        int state = ssh_session_is_known_server(session);
        if (state == SSH_SERVER_NOT_KNOWN) {
            std::string answer;
            std::cout <<  "The server is unknown. Do you trust the host key (yes/no)?" << std::endl;
            ssh_print_hash(SSH_PUBLICKEY_HASH_SHA256, hash, hlen);
            std::cin >> answer;

            if (answer == "yes" || answer == "y") {
                ssh_session_update_known_hosts(session);
            }

        }

    }catch(const std::exception& e) {
        ssh_key_free(srv_pubkey);
        ssh_clean_pubkey_hash(&hash);
        throw std::runtime_error(e.what());
    }
}

void oxideContext::ConnectToServerViaPassword(const std::string& host, const std::string& username, const std::string& password)  {
    try {
        if (host.empty()) {
            throw std::invalid_argument("please provide a valid host");
        }

        if (username.empty()) {
            throw std::invalid_argument("please provide a valid username");
        }

        if (password.empty()) {
            throw std::invalid_argument("please provide a valid password");
        }

        this->session = ssh_new();
        if (this->session == NULL) {
            throw std::invalid_argument("failed to create new ssh session");
        }

        if (ssh_options_set(session, SSH_OPTIONS_USER, username.c_str()) < 0){
              throw std::runtime_error(ssh_get_error(this->session));
        }

        if (ssh_options_set(session, SSH_OPTIONS_PASSWORD_AUTH, password.c_str()) < 0){
              throw std::runtime_error(ssh_get_error(this->session));
        }

        if (ssh_options_set(session, SSH_OPTIONS_KNOWNHOSTS, this->knownHostPath.c_str()) < 0){
              throw std::runtime_error(ssh_get_error(this->session));
        }


        size_t idx = host.find(":");
        if (idx != std::string::npos) {
            int port = std::atoi(host.substr(idx + 1).c_str());
            std::string h  = host.substr(0, idx);

            if (ssh_options_set(session, SSH_OPTIONS_PORT, &port) < 0) {
                throw std::runtime_error(ssh_get_error(this->session));
            }

            if (ssh_options_set(session, SSH_OPTIONS_HOST, h.c_str()) < 0){
                throw std::runtime_error(ssh_get_error(this->session));
            }
        }else {
            if (ssh_options_set(session, SSH_OPTIONS_HOST, host.c_str()) < 0){
                throw std::runtime_error(ssh_get_error(this->session));
            }
        }

        
        int status = ssh_connect(this->session);
        if (status != SSH_OK) {
            throw std::runtime_error(ssh_get_error(this->session));
        }

        verifyKnownHost(this->session);

        std::cout << "Connected!!!\n";

    }catch (const std::exception& e) {
        ssh_disconnect(this->session);
        ssh_free(this->session);
        throw std::runtime_error(e.what());
    }
}
