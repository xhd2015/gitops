## Steps
(Note: This is an edge case that may be hard to reproduce reliably in a unit test.
The scenario is: git status returns a directory entry, but between that call and expandDirsToFiles,
the directory is deleted by another process. expandDirsToFiles should handle the os.Stat error gracefully.)

1. No additional setup needed.
