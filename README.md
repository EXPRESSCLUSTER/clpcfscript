# clpcfscript
- This command converts clp.conf (cluster configuration file) to clpcfset commands to recreate clp.conf.

## How to Use 
1. Get the cluster configuration file as below.
   ```sh
   clpcfctrl --pull -l -x .
   ```
1. Delete the first line in clp.conf.
   ```xml
   <?xml version="1.0" encoding="UTF-8"?>
   ```
1. Download clpcfscript command and save the clp.conf file on the same directory.
1. Run clpcfscript command.
   ```sh
   clpcfscript
   ```
1. create-cluster.sh file will be created on conf directory.
1. Run create-cluster.sh and you can get the clp.conf file.

## Notes
- This command can convert the following parameters.
  - Cluster name, encode, OS.
  - Add servers.
  - Add LAN Kernel Heartbeat (lankhb).
  - Add groups.
  - Add resources.
  - Add Mirror Disk parameters.
    - Sync mode only.