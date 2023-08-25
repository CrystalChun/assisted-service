# Minimum Reuirements

The following tables show the minimum hardware requirements needed for a successful installation of each cluster type.

## Multi-Node Cluster

| Role | CPU (cores) | Memory (MiB) | Disk (GB)| Install Disk Speed Threshold (ms) | Packet Loss (%)|
|------|-------------|-----------|----------|-----------------------------------|--------|
| master | 4 | 16384 | 100 | 10 | 100 | 0 |
| worker | 2 | 8192 | 100 | 10 | 1000 | 10 |

## SNO Cluster

| Role | CPU (cores) | Memory (MiB) | Disk (GB)| Install Disk Speed Threshold (ms) | Packet Loss (%)|
|------|-------------|-----------|----------|-----------------------------------|--------|
| master | 8 | 16384 | 100 | 10 | 100 | -- |

## Edge Worker

| Role | CPU (cores) | Memory (MiB) | Disk (GB)| Install Disk Speed Threshold (ms) | Packet Loss (%)|
|------|-------------|-----------|----------|-----------------------------------|--------|
| worker | 2 | 8192 | 15 | 10 | -- | -- |
