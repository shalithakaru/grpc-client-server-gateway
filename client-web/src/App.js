import React, { useState } from 'react';
import './App.css';
import { Message } from './grpc/chat_pb';
import { ChatServiceClient } from './grpc/chat_grpc_web_pb';

const client = new ChatServiceClient('http://localhost:8080', null, null);

function App() {
  const [response, setResponse] = useState('');
  const [streamResponses, setStreamResponses] = useState([]);

  const makeUnaryCall = () => {
    const request = new Message();
    request.setBody('Hello World');
    
    client.unaryChat(request, {}, (err, response) => {
      if (err) {
        console.error(err);
      } else {
        setResponse(response.getBody());
      }
    });
  };

  const makeServerStreamingCall = () => {
    const request = new Message();
    request.setBody('Hello World');

    const stream = client.serverStreamChat(request, {});

    stream.on('data', (response) => {
      setStreamResponses((prev) => [...prev, response.getBody()]);
    });

    stream.on('end', () => {
      console.log('Stream ended.');
    });

    stream.on('error', (err) => {
      console.error(err);
    });
  };

  const makeClientStreamingCall = () => {
    const stream = client.clientStreamChat({}, (err, response) => {
      if (err) {
        console.error(err);
      } else {
        setResponse(response.getBody());
      }
    });

    for (let i = 0; i < 5; i++) {
      const request = new Message();
      request.setBody(`Message ${i}`);
      stream.write(request);
    }
    stream.end();
  };

  const makeBidirectionalStreamingCall = () => {
    const stream = client.bidirectionalStreamChat({});

    stream.on('data', (response) => {
      setStreamResponses((prev) => [...prev, response.getBody()]);
    });

    stream.on('end', () => {
      console.log('Stream ended.');
    });

    stream.on('error', (err) => {
      console.error(err);
    });

    for (let i = 0; i < 5; i++) {
      const request = new Message();
      request.setBody(`Message ${i}`);
      stream.write(request);
    }
    stream.end();
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>gRPC-Web with React</h1>
        <button onClick={makeUnaryCall}>Make Unary Call</button>
        <button onClick={makeServerStreamingCall}>Make Server Streaming Call</button>
        {/* Please read the README.md why below are not supported yet */}
        {/* <button onClick={makeClientStreamingCall}>Make Client Streaming Call</button> */}
        {/* <button onClick={makeBidirectionalStreamingCall}>Make Bidirectional Streaming Call</button> */}
        <p>Unary Response: {response}</p>
        <div>
          <p>Stream Responses:</p>
          <ul>
            {streamResponses.map((msg, idx) => (
              <li key={idx}>{msg}</li>
            ))}
          </ul>
        </div>
      </header>
    </div>
  );
}

export default App;
