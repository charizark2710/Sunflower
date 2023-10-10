#include "headers.h"
#include <Signal/DeviceSignal.h>
#include <cstring>
#include <source_location>
#include <string>
#include <utils/Log.h>

namespace Signal {
DeviceSignal::DeviceSignal(AMQP::TcpChannel &channel, std::string exchange,
                           std::string queue) {
  this->channel = &channel;
  this->exchange = exchange;
  this->queue = queue;

  this->init();
}

DeviceSignal::~DeviceSignal() { channel->close(); }

void DeviceSignal::init() {
  for (auto api : APIs) {
    this->bind(GW_ROUTING_KEY + api);
  }
  channel->setQos(QOS);
}

void DeviceSignal::listener(
    std::function<void(const AMQP::Message &message, uint64_t deliveryTag,
                       bool redelivered)>
        messageCb,
    std::function<void(const std::string &consumertag)> successCb,
    std::function<void(const char *message)> errorCb) {
  // start consuming from the queue, and install the callbacks
  channel->consume(queue).onReceived(messageCb).onSuccess(successCb).onError(
      errorCb);
}

void DeviceSignal::bind(std::string routingKey) {
  channel->bindQueue(exchange, queue, routingKey);
}

void DeviceSignal::Start() {

  auto startCb = [](const std::string &consumertag) {
    std::cout << "consume operation started" << std::endl;
  };

  // callback function that is called when the consume operation failed
  auto errorCb = [](const char *message) { std::cout << message << std::endl; };

  // callback operation when a message was received
  auto messageCb = [&](const AMQP::Message &message, uint64_t deliveryTag,
                       bool redelivered) {
    TraceLog(std::source_location::current(), "parent message received",
             message.body());

    // acknowledge the message
    channel->ack(deliveryTag);
    // channel->publish(EXCHANGE_NAME, "Test", "OK Confirm");
  };

  this->listener(messageCb, startCb, errorCb);
}
} // namespace Signal
