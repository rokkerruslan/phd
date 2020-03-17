# Mi - the migration tool for PostgreSQL

## Quick start

```shell script
$ go get github.com/rokkerruslan/mi

# Check that all right
$ mi --version  # or -V
...
```

```shell script
# Export variables
$ export DATABASE_URL=...

# Create new migration
$ mi new <name>

# Migrate to migration with 002 number
$ mi up 2

# If current migration is 3 we can down to 1
$ mi up 1

# Show migration history
$ mi history # or mi hi

# Show status (Unapplied migrations)
$ mi status # or mi st

# Squash migrations
$ mi squash 1:2

# Merge migrations
$ mi merge
```
