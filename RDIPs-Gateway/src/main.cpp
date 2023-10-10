#include <headers.h>
#include <utils/Log.h>

#include <algorithm>
#include <amqpcpp.h>
#include <amqpcpp/linux_tcp.h>
#include <amqpcpp/linux_tcp/tcpchannel.h>
#include <amqpcpp/login.h>

#include <connection/MyTcpHandler.h>
#include <cstdio>

void error(const char *msg) {
  perror(msg);
  exit(1);
}

const char *getEnvWithDefault(const char *key, const char *defaultVal) {
  char *result = getenv(key);
  if (result == NULL) {
    return defaultVal;
  }

  return result;
}

int main(int, char const *[]) {
  const char *brokerHost = getEnvWithDefault("BROKER_HOST", "localhost");
  const char *port = getEnvWithDefault("BROKER_PORT", "5672");
  const char *userName = getEnvWithDefault("BROKER_USER", "admin");
  const char *password = getEnvWithDefault("BROKER_PASSWORD", "admin");

  TraceLog(std::source_location::current(), "HOST: ", brokerHost,
           "\n PORT: ", port, "\n UserName: ", userName);
  MyTcpHandler myHandler;

  // address of the server
  AMQP::Address address(brokerHost, (uint16_t)std::atoi(port),
                        AMQP::Login(userName, password), "/");

  // create a AMQP connection object
  AMQP::TcpConnection connection(&myHandler, address);

  auto fds = &myHandler.fds;

  fds->events = POLLIN | POLLOUT;
  while (true) {
    int ret = poll(fds, 1, -1);

    if (ret == -1) {
      error("Error");
    }

    if (fds->revents & POLLIN) {
      connection.process(fds->fd, AMQP::readable);
    }
    if (fds->revents & POLLOUT) {
      connection.process(fds->fd, AMQP::writable);
    }
  }

  delete brokerHost;
  delete port;
  delete userName;
  delete password;

  return 0;
}
