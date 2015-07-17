# Optimize SQL ALTER TABLE

Simple tool to rewrite ALTER TABLE statements into optimized batches.

## Usage

Having one or more SQL files like this:
```sql
INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;
INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;
INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;
INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;

ALTER TABLE table ADD field0 INT;
ALTER TABLE table ADD field1 INT;
ALTER TABLE table ADD field2 INT;
ALTER TABLE table ADD field3 INT;
ALTER TABLE table ADD field4 INT;

INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;
INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;
```

Simply pipe the file (or files) to "osat" to produce the optimized version on standard
output:

```sh
$ osat < files*.sql >result.sql
```

The result will look something like this:
```sql
INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;
INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;
INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;
INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;

ALTER TABLE table ADD field0 INT,
		ADD field1 INT,
		ADD field2 INT,
		ADD field3 INT,
		ADD field4 INT;

INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;
INSERT INTO table (a,b,c) VALUES (1,2,3),(4,5,6)...;
````

## Warning

The whole system is designed to be dump and might get fooled easily.  Do check the output
manually for extra safety.

## The Problem

I run into a bunch of (auto-generated) SQL files. Each contains both DELETE/UPDATE/INSERT
statements and ALTER TABLE statements.  The files are going to be executed as a database
migration.

Having fewer ALTER TABLE statements decreases the total run time of the migration a lot,
since the SQL server will probably create temporary tables for each ALTER TABLE statement
and some of the tables are quite big.

The solution is simply to run as few ALTER TABLE statements as possible.

# The Solution

As the ALTER TABLE statement are surrounded by statements that expect the tables to have
a precise layout, the ALTER TABLE statements cannot be safetly grouped all together at
the beginning of the file (or files).

For this reason, only consecutive ALTER TABLE statements are grouped together. (Even empty
line might break grouping.)

As limited as it is, it serves its purpose.  Also, it's a basic example of how concurrent
programming[1] is a paradigm in itself to solve problems, including sequential routines.

## Building

```sh
$ go get github.com/dullgiulio/osat
```

## References

[1] R. Pike, Concurrency is not Parallelism; https://www.youtube.com/watch?v=cN_DpYBzKso
