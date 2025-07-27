#include <libssh/libssh.h>
#include <stdio.h>

ssh_session connect_ssh(const char *host, const char *user, int verbosity)
{
    ssh_session session = NULL;
    int auth = 0;

    session = ssh_new();
    if (session == NULL) {
        return NULL;
    }

    if (user != NULL) {
        if (ssh_options_set(session, SSH_OPTIONS_USER, user) < 0) {
            ssh_free(session);
            return NULL;
        }
    }

    if (ssh_options_set(session, SSH_OPTIONS_HOST, host) < 0) {
        ssh_free(session);
        return NULL;
    }
    ssh_options_set(session, SSH_OPTIONS_LOG_VERBOSITY, &verbosity);
    if (ssh_connect(session)) {
        fprintf(stderr, "Connection failed : %s\n", ssh_get_error(session));
        ssh_disconnect(session);
        ssh_free(session);
        return NULL;
    }


    ssh_disconnect(session);
    ssh_free(session);
    return NULL;
}
