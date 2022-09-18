# 1. Keep an architecture decision log

Date: 2022/05/10

## Status

Accepted

## Context

To allow others to know why an architectural decision was made before their time or without their input.

## Decision

Keep an architectural decision log that will record all decisions regarding architecture and infrastructure choices.

The records will use the following markdown.

    # {RecordNum}. {Title}

    Date: {YYYY/MM/DD}

    ## Status
    {Pending,Accepted,Rejected,Superceded By {RecordNum|Link},Deprecated}

    ## Context
    {What is the issue that we're seeing that is motivating this decision or change?}

    ## Decision
    {What is the change that we're proposing and/or doing?}

    ## Consequences
    {What becomes easier or more difficult to do because of this change}

## Consequences

All decisions will be kept in the /docs/ADL directory.
