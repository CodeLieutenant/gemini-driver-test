create table test.table1
(
    pk0  timeuuid,
    pk1  date,
    ck0  float,
    ck1  ascii,
    col0 inet,
    col1 frozen<list<timeuuid>>,
    col2 text,
    col3 date,
    col4 frozen<set<time>>,
    col5 double,
    col6 blob,
    col7 text,
    primary key ((pk0, pk1), ck0, ck1)
) with caching = {'keys': 'ALL', 'rows_per_partition': 'ALL'} and compaction = {'class': 'SizeTieredCompactionStrategy'}
   and compression = {'sstable_compression': 'org.apache.cassandra.io.compress.LZ4Compressor'}
   and speculative_retry = '99.0PERCENTILE';
