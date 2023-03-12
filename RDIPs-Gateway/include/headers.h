#pragma once

#define QUEUE_TEST "TEST"
#define EXCHANGE_TEST "TEST_EXCHANGE"
#define QOS 100

#include <iostream>
#include <amqpcpp.h>
#include <amqpcpp/linux_tcp.h>

#include <poll.h>
#include <sys/fcntl.h>

#include <string.h>