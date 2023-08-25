# Charting Wall Street

**Charting Wall Street: Detecting Trading Patterns in Today's Automated Landscape**

Written as part of the master's thesis of Rudolf Aelbrecht for the Master's in Computer Science course at Ghent University.

The complete text of the master's thesis is available [here](./charting-wall-street.pdf).

## Summary

### Introduction

Technical analysis has long relied on chart patterns like "head and shoulders" or "double top/bottom." While these patterns are often referenced, their empirical validity and profitability remain questionable. The advent of machine learning has given rise to new 'Generic Patterns,' adding another layer of complexity. This dissertation examines both these types of patterns, probing their differences and effectiveness.

### Methodology

The study centered on S&P 500 companies from 1986 to 2022, managing the large dataset through a tailored microservices architecture. Two approaches were tested for pattern detection:
1. **Machine Learning Clustering**: Targeting generic patterns.
2. **Script-based Method**: Targeting traditional patterns, showing greater effectiveness.

The pattern performance was assessed using the Triple Barrier Method and supplemented with statistical analysis.

### Findings

- **Traditional Patterns**: Varied efficacy was found across different patterns, with some demonstrating potential but modest effect size.
- **Generic Patterns**: Struggled to display significant advantages over traditional forms in real-world trading, hinting at potential intrinsic limitations.

### Conclusion

Traditional patterns do exhibit behavior distinct from random trading but often with modest effect sizes, questioning their realistic effectiveness in today's trading scenarios. Integration with other indicators may enhance reliability. The custom tools and methodologies developed open the door for further exploration, with the potential for advanced machine learning techniques to deepen the understanding of trading patterns and their validity in the modern trading landscape.
