const { Resource } = require("@opentelemetry/resources");

const {
  SemanticResourceAttributes,
} = require("@opentelemetry/semantic-conventions");

const {
  WebTracerProvider,
  SimpleSpanProcessor,
  ConsoleSpanExporter,
} = require("@opentelemetry/sdk-trace-web");

const {
  OTLPTraceExporter,
} = require("@opentelemetry/exporter-trace-otlp-http");

const { registerInstrumentations } = require("@opentelemetry/instrumentation");

const {
  FetchInstrumentation,
} = require("@opentelemetry/instrumentation-fetch");

const consoleExporter = new ConsoleSpanExporter();

const collectorExporter = new OTLPTraceExporter({
  headers: {},
});

const provider = new WebTracerProvider({
  resource: new Resource({
    [SemanticResourceAttributes.SERVICE_NAME]: process.env.REACT_APP_NAME,
  }),
});

const fetchInstrumentation = new FetchInstrumentation({});

fetchInstrumentation.setTracerProvider(provider);

provider.addSpanProcessor(new SimpleSpanProcessor(consoleExporter));

provider.addSpanProcessor(new SimpleSpanProcessor(collectorExporter));

provider.register();

registerInstrumentations({
  instrumentations: [fetchInstrumentation],

  tracerProvider: provider,
});

export default function TraceProvider({ children }) {
  return <>{children}</>;
}
