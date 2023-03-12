#pragma once

#include <Signal/SignalWrapper.h>
#include <headers.h>

namespace Signal {
class TestSignal : public SignalWrapper {
private:
  // void initQueue(std::string name, int flag, AMQP::Table table) override;
  void initQueue(std::string name) override;
  // void initExchange(std::string name, AMQP::ExchangeType type, int flag,
  // AMQP::Table table) override;
  void initExchange(std::string name) override;

  void listener(std::function<void(const AMQP::Message &message,
                                   uint64_t deliveryTag, bool redelivered)>
                    messageCb,
                std::function<void(const std::string &consumertag)> successCb,
                std::function<void(const char *message)> errorCb) override;

  void bind(std::string routingKey) override;
  std::string exchange;
  std::string queue;

public:
  TestSignal(AMQP::TcpChannel &channel, std::string exchange,
             std::string queue);
  ~TestSignal();
  void Start() override;
};

} // namespace Signal
