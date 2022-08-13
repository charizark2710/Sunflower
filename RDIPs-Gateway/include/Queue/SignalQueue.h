#pragma once

#include <headers.h>
#include <Queue/QueueWrapper.h>

#define QUEUE_NAME "TEST"
#define EXCHANGE_NAME "TEST_EXCHANGE"
#define QOS 100

namespace Queue
{
    class SignalQueue : public QueueWrapper
    {
    private:
        // virtual void initQueue(std::string name, int flag, AMQP::Table table) override;
        virtual void initQueue(std::string name) override;
        // virtual void initExchange(std::string name, AMQP::ExchangeType type, int flag, AMQP::Table table) override;
        virtual void initExchange(std::string name) override;

    public:
        SignalQueue(AMQP::TcpChannel &channel);
        ~SignalQueue();

        void listener(std::function<void(const AMQP::Message &message, uint64_t deliveryTag, bool redelivered)> messageCb,
                      std::function<void(const std::string &consumertag)> successCb,
                      std::function<void(const char *message)> errorCb) override;

        void bind(std::string routingKey) override;
    };

} // namespace Queue
