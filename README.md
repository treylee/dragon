## Architecture

This system is designed to execute performance tests and display real-time metrics. The React frontend acts as both the **Test Initiator** and the **Performance Dashboard**. When the test is started, the Go backend handles the performance test, and the React frontend updates with live metrics during the test.

```mermaid
graph LR
    A[React Frontend] --> B[Start Test]
    B --> C[Go Backend]
    C --> D[Run Performance Test]
    D --> A
    Frontend --> A[Live Update Dashboard]
 
 .