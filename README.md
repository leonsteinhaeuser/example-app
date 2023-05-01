# example-app

A simple microservice app

```mermaid
graph TD
  A[Client] --> B[View Service]

  B --> C[Number Service]

  B --> D[Article Service]
  D --> D1[Article DB]

  B --> E[User Service]
  E --> E1[User DB]

  B --> F[Keyword Service]
  F --> F1[Keyword DB]

  B --> G[Article Comment Service]
  G --> G1[Article Comment DB]
```
