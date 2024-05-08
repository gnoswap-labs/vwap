# VWAP

This repository implements a Volume Weighted Average Price (VWAP) calculator for the gno.land ecosystem. It retrieves trading data from the Gnoswap RPC API, calculates the exchange ratios between tokens, and computes the VWAP for each token pair.

## VWAP Calculation

The VWAP is calculated using the following formula:

$$
VWAP =
\begin{cases}
\frac{\sum_{i=1}^{n}{(Price_i \times Volume_i)}}{\sum_{i=1}^{n}{Volume_i}} & \text{if } \sum_{i=1}^{n}{Volume_i} > 0 \\
LastActivePrice & \text{if } \sum_{i=1}^{n}{Volume_i} = 0
\end{cases}
$$

where:

- $Price_i$ is the price of the $i$-th trade
- $Volume_i$ is the volume of the $i$-th trade
- $n$ is the total number of trades within the specified time window (e.g., last 10 minutes)
- $LastActivePrice$ is the last recorded price if there are no trades within the specified time window

The VWAP is calculated by summing up the product of price and volume for each trade within the specified time frame and dividing it by the total volume.

If there are no trades within the period, the last recorded price is used as the VWAP.

## Features

- Calculates the exchange ratios between tokens
- Computes the VWAP for each token pair based on the trading data from the last 10 minutes

## Pre-requisites

- Go version 1.22 or higher
