<!-- Limitations -->

Since there are many ways to evaluate patterns due to there being many financial markets, we will limit ourselves to
stock data. Stock data is still too broad. We will try to evaluate our patterns on a set of stocks chosen by their
market cap each year for as long as historical data is available. This makes sure that we do not only include companies
that are successful today, but have been successful at some point in history and may no longer exist.

We also require that a company has a large enough set of data, which means most modern tech companies are excluded from
our dataset. In the end, we settled with just over 500 stocks that contain more than 10 years of data each for this
experiment. Since patterns are claimed to work best on large time horizons, we will only evaluate our detection
algorithms on daily and larger time intervals.

Although we broadly evaluate the quality of the data, there might be some gaps in the dataset's quality. More
specifically, at the start of each stock's timeline, data was not yet digitized, and those data points do not provide as
much depth as more recent ones. Since we want to evaluate patterns over a large timeframe, we do not filter these points
out to gain insights into patterns from the past.