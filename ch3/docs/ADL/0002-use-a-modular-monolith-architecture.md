# 2. Use a modular monolith architecture

Date: 2022/05/10

## Status

Accepted

## Context

I want to avoid the deployment headaches by starting with a microservices architecture.

## Decision

I will use a modular monolith architecture to begin development with. The first modules will be the bounded contexts identified during the
Big Picture workshop.

- Store Management
- Kiosk Ordering
- Order Processing
- Automation & Fulfillment
- Payments

As we are using Go for our development; the code that belongs to a module will be placed under
the [/internal](https://go.dev/doc/go1.4#internalpackages) directory to keep it hidden from the other modules.

## Consequences

- Application runs as a single process
- Modules must use /internal to ensure autonomy
