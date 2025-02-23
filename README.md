## Architecture

This system is designed to execute performance tests and display real-time metrics. The React frontend acts as both the **Test Initiator** and the **Performance Dashboard**. When the test is started, the Go backend handles the performance test, and the React frontend updates with live metrics during the test.

```mermaid
graph TD;
    A[React Frontend (Performance Dashboard)] -->|Initiates Test| B[Go Backend (Test Processor)];
    B -->|Runs Performance Test & Collects Metrics| C[Performance Metrics (Avg Response Time, Requests per Second, etc.)];
    C -->|Aggregates Data| D[Live Updates to React Dashboard];
    D -->|Updates Dashboard| A;

    classDef frontend fill:#f9f,stroke:#333,stroke-width:4px;
    classDef backend fill:#fdf,stroke:#333,stroke-width:4px;
    classDef metrics fill:#f8f8f8,stroke:#333,stroke-width:4px;

    class A frontend;
    class B backend;
    class C metrics;
    class D frontend;
