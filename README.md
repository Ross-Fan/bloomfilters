## Bloom Filter for Recommendation System
The Bloom Filter is widely known for its mechanism of checking whtether an element exists in a large set. Here is an example of a Bloom filter used in a recommendation system for item deduplication for individuals.
Usually, a bloom filter denpend on three parameters: 
>    1. ***N***, the number of items in the filter (capacity limit of a filter)
>    2. ***M***, the number of bits in the filter
>    3. ***K***, the number of hash function

In the following picture, the bloom filter's capacity is 3000, the bits are 40000, the number of hash function is 4, the false positive rate is 0.4% , ref [Bloom Filter Calculator](https://hur.st/bloomfilter/?n=3000&p=&m=40000&k=4)
![alt](https://github.com/Ross-Fan/bloomfilters/blob/main/10.48.36.png)

In practice, the life cycle of a bloom fliter should be considered. So, when the total number of items reach the capacity (N), the replacement shall be executed. 
### Replacement
Assuming ***P*** bloom filters, each of the bloom filters has the same value of N,M,K, when all ***P*** bloom filters reach the capacity limit, the first bloom filter should be removed, and a new bloom filter with ***N***=0 should then be added at the tail of the bloom filters.
### Check
During the checking procedure, each bloom filter should be checked for a given item. If the bits (the hash result of the item) fulfill the condition, the item should be considered to exist.
### Set
In the Set phase, the Set should be continued with the bloom filter that has not yet reached the capacity limit. If all the bloom filter have reached the capacity, the replacement should be executed, and then return to the Set process if it is not finished.

### Hash Function
murmur3-32 is chosen as the hash function. In this program, the murmur32 is provided, the perfermance tested as the following result:
![alt](https://github.com/Ross-Fan/bloomfilters/blob/main/10.17.29.png)