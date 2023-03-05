#include "main.h"

#include <amqpcpp/login.h>
#include <amqpcpp.h>
#include <amqpcpp/linux_tcp.h>
#include <amqpcpp/linux_tcp/tcpchannel.h>

#include <connection/MyTcpHandler.h>

void error(const char *msg)
{
    perror(msg);
    exit(1);
}

const char* getEnvWithDefault(const char* key, const char* defaultVal) {
    char* result = getenv(key);
    if (result == NULL) {
        return defaultVal;
    }

    return  result;
}

int main(int argc, char const *argv[])
{
    printf("Start %s \n", "OK");
    
    const char *brokerHost = getEnvWithDefault("BROKER_HOST", "localhost");
    const char *port = getEnvWithDefault("BROKER_PORT", "5672");
    const char *userName = getEnvWithDefault("BROKER_USER", "admin");
    const char *password = getEnvWithDefault("BROKER_PASSWORD", "admin");


    MyTcpHandler myHandler;

    // address of the server
    AMQP::Address address( brokerHost, (uint16_t)std::atoi(port), AMQP::Login(userName, password), "/");

    // create a AMQP connection object
    AMQP::TcpConnection connection(&myHandler, address);
    printf("Connect %s \n", "OK");

    auto fds = &myHandler.fds;

    fds->events = POLLIN | POLLOUT;
    while (true)
    {
        int ret = poll(fds, 1, -1);

        if (fds->revents & POLLIN)
        {
            connection.process(fds->fd, AMQP::readable);
        }
        else if (fds->revents & POLLOUT)
        {
            connection.process(fds->fd, AMQP::writable);
        }

    }
    return 0;
}
