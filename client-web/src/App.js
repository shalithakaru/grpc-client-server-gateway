// Add these imports at the top
import { grpc } from '@improbable-eng/grpc-web';
import { HelloRequest, HelloReply } from './grpc/chat_pb';
import { GreeterClient } from './grpc/chat_grpc_web_pb';

function App() {
  const [response, setResponse] = useState('');
  const [streamResponses, setStreamResponses] = useState([]);

  const makeUnaryCall = () => {
    const request = new HelloRequest();
    request.setName('World');

    client.sayHello(request, {}, (err, response) => {
      if (err) {
        console.error(err);
      } else {
        setResponse(response.getMessage());
      }
    });
  };

  const makeServerStreamingCall = () => {
    const request = new HelloRequest();
    request.setName('World');

    const stream = client.sayHelloStream(request, {});

    stream.on('data', (response) => {
      setStreamResponses((prev) => [...prev, response.getMessage()]);
    });

    stream.on('end', () => {
      console.log('Stream ended.');
    });

    stream.on('error', (err) => {
      console.error(err);
    });
  };

  // Implement clientStreamingCall and bidirectionalStreamingCall

  return (
    <div className="App">
      <header className="App-header">
        <h1>gRPC-Web with React</h1>
        <button onClick={makeUnaryCall}>Make Unary Call</button>
        <button onClick={makeServerStreamingCall}>Make Server Streaming Call</button>
        {/* Add buttons for other calls */}
        <p>Unary Response: {response}</p>
        <div>
          <h2>Stream Responses:</h2>
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
