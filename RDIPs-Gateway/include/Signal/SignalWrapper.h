#pragma once

#include <headers.h>

namespace Signal
{
    class SignalWrapper
    {
    public:
        virtual ~SignalWrapper() = default;
    protected:
        AMQP::TcpChannel *channel;
        virtual void Start() = 0;

        virtual void initQueue(std::string name, int flag, AMQP::Table table){};
        virtual void initQueue(std::string name){};
        virtual void initExchange(std::string name, AMQP::ExchangeType type, int flag, AMQP::Table table){};
        virtual void initExchange(std::string name){};
        virtual void listener(std::function<void(const AMQP::Message &message, uint64_t deliveryTag, bool redelivered)> messageCb,
                              std::function<void(const std::string &consumertag)> successCb,
                              std::function<void(const char *message)> errorCb){};
        virtual void bind(std::string routingKey) = 0;
        u_int16_t getChannelId()
        {
            return channel->id();
        }

    };

} // namespace Signal
