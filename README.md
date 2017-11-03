# Bid Request Generator

Student programming challenge 2017

## Install

    go get -u github.com/hamaxx/bid_request_generator


## Run

    bid_request_generator [requests_per_sec] [concurrency]

Example:

    bid_request_generator 1000 1

## Output

Program writes randomly generated logs to stdout.
It produces 3 types of logs (bid, win, click) with different rates.

Rates:

- Bids: [requests\_per\_sec]
- Wins: Bids / 10
- Click: Wins / 100
