#include <amqpcpp.h>
#include <amqpcpp/linux_tcp.h>
#include <amqpcpp/linux_tcp/tcpchannel.h>
#include <connection/MyTcpHandler.h>
#include <amqpcpp/libev.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <poll.h>
#include <netdb.h>
#include <sys/select.h>
#include <chrono>

void error(const char *msg)
{
    perror(msg);
    exit(1);
}

// void setUpRabbitmq()
// {
//     MyconnectionHandler myHandler;
//     // create a AMQP connection object
//     AMQP::Connection connection(&myHandler, AMQP::Login("guest", "guest"), "/");
// }

struct hostent *server;

// void setUpSocket()
// {

//     server = gethostbyname("localhost");
//     struct pollfd fds[1];
//     int ret, pret;
//     char *buf = (char *)malloc(256);

//     int sockfd, new_socket, valread;
//     struct sockaddr_in address;
//     int opt = 1;
//     int addrlen = sizeof(address);
//     char buffer[1024] = {0};
//     char *hello = "Hello from server";
//     sockfd = socket(AF_INET, SOCK_STREAM, 0);
//     if (sockfd <= 0)
//     {
//         error("socket failed");
//     }

//     if (setsockopt(sockfd, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT,
//                    &opt, sizeof(opt)))
//     {
//         error("setsockopt");
//     }

//     address.sin_family = AF_INET;
//     address.sin_addr.s_addr = *server->h_addr_list[0];
//     address.sin_port = htons(5672);

//     if (inet_pton(AF_INET, "172.17.0.1", &address.sin_addr) <= 0)
//     {
//         error("\nInvalid address/ Address not supported \n");
//     }

//     if (connect(sockfd, (struct sockaddr *)&address, sizeof(address)) < 0)
//     {
//         error("\nConnection Failed \n");
//     }

//     while (true)
//     {
//         fds[0].fd = sockfd;
//         fds[0].events = 0;
//         fds[0].events |= POLLIN;

//         pret = poll(fds, 1, -1);

//         if (pret == 0)
//         {
//             error("timeout");
//         }
//         else
//         {
//             memset(buf, 0, 256);
//             ret = read(sockfd, buf, 256);
//             printf("%d \n", ret);

//             if (ret != -1)
//             {
//                 printf("%s \n", buf);
//             }
//         }
//     }
// }

int main(int argc, char const *argv[])
{
    MyTcpHandler myHandler;

    // address of the server
    AMQP::Address address("amqp://admin:admin@localhost/");

    // create a AMQP connection object
    AMQP::TcpConnection connection(&myHandler, address);

    AMQP::TcpChannel channel(&connection);

    channel.setQos(1);

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
        printf("parent message received %s at %d \n", message.body());

        // acknowledge the message
        channel.ack(deliveryTag);
        usleep(2000);
    };

    // start consuming from the queue, and install the callbacks
    channel.consume("hello")
        .onReceived(messageCb)
        .onSuccess(startCb)
        .onError(errorCb);

    char *c[] = {"info", "warning", "error"};

    for (int i = 0; i < 3; i++)
    {
        channel.bindQueue("weirdEx", "hello", c[i]);
    }
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
