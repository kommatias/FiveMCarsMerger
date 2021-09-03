# FiveM Cars Merger
Merge FiveM cars into a single resource

## Usage
- Download the binary from [here](https://github.com/waseem-h/FiveMCarsMerger/releases/latest) or build it
- Save it in a directory
- Create a new folder in the same directory with any name
- Copy all your car folders and files in that folder (It doesn't matter what heirarchy those files are placed)
- Run the program with these parameters:
```bash
      --clean                Clear output directory before merging
      --input-path string    Path to all cars (default ".")
      --output-path string   Output path (default "out")
      --verbose              Enable verbose logging
```
- Example:
```bash
./FiveMCarsMerger --clean --input-path "cars" --output-path "merged-cars"
```
