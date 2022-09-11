#include <amqpcpp.h>
#include <amqpcpp/linux_tcp.h>
#include <amqpcpp/linux_tcp/tcpchannel.h>
#include <connection/MyTcpHandler.h>
#include <unistd.h>
#include <poll.h>
#include <netdb.h>
#include <sys/select.h>

void error(const char *msg)
{
    perror(msg);
    exit(1);
}

int main(int argc, char const *argv[])
{
    printf("Start %s \n", "OK");

    MyTcpHandler myHandler;

    // address of the server
    AMQP::Address address("amqp://admin:admin@sfBroker:5672/");

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
