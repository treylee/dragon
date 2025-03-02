## Architecture

This system is designed to execute performance tests and display real-time metrics. The React frontend acts as both the **Test Initiator** and the **Performance Dashboard**. When the test is started, the Go backend handles the performance test, and the React frontend updates with live metrics during the test.


Alternatively, you can hit these endpoints directly:

- `POST /start` - Starts the performance test.
- `GET /status?testID={{testID}}` - Gets the current status of the test.
- `GET /tests?offset=0&limit=8` - Fetches paginated results of the tests.

## Postman Collection

You can download the Postman collection to get started with the available API endpoints for this system.

[Download Postman Collection](./assets/Tester.postman_collection)

To Monitor the live streaming results create a seperate websocket postman script
- ws://localhost:3001/ws?testid={{testID}}

```mermaid
graph LR
    A[React Frontend] -->|POST /start| B[Go Backend]
    B -->|Start Performance Test| C[Performance Test]
    C -->|Execute Test| D[Test Execution]
    D -->|Live Update| A[Live Update Dashboard]
    
    A -->|GET /status| E[Get Test Status]
    E -->|Query Test Status| B
    B -->|GET /metrics| F[Fetch Test Metrics]
    F -->|Return Metrics| C
    C -->|Send Metrics| A

    A -->|GET /results| G[Request Test Results]
    G -->|Query Results| B
    B -->|GET /results/:testID| H[Retrieve Test Results]
    H -->|Return Test Results| C
    C -->|Send Results| A
    
    A -->|GET /paginated-results| I[Paginate Test Results]
    I -->|Query Paginated Results| B
    B -->|GET /paginated-results?offset=0&limit=10| J[Fetch Paginated Results]
    J -->|Return Paginated Results| C
    C -->|Send Paginated Results| A

    classDef frontend fill:#E0F7FA,stroke:#00796B,stroke-width:2px;
    classDef backend fill:#FFEB3B,stroke:#FF9800,stroke-width:2px;
    classDef test fill:#C5CAE9,stroke:#3F51B5,stroke-width:2px;
    classDef textColor fill:#FFFFFF,stroke:none,color:#000000;
    
    class A frontend,textColor;
    class B backend,textColor;
    class C test,textColor;
    class D test,textColor;
    class E frontend,textColor;
    class F test,textColor;
    class G frontend,textColor;
    class H test,textColor;
    class I frontend,textColor;
    class J test,textColor;