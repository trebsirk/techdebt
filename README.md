# techdebt

### entropy

Entropy is a useful measure to model technical debt. Entropy here is calculated using Shannon's entropy formula, $H(X) = -\sum_{i} p(x_i)\log{p(x_i)}$. The values (the $X$ in $H(X)$) are an array of integers based on count of commits. The lower the entropy, the more technical debt. 

Entropy is calculated at the level of the file. 

```bash
make build
bin/techdebt 
```

##### functionality
    1. prints commits log style
    2. prints entropy for each file 

##### todos
    1. Entropy can be calculated at the level of the file, the repo, and the author. 
    2. Offer suggestions (prescriptive) for who could commit to which file to maximize repo entropy. (Low entropy is higher tech debt, and high entropy is low tech debt.)
    3. Command line arguments for different features, outputs, etc. 

