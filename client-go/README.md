# Text Client

Text client is a Go binary that can connect to the API and display the results
in a nicely formatted table:

```
+----+--------------------------------+----------+--------+------------+
| ID |             TITLE              | PRIORITY | STATUS |    DUE     |
+----+--------------------------------+----------+--------+------------+
|  1 | Write documentation            | medium   | to do  | 2020-05-18 |
|  2 | Fix a bug                      | high     | to do  | 2020-05-17 |
|  3 | Investigate Kubernetes         | low      | to do  | 2020-05-19 |
|    | deployment                     |          |        |            |
|  4 | Prepare demo                   | high     | to do  | 2020-05-17 |
+----+--------------------------------+----------+--------+------------+
```
