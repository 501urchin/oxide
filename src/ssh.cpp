#include "oxide.hpp"
#include <stdexcept>
#include <libssh/libssh.h>
#include <filesystem>


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
        if (state == SSH_SERVER_NOT_KNOWN || state == SSH_KNOWN_HOSTS_CHANGED || state == SSH_KNOWN_HOSTS_NOT_FOUND){
            std::string answer;
            std::cout << "The server is unknown. Do you trust the host key (yes/no)?" << std::endl;
            ssh_print_hash(SSH_PUBLICKEY_HASH_SHA256, hash, hlen);
            std::cin >> answer;

            if (answer == "yes" || answer == "y"){
                ssh_session_update_known_hosts(session);
            }
        }

        if (state == SSH_KNOWN_HOSTS_OTHER) {
            throw std::runtime_error("The host key for this server was not found but an other type of key exists. An attacker might change the default server key to confuse your client into thinking the key does not exist");
        }

        ssh_key_free(srv_pubkey);
        ssh_clean_pubkey_hash(&hash);
    }catch(const std::exception& e) {
        ssh_key_free(srv_pubkey);
        ssh_clean_pubkey_hash(&hash);
        throw std::runtime_error(e.what());
    } 
}

void oxideContext::ConnectToServer(const std::string& host, const std::string& username, oxideSshAuthMethod authMethod, const std::string& auth)  {
    if (host.empty()) {
        throw std::invalid_argument("please provide a valid host");
    }

    if (username.empty()) {
        throw std::invalid_argument("please provide a valid username");
    }

    if (auth.empty()) {
        throw std::invalid_argument("please provide a valid auth");
    }

    this->session = ssh_new();
    if (this->session == NULL) {
        throw std::invalid_argument("failed to create new ssh session");
    }
    
    try {
        if (ssh_options_set(this->session, SSH_OPTIONS_USER, username.c_str()) < 0){
              throw std::runtime_error(ssh_get_error(this->session));
        }
        
        if (ssh_options_set(this->session, SSH_OPTIONS_KNOWNHOSTS, knownHostPath.c_str()) < 0){
              throw std::runtime_error(ssh_get_error(this->session));
        }

        size_t idx = host.find(":");
        if (idx != std::string::npos) {
            long port = std::stol(host.substr(idx + 1).c_str());
            std::string h  = host.substr(0, idx);

            if (ssh_options_set(this->session, SSH_OPTIONS_PORT, &port) < 0) {
                throw std::runtime_error(ssh_get_error(this->session));
            }

            if (ssh_options_set(this->session, SSH_OPTIONS_HOST, h.c_str()) < 0){
                throw std::runtime_error(ssh_get_error(this->session));
            }
        }else {
            if (ssh_options_set(this->session, SSH_OPTIONS_HOST, host.c_str()) < 0){
                throw std::runtime_error(ssh_get_error(this->session));
            }
        }

        int status = ssh_connect(this->session);
        if (status != SSH_OK) {
            throw std::runtime_error(ssh_get_error(this->session));
        }

        verifyKnownHost(this->session);
        if(authMethod == oxideSshAuthMethod::SSH_KEY_AUTH) {  
            if (!std::filesystem::exists(auth)) {
                throw std::runtime_error("private key file doesn't exist");
            }

            ssh_key privkey = nullptr;
            int rc = ssh_pki_import_privkey_file(auth.c_str(), nullptr, nullptr, nullptr, &privkey);
            if (rc != SSH_OK) {
                throw std::runtime_error(ssh_get_error(this->session));
            }

            rc = ssh_userauth_publickey(session, username.c_str(), privkey);
            ssh_key_free(privkey);
            if (rc < 0){
                throw std::runtime_error(ssh_get_error(this->session));
            }
        }else if (authMethod == oxideSshAuthMethod::SSH_PASSWORD_AUTH) {
            if (ssh_userauth_password(this->session, username.c_str(), auth.c_str()) != SSH_AUTH_SUCCESS) {
                throw std::runtime_error(ssh_get_error(this->session));
            }
        }else {
            throw std::runtime_error("invalid authMethod flag");
        }
        

    }catch (const std::exception& e) {
        if (this->session != nullptr) {
            if (ssh_is_connected(this->session)) {
                ssh_disconnect(this->session);
            }
            ssh_free(this->session);
            this->session = nullptr;
        }
        throw;
    }
}



const std::string& oxideContext::Execute(const std::string& command) {
    ssh_channel channel = ssh_channel_new(this->session);
    if (channel == NULL) {
        throw std::runtime_error("failed to open new ssh channel");
    }

    try {

        int rc = ssh_channel_open_session(channel);
        if (rc < 0) {
            throw std::runtime_error("failed to open session on channel");
        }

        rc = ssh_channel_request_exec(channel, command.c_str());
        if (rc != SSH_OK) {
            throw std::runtime_error(std::format("failed to execute command '{}' on channel", command));
        }


        // next we need to read the output of the command 
        // i think we can use this to do that but im not sure if this does what i want and idk how to dynamically allocate the needed buf
        // also we need to decide how to return the result, as a cpp string or a c buf 
        // rbytes = ssh_channel_read(channel, buffer, sizeof(buffer), 0);
        // if (rbytes <= 0) {
        //     goto failed;
        // }


    }catch(const std::exception& e) {
        ssh_channel_send_eof(channel);
        ssh_channel_close(channel);
        ssh_channel_free(channel);
        throw;
    }
}