# techdebt

### entropy

Entropy is a useful measure to model technical debt. Entropy here is calculated using Shannon's entropy formula, $H(X) = -\sum_{i} p(x_i)\log{p(x_i)}$. The values (the $X$ in $H(X)$) are an array of integers based on count of commits. 

Entropy is calculated at the level of the file, the repo, and the author. 

```bash
techdebt git entropy
```