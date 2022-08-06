#include <amqpcpp.h>
#include <sys/socket.h>
#include <poll.h>
#include <arpa/inet.h>
#include <unistd.h>

class MyconnectionHandler : public AMQP::ConnectionHandler
{

    virtual void onProperties(AMQP::Connection *connection, const AMQP::Table &server, AMQP::Table &client)
    {
        printf("TEST");
    }

    /**
     *  Method that is called by the AMQP library every time it has data
     *  available that should be sent to RabbitMQ.
     *  @param  connection  pointer to the main connection object
     *  @param  data        memory buffer with the data that should be sent to RabbitMQ
     *  @param  size        size of the buffer
     */
    virtual void onData(AMQP::Connection *connection, const char *data, size_t size)
    {
        // @todo
        //  Add your own implementation, for example by doing a call to the
        //  send() system call. But be aware that the send() call may not
        //  send all data at once, so you also need to take care of buffering
        //  the bytes that could not immediately be sent, and try to send
        //  them again when the socket becomes writable again
    }

    /**
     *  Method that is called by the AMQP library when the login attempt
     *  succeeded. After this method has been called, the connection is ready
     *  to use.
     *  @param  connection      The connection that can now be used
     */
    virtual void
    onReady(AMQP::Connection *connection)
    {
        // @todo
        //  add your own implementation, for example by creating a channel
        //  instance, and start publishing or consuming

        // and create a channel
        AMQP::Channel channel(connection);

        auto startCb = [](const std::string &consumertag)
        {
            std::cout << "consume operation started" << std::endl;
        };

        // callback function that is called when the consume operation failed
        auto errorCb = [](const char *message)
        {
            std::cout << "consume operation failed" << std::endl;
        };

        // callback operation when a message was received
        auto messageCb = [&channel](const AMQP::Message &message, uint64_t deliveryTag, bool redelivered)
        {
            std::cout << "message received" << std::endl;

            // acknowledge the message
            channel.ack(deliveryTag);
        };

        // start consuming from the queue, and install the callbacks
        channel.consume("my-queue")
            .onReceived(messageCb)
            .onSuccess(startCb)
            .onError(errorCb);
    }

    /**
     *  Method that is called by the AMQP library when a fatal error occurs
     *  on the connection, for example because data received from RabbitMQ
     *  could not be recognized.
     *  @param  connection      The connection on which the error occured
     *  @param  message         A human readable error message
     */
    virtual void onError(AMQP::Connection *connection, const char *message)
    {
        // @todo
        //  add your own implementation, for example by reporting the error
        //  to the user of your program, log the error, and destruct the
        //  connection object because it is no longer in a usable state

        printf("%s \n", message);
    }

    /**
     *  Method that is called when the connection was closed. This is the
     *  counter part of a call to Connection::close() and it confirms that the
     *  AMQP connection was correctly closed.
     *
     *  @param  connection      The connection that was closed and that is now unusable
     */
    virtual void onClosed(AMQP::Connection *connection)
    {
        // @todo
        //  add your own implementation, for example by closing down the
        //  underlying TCP connection too
        printf("CLOSING \n");
    }
};