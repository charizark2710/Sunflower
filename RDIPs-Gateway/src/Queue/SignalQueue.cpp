#include <Queue/SignalQueue.h>

namespace Queue
{
    SignalQueue::SignalQueue(AMQP::TcpChannel &channel)
    {
        this->channel = &channel;
        this->initQueue(QUEUE_NAME);
        this->initExchange(EXCHANGE_NAME);
    }

    SignalQueue::~SignalQueue()
    {
        channel->close();
        delete channel;
    }

    void SignalQueue::initQueue(std::string name)
    {
        std::cout << "Init queue \n";
        channel->declareQueue(name);
        channel->setQos(QOS);
        std::cout << "Init queue Succeed \n";
    }

    void SignalQueue::initExchange(std::string name)
    {
        channel->declareExchange(name);
    }

    void SignalQueue::listener(std::function<void(const AMQP::Message &message, uint64_t deliveryTag, bool redelivered)> messageCb,
                               std::function<void(const std::string &consumertag)> successCb,
                               std::function<void(const char *message)> errorCb)
    {
        // start consuming from the queue, and install the callbacks
        channel->consume(QUEUE_NAME)
            .onReceived(messageCb)
            .onSuccess(successCb)
            .onError(errorCb);
    }

    void SignalQueue::bind(std::string routingKey)
    {
        channel->bindQueue(EXCHANGE_NAME, QUEUE_NAME, routingKey);
        std::cout << "succeed \n";
    }
} // namespace Queue
