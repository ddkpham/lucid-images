from flask import Flask, request, jsonify
from flask_cors import CORS
import pika
import uuid
import json
import os 



class DemoRpcClient(object):
    def __init__(self):
        self.connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
        self.channel = self.connection.channel()

        result = self.channel.queue_declare(queue='', exclusive=True)
        self.callback_queue = result.method.queue

        self.channel.basic_consume(
            queue=self.callback_queue,
            on_message_callback=self.on_response,
            auto_ack=True)

    def on_response(self, ch, method, props, body):
        if self.corr_id == props.correlation_id:
            self.response = json.loads(body.decode('utf-8'))

    def call(self, arg):
        self.response = None
        self.corr_id = str(uuid.uuid4())
        body = json.dumps(arg).encode('utf8')
        
        self.channel.basic_publish(
            exchange='',
            routing_key='rpc_queue',
            properties=pika.BasicProperties(reply_to=self.callback_queue, correlation_id=self.corr_id),
            body=body)
        while self.response is None:
            self.connection.process_data_events()
        return self.response


rpc = DemoRpcClient()

# Initialize app
app = Flask(__name__)

# Enable CORS 
CORS(app)
print("CORS", CORS)

# Add routes 
@app.route('/', methods=['GET'])
def get():
  response = rpc.call({'image': 'argument.png'})
  print(response['length'])
  return jsonify({'image': 'Hello World', 'length': response['length']})

# Image contrast enhancement route 
@app.route('/enhance', methods=['POST'])
def enhance():
    print("request", request.json)
    reqBody = request.json
    response = rpc.call({'image': reqBody['image']})
    print("finished request! ", reqBody)
    return jsonify({'image': 'post test', 'length': '42'})

# Run Server
if __name__ == '__main__':
  app.run(host='0.0.0.0', debug=True, port=5000)