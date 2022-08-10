#pragma once

#include <headers.h>

namespace Queue
{
    class QueueWrapper
    {
    private:
        AMQP::TcpChannel &channel;

    public:
        QueueWrapper(AMQP::TcpChannel &channel, int &flag);
        ~QueueWrapper();
        virtual void initQueue(std::string name, int flag, AMQP::Table table) = 0;
        virtual void initQueue(std::string name) = 0;
        virtual void initExchange(std::string name, AMQP::ExchangeType type, int flag, AMQP::Table table) = 0;
        virtual void initExchange(std::string name) = 0;
        virtual void listener(std::string name,
                      std::function<void(const AMQP::Message &message, uint64_t deliveryTag, bool redelivered)> messageCb,
                      std::function<void(const std::string &consumertag)> successCb,
                      std::function<void(const char *message)> errorCb) = 0;
    };
} // namespace Queue
