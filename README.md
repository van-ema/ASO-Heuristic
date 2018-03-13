# ASO-Heuristic

## Usage
To build and execute the script you need a goLang release. You can download it from the official site:
https://golang.org/dl/ .
After the installation of the go compiler to execute the script navigate into the project directory "orderSchedulingAlgorithm"

```
go build .
./orderSchedulingAlgorithm
```

**Typical usage:**

Optionally you can specify flags value to use different files containing the distance matrix and the delivery times vector.

```
Usage: orderSchedulingAlgorithm [-d value] [-m value] [-n value] [-t value] [parameters ...]
 
 -d, --distanceMat=value
       distance matrix filename
       
 -m, --nmover=value
       number of movers
       
 -n, --nOrder=value
       number of orders
       
 -t, --deliveryTimes=value
       delivery times vector filename
```

