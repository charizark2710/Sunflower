#include <amqpcpp.h>
#include <amqpcpp/linux_tcp.h>
#include <amqpcpp/linux_tcp/tcpchannel.h>
#include <connection/MyTcpHandler.h>
#include <unistd.h>
#include <poll.h>
#include <netdb.h>
#include <sys/select.h>
#include <Queue/SignalQueue.h>

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
    AMQP::Address address("amqp://admin:admin@localhost/");

    // create a AMQP connection object
    AMQP::TcpConnection connection(&myHandler, address);

    AMQP::TcpChannel channel(&connection);

    Queue::SignalQueue *signal = new Queue::SignalQueue(channel);

    auto startCb = [](const std::string &consumertag)
    {
        std::cout << "consume operation started" << std::endl;
    };

    // callback function that is called when the consume operation failed
    auto errorCb = [](const char *message)
    {
        std::cout << message << std::endl;
    };

    // callback operation when a message was received
    auto messageCb = [&](const AMQP::Message &message, uint64_t deliveryTag, bool redelivered)
    {
        printf("parent message received %s \n", message.body());

        // acknowledge the message
        channel.ack(deliveryTag);
        usleep(2000);
    };

    signal->bind("Test");

    signal->listener(messageCb, startCb, errorCb);

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
