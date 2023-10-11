#pragma once

#include <headers.h>
#include <vector>

namespace Signal {
class SignalWrapper {
public:
  virtual ~SignalWrapper() = default;

protected:
  AMQP::TcpChannel *channel;
  std::string exchange;
  std::string queue;

  virtual void Start() = 0;

  virtual void init() = 0;
  virtual void
      listener(std::function<void(const AMQP::Message &message,
                                  uint64_t deliveryTag, bool redelivered)>,
               std::function<void(const std::string &consumertag)>,
               std::function<void(const char *message)>){};
  virtual void bind(std::string routingKey) = 0;
  u_int16_t getChannelId() { return channel->id(); }
};

} // namespace Signal
