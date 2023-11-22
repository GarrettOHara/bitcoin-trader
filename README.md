# Bitcoin Trader
This repository defines infrastructure for automated Bitcoin trades placed on Coinbase. Each trade will be executed on a schedule.

The goal of this project is to enable automated Dollar Cost Averaging on BTC without high overhead/exchange fees.

# FQA

### What is `quote_size` ?
> A currency pair such as BTC-USD has a format of a base:quote currency. When making market buy orders, you would need to fill in the quote_size wherein for this scenario it would be USD. If you would input “50” in the quote_size field, you will be buying $50 worth of BTC.

*source*
https://forums.coinbasecloud.dev/t/market-order-example-on-how-to-use-quote-size-and-base-size/2640/2

