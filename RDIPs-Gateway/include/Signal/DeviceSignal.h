#pragma once

#include <Signal/SignalWrapper.h>
#include <headers.h>
#include <string>

const inline std::string APIs[] = {
    "GetAllDevices", "PostDevice", "GetDetailDevice", "PutDetailDevice",
    "DeleteDevice",
    // Performances
    "GetAllPerformances", "GetDetailPerformance", "PutDetailPerformance",
    // History
    "GetDetailHistory", "PutDetailHistory",
    // Weather
    "GetWeatherForecast"};

namespace Signal {
class DeviceSignal : public SignalWrapper {
private:
  // void init(std::string name, int flag, AMQP::Table table) override;
  void init() override;

  void listener(std::function<void(const AMQP::Message &message,
                                   uint64_t deliveryTag, bool redelivered)>
                    messageCb,
                std::function<void(const std::string &consumertag)> successCb,
                std::function<void(const char *message)> errorCb) override;

  void bind(std::string routingKey) override;

public:
  DeviceSignal(AMQP::TcpChannel &channel, std::string exchange,
               std::string queue);
  ~DeviceSignal();
  void Start() override;
};

} // namespace Signal
