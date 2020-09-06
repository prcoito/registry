# registry
Go lib that allows to read Windows registry files

## Objectives
This lib pretends to offer an alternative to the golang's sys/windows/registry module. The exported functions are equivalent to the exported functions of golang's registry module to allow easier switching.

## Use cases
 - Registry access in Unix systems.
 - Access to registry files without loading to Windows registry


## Project status
Currently this project is still on initial phase. Accessing keys, subkeys and values is possible (read only).
There is work to be done in error handling and optimizations to be done.

## Thanks
This project would not be possible without the work of Timothy D. Morgan (http://www.sentinelchicken.com/data/TheWindowsNTRegistryFileFormat.pdf) and Joachim Metz (https://github.com/libyal/libregf)
