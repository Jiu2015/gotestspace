## Change Log
### [1.0.0] - 2021-12-30
The gotestspace offical release  
### ADD
* Add internal environments: CALLER and CALLER_DIR, these two environments will hold the file name and dirname for source file while call `testspace.Create()`
* Add GetMultiPath method for join multiple dir names.
* Add `Cleaner` and `Cleaner` custom method, it will be run while gotestspace run `cleanup`.
* Add gotestspace Error type, it wraps the gotestspace common error, you can find stdout and stderr on it.

### Changed
* Update the command type with interface, adds better scalability.

### Fixed
* Fix the `Cleaner`  does not work bug
* Fix internal output and outErr doesn't have value bug
