#include <Signal/TestSignal.h>

namespace Signal {
TestSignal::TestSignal(AMQP::TcpChannel &channel, std::string exchange,
                       std::string queue) {
  this->channel = &channel;
  this->initExchange(exchange);
  this->initQueue(queue);
}

TestSignal::~TestSignal() { channel->close(); }

void TestSignal::initQueue(std::string name) {
  queue = name;
  channel->declareQueue(name);
  channel->setQos(QOS);
}

void TestSignal::initExchange(std::string name) {
  exchange = name;
  channel->declareExchange(name, AMQP::topic);
}

void TestSignal::listener(
    std::function<void(const AMQP::Message &message, uint64_t deliveryTag,
                       bool redelivered)>
        messageCb,
    std::function<void(const std::string &consumertag)> successCb,
    std::function<void(const char *message)> errorCb) {
  // start consuming from the queue, and install the callbacks
  channel->consume(queue).onReceived(messageCb).onSuccess(successCb).onError(
      errorCb);
}

void TestSignal::bind(std::string routingKey) {
  channel->bindQueue(exchange, queue, routingKey);
}

void TestSignal::Start() {

  auto startCb = [](const std::string &consumertag) {
    std::cout << "consume operation started" << std::endl;
  };

  // callback function that is called when the consume operation failed
  auto errorCb = [](const char *message) { std::cout << message << std::endl; };

  // callback operation when a message was received
  auto messageCb = [&](const AMQP::Message &message, uint64_t deliveryTag,
                       bool redelivered) {
    printf("parent message received %s \n", message.body());

    // acknowledge the message
    channel->ack(deliveryTag);
    // channel->publish(EXCHANGE_NAME, "Test", "OK Confirm");
  };

  // get everythings
  this->bind("#");

  this->listener(messageCb, startCb, errorCb);

  //   channel->publish(exchange, "Test", "TEEST");
}
} // namespace Signal
