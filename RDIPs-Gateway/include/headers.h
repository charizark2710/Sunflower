#pragma once

#include <cstring>
#include <format>
#include <string>

#include <amqpcpp.h>
#include <amqpcpp/linux_tcp.h>
#include <iostream>

#include <poll.h>
#include <sys/fcntl.h>

#include <string.h>

#include <cstdlib>
#include <netdb.h>
#include <source_location>
#include <sstream>
#include <sys/select.h>
#include <unistd.h>


#define QUEUE_API "API"
#define GW_ROUTING_KEY "gateway.*"
#define DEVICE_ROUTING_KEY "device.*"
#define EXCHANGE "amq.topic"
#define QOS 100
