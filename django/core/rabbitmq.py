from kombu import Connection

def create_rabbitmq_connection()-> Connection:
    return Connection("amqp://guest:guest@host.docker.internal:5672//")